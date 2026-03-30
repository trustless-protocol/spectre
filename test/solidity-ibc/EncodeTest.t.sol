// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable custom-errors,max-line-length

import { Test, console } from "forge-std/Test.sol";

import { Encode } from "../../contracts/utils/Encode.sol";
import { Header } from "../../contracts/utils/Header.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";

/// @title EncodeTest
/// @notice Cross-validates Encode.sol output against Go proto.Marshal() reference data.
///         Go proto.Marshal() is the SOURCE OF TRUTH.
///         Run Go reference: cd operator && go run ./cmd/encode_debug/
contract EncodeTest is Test {
    // ─── Test data (same as Go cmd/encode_debug/main.go) ───

    function _pubKey() internal pure returns (bytes32) {
        bytes32 pk;
        assembly {
            pk := 0xaaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9
        }
        return pk;
    }

    function _hash1() internal pure returns (bytes32) {
        return bytes32(0x101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f);
    }

    function _hash2() internal pure returns (bytes32) {
        return bytes32(0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f);
    }

    function _hash3() internal pure returns (bytes32) {
        return bytes32(0x303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f);
    }

    function _hash4() internal pure returns (bytes32) {
        return bytes32(0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f);
    }

    function _hash5() internal pure returns (bytes32) {
        return bytes32(0x505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f);
    }

    function _hash6() internal pure returns (bytes32) {
        return bytes32(0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f);
    }

    function _hash7() internal pure returns (bytes32) {
        return bytes32(0x707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f);
    }

    // ─── encodeVarint tests ───

    function test_encodeVarint_zero() public pure {
        assertEq(Encode.encodeVarint(0), hex"00");
    }

    function test_encodeVarint_one() public pure {
        assertEq(Encode.encodeVarint(1), hex"01");
    }

    function test_encodeVarint_127() public pure {
        assertEq(Encode.encodeVarint(127), hex"7f");
    }

    function test_encodeVarint_128() public pure {
        assertEq(Encode.encodeVarint(128), hex"8001");
    }

    function test_encodeVarint_255() public pure {
        assertEq(Encode.encodeVarint(255), hex"ff01");
    }

    function test_encodeVarint_256() public pure {
        assertEq(Encode.encodeVarint(256), hex"8002");
    }

    function test_encodeVarint_300() public pure {
        assertEq(Encode.encodeVarint(300), hex"ac02");
    }

    function test_encodeVarint_16384() public pure {
        assertEq(Encode.encodeVarint(16384), hex"808001");
    }

    function test_encodeVarint_1000000() public pure {
        assertEq(Encode.encodeVarint(1000000), hex"c0843d");
    }

    // ─── encodeString tests ───

    function test_encodeString_empty() public pure {
        assertEq(Encode.encodeString(""), hex"00");
    }

    function test_encodeString_cosmoshub4() public pure {
        assertEq(Encode.encodeString("cosmoshub-4"), hex"0b636f736d6f736875622d34");
    }

    function test_encodeString_testChain() public pure {
        assertEq(Encode.encodeString("test-chain"), hex"0a746573742d636861696e");
    }

    function test_encodeString_singleChar() public pure {
        assertEq(Encode.encodeString("a"), hex"0161");
    }

    // ─── encodeValidator tests ───
    // Go reference: proto.Marshal(SimpleValidator{PubKey: &PublicKey{Ed25519: pubKey}, VotingPower: vp})

    function test_encodeValidator_power100() public pure {
        IICS07TendermintMsgs.SimpleValidator memory v = IICS07TendermintMsgs.SimpleValidator({
            pubKey: _pubKey(),
            votingPower: 100
        });
        assertEq(
            Encode.encodeValidator(v),
            hex"0a220a20aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c91064"
        );
    }

    function test_encodeValidator_power1() public pure {
        IICS07TendermintMsgs.SimpleValidator memory v = IICS07TendermintMsgs.SimpleValidator({
            pubKey: _pubKey(),
            votingPower: 1
        });
        assertEq(
            Encode.encodeValidator(v),
            hex"0a220a20aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c91001"
        );
    }

    function test_encodeValidator_power0() public pure {
        IICS07TendermintMsgs.SimpleValidator memory v = IICS07TendermintMsgs.SimpleValidator({
            pubKey: _pubKey(),
            votingPower: 0
        });
        // Go proto.Marshal skips votingPower=0
        assertEq(
            Encode.encodeValidator(v),
            hex"0a220a20aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9"
        );
    }

    function test_encodeValidator_power1000000() public pure {
        IICS07TendermintMsgs.SimpleValidator memory v = IICS07TendermintMsgs.SimpleValidator({
            pubKey: _pubKey(),
            votingPower: 1000000
        });
        assertEq(
            Encode.encodeValidator(v),
            hex"0a220a20aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c910c0843d"
        );
    }

    // ─── encodeVersion tests ───

    function test_encodeVersion_11_0() public pure {
        IICS07TendermintMsgs.Version memory v = IICS07TendermintMsgs.Version({ blockVersion: 11, appVersion: 0 });
        // Go proto.Marshal: appVersion=0 is skipped
        assertEq(Encode.encodeVersion(v), hex"080b");
    }

    function test_encodeVersion_11_2() public pure {
        IICS07TendermintMsgs.Version memory v = IICS07TendermintMsgs.Version({ blockVersion: 11, appVersion: 2 });
        assertEq(Encode.encodeVersion(v), hex"080b1002");
    }

    function test_encodeVersion_1_1() public pure {
        IICS07TendermintMsgs.Version memory v = IICS07TendermintMsgs.Version({ blockVersion: 1, appVersion: 1 });
        assertEq(Encode.encodeVersion(v), hex"08011001");
    }

    // ─── encodePartSetHeader tests ───

    function test_encodePartSetHeader() public pure {
        IICS07TendermintMsgs.PartSetHeader memory psh = IICS07TendermintMsgs.PartSetHeader({
            total: 1,
            hashData: _hash1()
        });
        assertEq(
            Encode.encodePartSetHeader(psh),
            hex"08011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"
        );
    }

    // ─── encodeBlockId tests ───

    function test_encodeBlockId() public pure {
        IICS07TendermintMsgs.BlockId memory bid = IICS07TendermintMsgs.BlockId({
            hashData: _hash2(),
            partSetHeader: IICS07TendermintMsgs.PartSetHeader({ total: 1, hashData: _hash1() })
        });
        assertEq(
            Encode.encodeBlockId(bid),
            hex"0a20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f122408011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"
        );
    }

    // ─── cdcEncode tests ───

    function test_cdcEncodeString() public pure {
        // Go: cdcEncode("cosmoshub-4") = StringValue{Value: "cosmoshub-4"}.Marshal()
        assertEq(Encode.cdcEncodeString("cosmoshub-4"), hex"0a0b636f736d6f736875622d34");
    }

    function test_cdcEncodeString_empty() public pure {
        assertEq(Encode.cdcEncodeString(""), hex"");
    }

    function test_cdcEncodeInt64() public pure {
        // Go: cdcEncode(int64(12345)) = Int64Value{Value: 12345}.Marshal()
        assertEq(Encode.cdcEncodeInt64(12345), hex"08b960");
    }

    function test_cdcEncodeInt64_zero() public pure {
        assertEq(Encode.cdcEncodeInt64(0), hex"");
    }

    function test_cdcEncodeBytes32() public pure {
        // Go: cdcEncode(hash5) = BytesValue{Value: hash5}.Marshal()
        assertEq(
            Encode.cdcEncodeBytes32(_hash5()),
            hex"0a20505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f"
        );
    }

    function test_encodeTimestamp() public pure {
        // Go: gogotypes.StdTimeMarshal(time.Unix(1700000000, 0))
        assertEq(Encode.encodeTimestamp(1700000000), hex"0880e2cfaa06");
    }

    // ─── hashValSet tests ───

    function test_hashValSet() public pure {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](2);
        vals[0] = IICS07TendermintMsgs.ValidatorInfo({
            valAddress: hex"",
            pubKey: _pubKey(),
            votingPower: 100,
            proposerPriority: 0
        });
        vals[1] = IICS07TendermintMsgs.ValidatorInfo({
            valAddress: hex"",
            pubKey: bytes32(0xbbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9da),
            votingPower: 200,
            proposerPriority: 0
        });

        IICS07TendermintMsgs.ValidatorSet memory valSet = IICS07TendermintMsgs.ValidatorSet({
            validators: vals,
            hasProposer: false,
            proposer: IICS07TendermintMsgs.ValidatorInfo({
                valAddress: hex"",
                pubKey: bytes32(0),
                votingPower: 0,
                proposerPriority: 0
            }),
            totalVotingPower: 300
        });

        // Go reference: merkle.HashFromByteSlices(proto-encoded validators)
        assertEq(
            Header.hashValSet(valSet),
            bytes32(0x1bf78c35d508c8720ea6cca7225cf9fe1a628438d1aa6c7e2d2dbddff036023f)
        );
    }

    // ─── hashHeader test ───

    function test_hashHeader() public pure {
        IICS07TendermintMsgs.BlockHeader memory header = IICS07TendermintMsgs.BlockHeader({
            version: IICS07TendermintMsgs.Version({ blockVersion: 11, appVersion: 0 }),
            chainId: "cosmoshub-4",
            height: 12345,
            time: 1700000000,
            hasLastBlockId: true,
            lastBlockId: IICS07TendermintMsgs.BlockId({
                hashData: _hash2(),
                partSetHeader: IICS07TendermintMsgs.PartSetHeader({ total: 1, hashData: _hash1() })
            }),
            hasLastCommitHash: true,
            lastCommitHash: _hash3(),
            hasDataHash: true,
            dataHash: _hash4(),
            validatorsHash: _hash5(),
            nextValidatorsHash: _hash6(),
            consensusHash: _hash7(),
            appHash: bytes32(0xa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf),
            hasLastResultsHash: true,
            lastResultsHash: _hash1(),
            hasEvidenceHash: true,
            evidenceHash: _hash2(),
            proposerAddress: hex"f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff00010203"
        });

        // Go reference: types.Header.Hash()
        assertEq(
            Header.hashHeader(header),
            bytes32(0x94db05d2ad72b0ad975ce0b1a832c5cd4b7bfdc40008a93a815509b12cfebc20)
        );
    }

    // ─── merkleHash tests ───

    function test_merkleHash_empty() public pure {
        bytes[] memory items = new bytes[](0);
        assertEq(Header.merkleHash(items), bytes32(0));
    }

    function test_merkleHash_singleLeaf() public pure {
        bytes[] memory items = new bytes[](1);
        items[0] = hex"deadbeef";
        // Go: merkle.HashFromByteSlices (1-byte 0x00 leaf prefix)
        assertEq(
            Header.merkleHash(items),
            bytes32(0x48c90c8ae24688d6bef5d48a30c2cc8b6754335a8db21793cc0a8e3bed321729)
        );
    }

    function test_merkleHash_twoLeaves() public pure {
        bytes[] memory items = new bytes[](2);
        items[0] = hex"aa";
        items[1] = hex"bb";
        // Go: merkle.HashFromByteSlices (1-byte prefixes)
        assertEq(
            Header.merkleHash(items),
            bytes32(0x3a6c27fde711243b24095e21cfa3ab2fd6c4e186412e526b2684850209247eca)
        );
    }

    // ─── Debug: emit encoded bytes for manual inspection ───

    function test_debug_printEncodings() public pure {
        console.log("=== encodeVarint ===");
        console.logBytes(Encode.encodeVarint(0));
        console.logBytes(Encode.encodeVarint(128));
        console.logBytes(Encode.encodeVarint(300));
        console.logBytes(Encode.encodeVarint(1000000));

        console.log("=== encodeVersion ===");
        console.logBytes(
            Encode.encodeVersion(IICS07TendermintMsgs.Version({ blockVersion: 11, appVersion: 0 }))
        );

        console.log("=== encodeValidator ===");
        console.logBytes(
            Encode.encodeValidator(
                IICS07TendermintMsgs.SimpleValidator({ pubKey: _pubKey(), votingPower: 100 })
            )
        );

        console.log("=== encodeBlockId ===");
        console.logBytes(
            Encode.encodeBlockId(
                IICS07TendermintMsgs.BlockId({
                    hashData: _hash2(),
                    partSetHeader: IICS07TendermintMsgs.PartSetHeader({ total: 1, hashData: _hash1() })
                })
            )
        );
    }
}
