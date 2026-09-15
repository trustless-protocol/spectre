# ICS26Router Runtime Byte Budget

Forced fresh baseline and refactored runtime size: **24,465 bytes**, leaving **111 bytes** under EIP-170. The 01 target of 1 KiB margin would require at least 913 bytes of additional safe reduction.

Reachable runtime ownership remains:

- ICS-26 packet validation and commitment transitions;
- ICS-02 client registry, proof dispatch, and migration entrypoints;
- ICS-24 path/commitment logic;
- pause, access-management, upgrade, and deprecated-admin compatibility boundaries;
- application callback orchestration and emitted compatibility surface.

No dead or duplicate candidate with at least 913 bytes of proven behavior-preserving savings was identified. Therefore this refactor performs **no bytecode optimization**. Removing deprecated admin code, reordering flows, changing errors, or altering module calls would violate the compatibility contract without separate deployment inventory and review. Runtime bytecode remains exactly equal to origin/main.
