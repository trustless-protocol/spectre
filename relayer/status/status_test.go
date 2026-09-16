package status

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"relayer/relay"
	"relayer/services"
)

func fixedSnapshot(_ context.Context) Snapshot {
	return Snapshot{
		Now:    time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC),
		Paths:  []relay.Status{{Path: "cosmos->eth", ClientID: "07-tendermint-0", TrustedHeight: 41, RelayableHeight: 42}},
		Queues: []services.QueueStatus{{Direction: "cosmos", Queued: 3}},
	}
}

// The endpoint exposes packet identities and client ids. It must never grow a
// way to CHANGE them by accident, so anything that is not a read is refused
// rather than routed.
func TestHandlerIsReadOnly(t *testing.T) {
	srv := httptest.NewServer(Handler(fixedSnapshot, nil))
	t.Cleanup(srv.Close)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, srv.URL+Path, nil)
			if err != nil {
				t.Fatalf("build request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("%s: %v", method, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Fatalf("%s returned %s, want 405: a read-only endpoint must refuse writes, not ignore them",
					method, resp.Status)
			}
			if allow := resp.Header.Get("Allow"); allow != http.MethodGet {
				t.Errorf("Allow = %q, want GET so a caller learns what is permitted", allow)
			}
		})
	}
}

// Fetch and Handler are two halves of one wire format. Testing them apart would
// let the pair drift while both stay green.
func TestFetchRoundTripsTheSnapshot(t *testing.T) {
	srv := httptest.NewServer(Handler(fixedSnapshot, nil))
	t.Cleanup(srv.Close)

	got, err := Fetch(context.Background(), strings.TrimPrefix(srv.URL, "http://"))
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	want := fixedSnapshot(context.Background())
	if len(got.Paths) != 1 || got.Paths[0].Path != want.Paths[0].Path {
		t.Fatalf("paths = %+v, want the served one", got.Paths)
	}
	if got.Paths[0].RelayableHeight != 42 || got.Paths[0].TrustedHeight != 41 {
		t.Fatalf("heights = %d/%d, want 41/42: these are the two numbers the command exists to show",
			got.Paths[0].TrustedHeight, got.Paths[0].RelayableHeight)
	}
	if len(got.Queues) != 1 || got.Queues[0].Queued != 3 {
		t.Fatalf("queues = %+v, want the served one", got.Queues)
	}
}

// An operator who cannot reach the endpoint is the common case -- it is off by
// default -- so the error has to say that rather than surface a bare dial error.
func TestFetchSaysWhyItCouldNotConnect(t *testing.T) {
	srv := httptest.NewServer(Handler(fixedSnapshot, nil))
	addr := strings.TrimPrefix(srv.URL, "http://")
	srv.Close() // nothing is listening now

	_, err := Fetch(context.Background(), addr)
	if err == nil {
		t.Fatal("Fetch succeeded against a closed endpoint")
	}
	if !strings.Contains(err.Error(), "status_address") {
		t.Fatalf("error = %v; want it to name the setting that turns the endpoint on", err)
	}
}

