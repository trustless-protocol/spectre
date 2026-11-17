package witness

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
)

// encoding returns a 32-byte compressed point: first and last bytes as given,
// every byte in between set to fill.
func encoding(first, fill, last byte) []byte {
	out := make([]byte, 32)
	for i := range out {
		out[i] = fill
	}
	out[0], out[31] = first, last
	return out
}

func TestCheckCircuitEncoding(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	sig := ed25519.Sign(priv, []byte("vote"))

	withR := func(r []byte) []byte { return append(append([]byte(nil), r...), sig[32:]...) }
	withS := func(s []byte) []byte { return append(append([]byte(nil), sig[:32]...), s...) }
	// ℓ itself, little-endian: the smallest S the circuit rejects.
	ell := []byte{
		0xed, 0xd3, 0xf5, 0x5c, 0x1a, 0x63, 0x12, 0x58, 0xd6, 0x9c, 0xf7, 0xa2, 0xde, 0xf9, 0xde, 0x14,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10,
	}
	ellMinusOne := append([]byte(nil), ell...)
	ellMinusOne[0]--

	cases := []struct {
		name string
		pub  []byte
		sig  []byte
		want error
	}{
		{"honest signature", pub, sig, nil},
		{"A identity (0, 1)", encoding(0x01, 0x00, 0x00), sig, ErrPointXZero},
		{"A identity with the sign bit set", encoding(0x01, 0x00, 0x80), sig, ErrPointXZero},
		{"A order-2 point (0, -1)", encoding(0xec, 0xff, 0x7f), sig, ErrPointXZero},
		{"A y = p + 1 (identity written non-canonically)", encoding(0xee, 0xff, 0x7f), sig, ErrPointNotCanonical},
		{"A y = p (order-4 point written non-canonically)", encoding(0xed, 0xff, 0x7f), sig, ErrPointNotCanonical},
		{"A y = 2^255 - 1", encoding(0xff, 0xff, 0x7f), sig, ErrPointNotCanonical},
		// y = 0 is the canonical order-4 point: low order, but x != 0, and the
		// circuit proves it (TestCheckCircuitEncodingMatchesCircuit).
		{"A order-4 point y = 0", encoding(0x00, 0x00, 0x00), sig, nil},
		{"A y = p - 2", encoding(0xeb, 0xff, 0x7f), sig, nil},
		{"R identity (nonce r = 0)", pub, withR(encoding(0x01, 0x00, 0x00)), ErrPointXZero},
		{"R order-2 point", pub, withR(encoding(0xec, 0xff, 0x7f)), ErrPointXZero},
		{"R y = p", pub, withR(encoding(0xed, 0xff, 0x7f)), ErrPointNotCanonical},
		{"S = ℓ", pub, withS(ell), ErrScalarNotCanonical},
		{"S = ℓ - 1", pub, withS(ellMinusOne), nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CheckCircuitEncoding(tc.pub, tc.sig)
			if tc.want == nil {
				if got != nil {
					t.Fatalf("CheckCircuitEncoding = %v, want nil", got)
				}
				return
			}
			if !errors.Is(got, tc.want) {
				t.Fatalf("CheckCircuitEncoding = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCheckCircuitEncodingRejectsBadLengths(t *testing.T) {
	if err := CheckCircuitEncoding(make([]byte, 31), make([]byte, 64)); err == nil {
		t.Fatal("a 31-byte public key was accepted")
	}
	if err := CheckCircuitEncoding(make([]byte, 32), make([]byte, 63)); err == nil {
		t.Fatal("a 63-byte signature was accepted")
	}
}

// TestCheckCircuitEncodingNamesTheField keeps the field in the message, which is
// what an operator reads in the relayer's warning.
func TestCheckCircuitEncodingNamesTheField(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	sig := ed25519.Sign(priv, []byte("vote"))
	err := CheckCircuitEncoding(encoding(0x01, 0x00, 0x00), sig)
	if err == nil || err.Error() != "public key: "+ErrPointXZero.Error() {
		t.Fatalf("public key error = %v", err)
	}
	badR := append(encoding(0x01, 0x00, 0x00), sig[32:]...)
	err = CheckCircuitEncoding(pub, badR)
	if err == nil || err.Error() != "signature R: "+ErrPointXZero.Error() {
		t.Fatalf("R error = %v", err)
	}
}
