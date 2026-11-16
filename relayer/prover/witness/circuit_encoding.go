package witness

import (
	"errors"
	"fmt"
	"math/big"
)

// The batch circuit is stricter than crypto/ed25519.Verify about how a
// signature's points and scalar are encoded. Verify follows the ZIP-215 reading
// of RFC 8032 and accepts a public key A whose y coordinate is written as y + p,
// and it accepts any point with x = 0. The circuit rejects all of them:
//
//   - it asserts x != 0 for both A and R (ecip-gnark signature/eddsa
//     assertOnEdwardsCurve), which rules out the identity (0, 1) and the
//     order-2 point (0, -1). R is the identity when the signer used nonce r = 0;
//   - it hashes A's canonical re-encoding into the witness digest, while the
//     relayer and SignatureVerifier hash the raw bytes, so the public input never
//     matches when A is written non-canonically;
//   - it asserts S < ℓ (Verify does too).
//
// An honest signer produces none of these: a random key or nonce lands on them
// with negligible probability. A validator that does is choosing to, and if the
// relayer selected it, every proof over that quorum would fail and the client
// would stop updating for as long as it stays selected.
var (
	// ErrPointNotCanonical: y is written as y + p (y >= 2^255 - 19).
	ErrPointNotCanonical = errors.New("point is not canonically encoded (y >= 2^255-19)")
	// ErrPointXZero: the point is the identity or the order-2 point (0, -1).
	ErrPointXZero = errors.New("point has x = 0 (the identity or the order-2 point)")
	// ErrScalarNotCanonical: S >= ℓ.
	ErrScalarNotCanonical = errors.New("signature scalar S is not below the group order")
)

var (
	// curveP is the Curve25519 base-field modulus 2^255 - 19.
	curveP = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(19))
	// curvePMinusOne is the y coordinate of the order-2 point (0, -1).
	curvePMinusOne = new(big.Int).Sub(curveP, big.NewInt(1))
	// groupOrder is ℓ = 2^252 + 27742317777372353535851937790883648493.
	groupOrder, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)
)

// CheckCircuitEncoding reports why the batch circuit cannot prove a signature
// that crypto/ed25519.Verify accepts, or nil when it can. It only inspects
// encodings: call it on a signature Verify has already accepted, which is what
// guarantees the points decode onto the curve and the signature equation holds.
// The circuit checks that equation cofactored, which every signature Verify
// accepts also satisfies, so the encodings are the only gap.
func CheckCircuitEncoding(publicKey, signature []byte) error {
	if len(publicKey) != 32 {
		return fmt.Errorf("public key must be 32 bytes, got %d", len(publicKey))
	}
	if len(signature) != 64 {
		return fmt.Errorf("signature must be 64 bytes, got %d", len(signature))
	}
	if err := checkPointEncoding(publicKey); err != nil {
		return fmt.Errorf("public key: %w", err)
	}
	if err := checkPointEncoding(signature[:32]); err != nil {
		return fmt.Errorf("signature R: %w", err)
	}
	if littleEndianInt(signature[32:]).Cmp(groupOrder) >= 0 {
		return ErrScalarNotCanonical
	}
	return nil
}

// checkPointEncoding checks a compressed point: y in the low 255 bits, the sign
// of x in the top bit. On the curve -x^2 + y^2 = 1 + d*x^2*y^2, x = 0 exactly
// when y = 1 or y = p - 1, whatever the sign bit says.
func checkPointEncoding(enc []byte) error {
	y := littleEndianInt(enc)
	y.SetBit(y, 255, 0)
	if y.Cmp(curveP) >= 0 {
		return ErrPointNotCanonical
	}
	if y.Cmp(big.NewInt(1)) == 0 || y.Cmp(curvePMinusOne) == 0 {
		return ErrPointXZero
	}
	return nil
}

func littleEndianInt(le []byte) *big.Int {
	be := make([]byte, len(le))
	for i := range le {
		be[len(le)-1-i] = le[i]
	}
	return new(big.Int).SetBytes(be)
}