// Off is the default, and an engine must not have to know whether anyone is
// listening: Serve with no address returns rather than binding something.
func TestServeWithNoAddressBindsNothing(t *testing.T) {
	done := make(chan error, 1)
	go func() { done <- Serve(context.Background(), "", fixedSnapshot, nil) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Serve(\"\") = %v, want nil", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Serve(\"\") blocked; an unconfigured endpoint must not hold a goroutine open")
	}
}

// A bind failure has to be reported, not logged and forgotten: an operator told
// to query the endpoint must not find out during an incident that it never
// started.
func TestServeReportsABindFailure(t *testing.T) {
	taken := httptest.NewServer(Handler(fixedSnapshot, nil))
	t.Cleanup(taken.Close)

	err := Serve(context.Background(), strings.TrimPrefix(taken.URL, "http://"), fixedSnapshot, nil)
	if err == nil {
		t.Fatal("Serve returned nil for an address already in use")
	}
	if !strings.Contains(err.Error(), "status endpoint") {
		t.Fatalf("error = %v, want it to name what failed to start", err)
	}
}

// The write route is ABSENT rather than forbidden when controls are off. 405
// would say "this exists but not for you", which is the wrong thing to tell an
// operator whose next step is to turn it on.
func TestReactivateRouteExistsOnlyWhenControlsAreOn(t *testing.T) {
	t.Run("off by default", func(t *testing.T) {
		srv := httptest.NewServer(Handler(fixedSnapshot, nil))
		t.Cleanup(srv.Close)

		err := Reactivate(context.Background(), strings.TrimPrefix(srv.URL, "http://"), DeadLetterSelector{Path: "cosmos<->eth", Direction: "cosmos", SourceClient: "src-0", Sequence: 7})
		if err == nil {
			t.Fatal("reactivate succeeded against an endpoint with no write route")
		}
		if !strings.Contains(err.Error(), "status_allow_control") {
			t.Fatalf("error = %v; want it to name the setting that enables the route", err)
		}
	})

	t.Run("enabled explicitly", func(t *testing.T) {
		var got DeadLetterSelector
		ctl := &Controls{Reactivate: func(sel DeadLetterSelector) error {
			got = sel
			return nil
		}}
		srv := httptest.NewServer(Handler(fixedSnapshot, ctl))
		t.Cleanup(srv.Close)

		if err := Reactivate(context.Background(), strings.TrimPrefix(srv.URL, "http://"), DeadLetterSelector{Path: "cosmos<->eth", Direction: "cosmos", SourceClient: "src-0", Sequence: 7}); err != nil {
			t.Fatalf("reactivate: %v", err)
		}
		if want := (DeadLetterSelector{Path: "cosmos<->eth", Direction: "cosmos", SourceClient: "src-0", Sequence: 7}); got != want {
			t.Fatalf("handler received %+v, want %+v; every field has to survive the round trip, "+
				"because path and direction are what make the write unambiguous", got, want)
		}
	})
}

// Enabling the write route must not make the READ route writable, nor the write
// route readable: two routes, two methods, no overlap.
func TestControlRouteKeepsItsOwnMethod(t *testing.T) {
	ctl := &Controls{Reactivate: func(DeadLetterSelector) error { return nil }}
	srv := httptest.NewServer(Handler(fixedSnapshot, ctl))
	t.Cleanup(srv.Close)

	t.Run("the snapshot stays read-only", func(t *testing.T) {
		resp, err := http.Post(srv.URL+Path, "application/json", strings.NewReader("{}"))
		if err != nil {
			t.Fatalf("post: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("POST %s returned %s, want 405 even with controls on", Path, resp.Status)
		}
	})

	t.Run("the write route refuses a read", func(t *testing.T) {
		resp, err := http.Get(srv.URL + ReactivatePath)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("GET %s returned %s, want 405", ReactivatePath, resp.Status)
		}
	})
}

// A request that names no packet must be refused, not turned into a lookup for
// the empty client at sequence zero -- which would reach the tracker and refuse
// there, with a worse message.
func TestReactivateRefusesAnUnnamedPacket(t *testing.T) {
	called := false
	ctl := &Controls{Reactivate: func(DeadLetterSelector) error { called = true; return nil }}
	srv := httptest.NewServer(Handler(fixedSnapshot, ctl))
	t.Cleanup(srv.Close)

	// Each of these omits exactly one required field, plus the undecodable case.
	// Path and direction are in the list because they are what makes the write
	// name one tracker; accepting a request without them is the reported bug.
	for _, body := range []string{
		`{}`,
		`{"source_client":"src-0"}`,
		`{"sequence":7}`,
		`{"path":"cosmos<->eth","direction":"cosmos","source_client":"src-0"}`,
		`{"path":"cosmos<->eth","direction":"cosmos","sequence":7}`,
		`{"path":"cosmos<->eth","source_client":"src-0","sequence":7}`,
		`{"direction":"cosmos","source_client":"src-0","sequence":7}`,
		`not json`,
	} {
		resp, err := http.Post(srv.URL+ReactivatePath, "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("post %q: %v", body, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("POST %q returned %s, want 400", body, resp.Status)
		}
	}
	if called {
		t.Fatal("an unnamed packet reached the tracker")
	}
}

// The tracker's refusal is the caller's answer -- "that packet is not where you
// think it is" -- so it must arrive intact rather than as a bare 500.
func TestReactivateCarriesTheTrackersRefusal(t *testing.T) {
	ctl := &Controls{Reactivate: func(DeadLetterSelector) error {
		return errors.New("no dead-lettered packet for client src-0 sequence 7")
	}}
	srv := httptest.NewServer(Handler(fixedSnapshot, ctl))
	t.Cleanup(srv.Close)

	err := Reactivate(context.Background(), strings.TrimPrefix(srv.URL, "http://"), DeadLetterSelector{Path: "cosmos<->eth", Direction: "cosmos", SourceClient: "src-0", Sequence: 7})
	if err == nil {
		t.Fatal("a refused reactivation reported success")
	}
	if !strings.Contains(err.Error(), "no dead-lettered packet") {
		t.Fatalf("error = %v; want the tracker's own reason", err)
	}
}

// The write route must not be served where other people can reach it.
//
// Reported by @DongLieu: warnIfReachable covers the snapshot, whose risk is
// disclosure, and a warning is the right weight there. Reactivation CHANGES
// relayer state and nothing in this package authenticates it, so binding it to
// 0.0.0.0 publishes an unauthenticated mutation -- a decision an operator
// inherits from an example config rather than makes.
func TestServeRefusesControlsOffLoopback(t *testing.T) {
	ctl := &Controls{Reactivate: func(DeadLetterSelector) error { return nil }}

	t.Run("a non-loopback bind with controls on is refused", func(t *testing.T) {
		// Bounded on purpose: without the guard Serve does not return, it SERVES,
		// so an unbounded call here reports the defect as a wall-clock timeout
		// rather than as the defect. The mutation found that the first time.
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() { done <- Serve(ctx, "0.0.0.0:0", fixedSnapshot, ctl) }()

		var err error
		select {
		case err = <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("Serve did not return: it is serving the write route on 0.0.0.0 instead of refusing it")
		}
		if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("served the write route on 0.0.0.0 (err = %v); anyone who can reach it could reactivate packets", err)
		}
		for _, want := range []string{"status_allow_control", ReactivatePath, "loopback"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not mention %q; an operator has to be able to act on it", err, want)
			}
		}
	})

	// The read-only endpoint keeps binding anywhere: its risk is disclosure, it
	// already warns, and taking it away would remove the diagnosis this whole
	// package exists for.
	t.Run("the same bind is allowed read-only", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() { done <- Serve(ctx, "0.0.0.0:0", fixedSnapshot, nil) }()

		select {
		case err := <-done:
			t.Fatalf("read-only serve returned early: %v", err)
		case <-time.After(150 * time.Millisecond):
		}
		cancel()
		if err := <-done; err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf("read-only serve: %v", err)
		}
	})

	t.Run("loopback with controls is allowed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() { done <- Serve(ctx, "127.0.0.1:0", fixedSnapshot, ctl) }()

		select {
		case err := <-done:
			t.Fatalf("loopback serve with controls returned early: %v", err)
		case <-time.After(150 * time.Millisecond):
		}
		cancel()
		if err := <-done; err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf("loopback serve: %v", err)
		}
	})
}

