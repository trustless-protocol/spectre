package witness

import "testing"

func TestSmallestBucketGEQ(t *testing.T) {
	loaded := []int{4, 16} // a set with a hole: n8 not built
	cases := []struct {
		name    string
		n       int
		want    int
		wantErr bool
	}{
		{"one signer -> smallest loaded", 1, 4, false},
		{"exact fit", 4, 4, false},
		{"hole in the set: skips to the next loaded bucket", 6, 16, false},
		{"zero rejected", 0, 0, true},
		{"negative rejected", -1, 0, true},
		{"over the largest loaded", 17, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SmallestBucketGEQ(tc.n, loaded)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got bucket %d", got)
				}
				return
			}
			if err != nil || got != tc.want {
				t.Fatalf("SmallestBucketGEQ(%d, %v) = %d, %v; want %d", tc.n, loaded, got, err, tc.want)
			}
		})
	}
}

func TestIsSupported(t *testing.T) {
	for _, n := range SupportedBuckets {
		if !IsSupported(n) {
			t.Fatalf("supported bucket %d reported unsupported", n)
		}
	}
	for _, n := range []int{0, 5, 129} {
		if IsSupported(n) {
			t.Fatalf("bucket %d must not be supported", n)
		}
	}
}
