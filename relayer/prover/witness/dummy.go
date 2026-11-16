package witness

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sync"
)

// Dummy padding fills the unused slots of a bucket. Every slot in a batch is
// Ed25519-verified in-circuit — padding is not gated — so each dummy carries a
// real signature from a deterministic keypair. The padding keys differ from one
// another by derivation, but nothing keeps them apart from the active slots'
// points, and the circuit does not need that: a batch whose active signer
// repeats a padding slot's A or R still proves (#579 item 4,
// TestBatchCircuit_SolvesWhenASignerCollidesWithADummy). What marks a slot as
// padding is the committed Active byte plus the on-chain quorum check.
//
// Neither string below is checked on-chain: SignatureVerifier hashes only each
// slot's Active byte and pubkey, and the dummy pubkeys travel in calldata. Only
// DummyMsgLen reaches the circuit (as its SHA-512 MinMsgLen), and it is fixed
// independently of dummyMagic, so renaming either string leaves the R1CS and
// every verifying key unchanged.

// DummyMsgLen is the length of every padding slot's signed-bytes buffer — the
// shortest message any slot carries, hence the floor the circuit hands SHA-512
// as minimal input length (ZK-05). It is baked into the R1CS: changing it
// requires re-setup + redeploying every Groth16Verifier_N{N} + setBucket.
const DummyMsgLen = 18

// dummyMagicWidth is the fixed field dummyMagic is zero-padded into. Full
// layout: magic field(14) || bucket(2 BE) || slot(2 BE) = DummyMsgLen bytes.
const dummyMagicWidth = DummyMsgLen - 4

// dummyMagic prefixes every dummy slot's signed-bytes buffer. It may be any
// length up to dummyMagicWidth; a longer value fails to compile below.
const dummyMagic = "spectre-dummy"

var _ [dummyMagicWidth - len(dummyMagic)]struct{} // compile-time: len(dummyMagic) <= dummyMagicWidth

// dummySeedFmt derives each slot's deterministic keypair seed.
const dummySeedFmt = "spectre-dummy-seed|bucket=%d|slot=%d"

// DummyMsgBytes is the deterministic signed-bytes buffer for a padding slot:
// dummyMagic zero-padded to dummyMagicWidth || bucket(2 BE) || slot(2 BE).
func DummyMsgBytes(bucket, slot uint16) []byte {
	out := make([]byte, DummyMsgLen)
	copy(out, dummyMagic)
	binary.BigEndian.PutUint16(out[dummyMagicWidth:], bucket)
	binary.BigEndian.PutUint16(out[dummyMagicWidth+2:], slot)
	return out
}

// PadWithDummies extends sigs to exactly n slots with deterministic dummy
// signatures (Active=false). Deterministic in (n, slot), so calldata is
// reproducible across runs.
func PadWithDummies(sigs []ValidatorSignature, n int) ([]ValidatorSignature, error) {
	if len(sigs) > n {
		return nil, fmt.Errorf("bucket n=%d too small for %d signers", n, len(sigs))
	}
	if n > 0xFFFF {
		return nil, fmt.Errorf("bucket size %d exceeds 16-bit dummy slot encoding", n)
	}
	out := make([]ValidatorSignature, n)
	copy(out, sigs)
	for i := len(sigs); i < n; i++ {
		priv := dummyKey(n, i)
		msg := DummyMsgBytes(uint16(n), uint16(i))
		out[i] = ValidatorSignature{
			Signature:   ed25519.Sign(priv, msg),
			PublicKey:   priv.Public().(ed25519.PublicKey),
			SignedBytes: msg,
			Active:      false,
		}
	}
	return out, nil
}

// dummyKey derives the padding keypair of slot i in bucket n.
func dummyKey(n, i int) ed25519.PrivateKey {
	seed := sha256.Sum256(fmt.Appendf(nil, dummySeedFmt, n, i))
	return ed25519.NewKeyFromSeed(seed[:])
}

// dummyKeyIndex maps the padding public key of every slot of every supported
// bucket to that (bucket, slot). It is built once, on first use, and only read
// after that.
var (
	dummyKeyIndexOnce sync.Once
	dummyKeyIndex     map[string][2]int
)

// DummyKeySlot reports whether publicKey is the padding key PadWithDummies uses
// for a slot of a supported bucket, and which one. The padding keys come from
// public seeds, so their private keys are public: anyone can sign as a validator
// that uses one, and no honest key generation produces one. The circuit still
// proves a batch in which such a validator's key repeats a padding slot's key
// (#579 item 4), so this only tells the relayer whom to warn about.
func DummyKeySlot(publicKey []byte) (bucket, slot int, ok bool) {
	dummyKeyIndexOnce.Do(func() {
		dummyKeyIndex = make(map[string][2]int)
		for _, n := range SupportedBuckets {
			for i := 0; i < n; i++ {
				dummyKeyIndex[string(dummyKey(n, i).Public().(ed25519.PublicKey))] = [2]int{n, i}
			}
		}
	})
	at, ok := dummyKeyIndex[string(publicKey)]
	return at[0], at[1], ok
}