// Reported by @neitdung. Refusing a non-loopback bind stops other HOSTS from
// reaching the control route; it does nothing about the browser already running
// on this one. A page on any site can POST to http://127.0.0.1:<port>/reactivate
// as a CORS "simple request" -- no preflight, so nothing on this server is ever
// asked for permission. The same-origin policy hides the RESPONSE from the page;
// it does not stop the state change, and a state change is the whole risk here.
func TestReactivateRefusesABrowserCrossOriginPost(t *testing.T) {
	body := `{"path":"cosmos->eth","direction":"cosmos","source_client":"08-wasm-0","sequence":7}`

	for _, tc := range []struct {
		name        string
		contentType string
		origin      string
	}{
		// text/plain is what makes the request "simple": the browser sends it
		// without asking this server first.
		{name: "simple content type from a page", contentType: "text/plain;charset=UTF-8", origin: "https://attacker.example"},
		{name: "form encoding from a page", contentType: "application/x-www-form-urlencoded", origin: "https://attacker.example"},
		{name: "multipart from a page", contentType: "multipart/form-data; boundary=x", origin: "https://attacker.example"},
		// An Origin header at all means a browser sent this. Operators use curl.
		{name: "json claimed by a page", contentType: "application/json", origin: "https://attacker.example"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var called int
			handler := reactivateHandler(func(DeadLetterSelector) error { called++; return nil })

			req := httptest.NewRequest(http.MethodPost, ReactivatePath, strings.NewReader(body))
			req.Header.Set("Content-Type", tc.contentType)
			req.Header.Set("Origin", tc.origin)
			rec := httptest.NewRecorder()
			handler(rec, req)

			if called != 0 {
				t.Fatalf("the control ran for a cross-origin browser POST (%s / %s); "+
					"the same-origin policy hides the reply, not the mutation", tc.contentType, tc.origin)
			}
			if rec.Code != http.StatusForbidden && rec.Code != http.StatusUnsupportedMediaType {
				t.Fatalf("status = %d, want the request refused before it reached the control", rec.Code)
			}
		})
	}

	// The content-type half on its own. Every case above also carries an Origin,
	// so the first check catches them and the second is never the thing that
	// refuses -- removing it left this whole test green, which is how a second
	// defence quietly becomes decoration. These carry NO Origin, so only the
	// content type can refuse them.
	t.Run("a simple content type is refused on its own merits", func(t *testing.T) {
		for _, contentType := range []string{
			"text/plain;charset=UTF-8",
			"application/x-www-form-urlencoded",
			"multipart/form-data; boundary=x",
			"", // a POST with no declared type at all
		} {
			var called int
			handler := reactivateHandler(func(DeadLetterSelector) error { called++; return nil })

			req := httptest.NewRequest(http.MethodPost, ReactivatePath, strings.NewReader(body))
			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}
			rec := httptest.NewRecorder()
			handler(rec, req)

			if called != 0 {
				t.Fatalf("the control ran for a %q body with no Origin; requiring application/json "+
					"is what forces a preflight this server never answers", contentType)
			}
			if rec.Code != http.StatusUnsupportedMediaType {
				t.Fatalf("status = %d for %q, want 415", rec.Code, contentType)
			}
		}
	})

	// The control must still work for the caller it exists for.
	t.Run("an operator's own request still works", func(t *testing.T) {
		var called int
		handler := reactivateHandler(func(DeadLetterSelector) error { called++; return nil })

		req := httptest.NewRequest(http.MethodPost, ReactivatePath, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler(rec, req)

		if called != 1 || rec.Code != http.StatusOK {
			t.Fatalf("a plain application/json POST with no Origin was refused: status=%d calls=%d body=%q",
				rec.Code, called, rec.Body.String())
		}
	})
}

// Reported by @neitdung. Serve waits on the shutdown goroutine unconditionally
// after srv.Serve returns, and that goroutine closes its channel only after the
// parent context is cancelled. So a listener that dies AFTER startup -- the one
// case where the error is news -- never reaches the caller: Serve blocks until
// something unrelated cancels the context, and the relayer runs on with a status
// endpoint that answers nothing.
//
// Bounded, because the defect is "does not return": an unbounded version reports
// it as a timeout of the whole package rather than as this test failing.
func TestServeReturnsAListenerFailureWithoutWaitingForShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // NOT cancelled during the assertion: that is the point

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	served := make(chan error, 1)
	go func() { served <- serveOn(ctx, listener, fixedSnapshot, nil) }()

	// Kill the listener out from under the running server.
	time.Sleep(20 * time.Millisecond)
	_ = listener.Close()

	select {
	case err := <-served:
		if err == nil {
			t.Fatal("a listener that died mid-run reported a clean stop")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Serve did not return after its listener failed; it is waiting for a context " +
			"cancellation that may never come, and the endpoint is dead in the meantime")
	}
}
