//! OP optimistic update types and verifier.

/// Verifies an OP dispute-game commitment without waiting for resolution.
pub use op_stack_verifier::verify;
/// OP optimistic update header.
pub use op_stack_verifier::Header;
/// OP Stack output-root proof supplied by the relayer.
pub use op_stack_verifier::OutputRootProof;
/// Conservative EVM proof limits.
pub use op_stack_verifier::PROOF_LIMITS;
