# Base attestor-trusted L2 client

This client trusts the relayer to forward an attestor-produced L2 execution header. It verifies
the header hash, configured fork, and router account proof, but does not verify L1, dispute games,
or attestor signatures. This is a bring-up trust boundary, not a trustless light client.
