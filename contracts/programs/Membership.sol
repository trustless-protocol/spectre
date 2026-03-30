// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { IMembershipMsgs } from "../light-clients/msgs/IMembershipMsgs.sol";
import { IMembership } from "../interfaces/IMembership.sol";

import "@openzeppelin-contracts/utils/math/Math.sol";
import "@openzeppelin-contracts/utils/Bytes.sol";
/**
 * @title MerkleTreeMembership
 * @dev Contract to verify membership of key-value pairs in a Merkle tree
 * Converted from Rust zkVM code for Cosmos SDK proof verification
 */
contract Membership  is IMembership {
    using Math for uint256;
    using Bytes for *;
     
    // Events
    event MembershipVerified(bytes32 indexed commitmentRoot, uint256 pairsCount);
    event NonMembershipVerified(bytes32 indexed commitmentRoot, bytes[] key);
    
    // Custom errors
    error InvalidLength();
    error EmptyRequest();
    error VerificationMembershipFailed();
    error VerificationNonMembershipFailed();
    error MissingMerkleProof();
    error MissingMerkleRoot();
    error MismatchedNumberOfProofs(uint256 expected, uint256 actual);
    error MissingVerifiedValue();
    error InvalidMerkleProof();
    error InvalidExistenceProof();
    error FailedToVerifyMembership();
    error InputDataMissing();
    error ProvidedKeyValueMismatch();
    error RootMismatch();
    error MissingLeafSpec();
    error InvalidVarint();
    error InvalidOffset();
    error BranchNotFound(uint256 branch);

    /**
     * @dev Verify membership of multiple key-value pairs in the Merkle tree
     * @param appHash The root hash of the Merkle tree (32 bytes)
     * @param kvPairs Array of key-value pairs to verify
     * @param merkleProofs Array of corresponding Merkle proofs
     * @return output The membership verification result
     */
    function membership(
        bytes32 appHash,
        IMembershipMsgs.KVPair[] calldata kvPairs,
        IMembershipMsgs.MerkleProof[] calldata merkleProofs
    ) public returns (IMembershipMsgs.MembershipOutput memory output) {
        if (kvPairs.length == 0) {
            revert EmptyRequest();
        }
        
        if (kvPairs.length != merkleProofs.length) {
            revert InvalidLength();
        }
        
        bytes32 commitmentRoot = appHash;
        IMembershipMsgs.KVPair[] memory verifiedPairs = new IMembershipMsgs.KVPair[](kvPairs.length);
        
        for (uint256 i = 0; i < kvPairs.length; i++) {
            IMembershipMsgs.KVPair memory kvPair = kvPairs[i];
            IMembershipMsgs.MerkleProof memory merkleProof = merkleProofs[i];

            IMembershipMsgs.ProofSpec[] memory proofSpecs = new IMembershipMsgs.ProofSpec[](2);
            proofSpecs[0] = iavlSpec();
            proofSpecs[1] = tendermintSpec();
            // Check if this is a non-membership proof (empty value)
            if (kvPair.value.length == 0) {
                // Verify non-membership
                if (!verifyNonMembership(proofSpecs, commitmentRoot, kvPair.path, merkleProof)) {
                    revert VerificationNonMembershipFailed();
                }
                emit NonMembershipVerified(commitmentRoot, kvPair.path);
            } else {
                // Verify membership
                verifyMembership(proofSpecs, commitmentRoot, kvPair.path, kvPair.value, 0, merkleProof);
            }
            
            verifiedPairs[i] = kvPair;
        }

        output = IMembershipMsgs.MembershipOutput({
            commitmentRoot: commitmentRoot,
            kvPairs: verifiedPairs
        });
        
        emit MembershipVerified(commitmentRoot, kvPairs.length);
        return output;
    }
    
    function verifyMembership(
        IMembershipMsgs.ProofSpec[] memory proofSpecs,
        bytes32 root,
        bytes[] memory path,
        bytes32 value,
        uint256 startIndex,
        IMembershipMsgs.MerkleProof memory proof
    ) internal view {
        if (proof.proofs.length == 0) {
            revert MissingMerkleProof();
        }

        if (root == bytes32(0)) {
            revert MissingMerkleRoot();
        }

        uint256 proofLength = proof.proofs.length;
        if (proofSpecs.length != proofLength) {
            revert MismatchedNumberOfProofs(proofSpecs.length, proofLength);
        }

        if (path.length != proofLength) {
            revert MismatchedNumberOfProofs(path.length, proofLength);
        }

        if (value.length == 0) {
            revert MissingVerifiedValue();
        }

        // Process proofs from startIndex onwards
        // Keys are represented from root-to-leaf, so we iterate in reverse
        bytes32 subroot = value;
        bytes32 valueUpdate = value;

        uint256 pathLength = path.length;
        for (uint256 i = startIndex; i < proofLength; i++) {
            bytes memory keyPath = path[pathLength - i - 1];

            IMembershipMsgs.CommitmentProof memory commitmentProof = proof.proofs[i];
            if (commitmentProof.proofType != IMembershipMsgs.ProofType.EXIST) {
                revert InvalidMerkleProof();
            }

            subroot = calculateExistenceRoot(commitmentProof.existenceProof);
            if (!verifyExistenceProof(
                commitmentProof.existenceProof,
                proofSpecs[i],
                subroot,
                keyPath,
                valueUpdate
                )
            ) {
                revert FailedToVerifyMembership();
            }
            valueUpdate = subroot;
        }

        if (!(keccak256(abi.encode(root)) == keccak256(abi.encode(subroot)))) {
            revert FailedToVerifyMembership();
        }
    }

    /**
     * @dev Verify non-membership of a key in the Merkle tree
     * @param proofSpecs Array of proof specifications
     * @param root The Merkle root
     * @param path The key to verify non-existence
     * @param proof The Merkle proof for non-membership
     * @return True if the non-membership proof is valid
     */
    function verifyNonMembership(
        IMembershipMsgs.ProofSpec[] memory proofSpecs,
        bytes32 root,
        bytes[] memory path,
        IMembershipMsgs.MerkleProof memory proof
    ) internal view returns (bool) {
        if (proof.proofs.length == 0) {
            revert MissingMerkleProof();
        }

        if (root == bytes32(0)) {
            revert MissingMerkleRoot();
        }

        uint256 proofLength = proof.proofs.length;
        if (proofSpecs.length != proofLength) {
            revert MismatchedNumberOfProofs(proofSpecs.length, proofLength);
        }

        if (path.length != proofLength) {
            revert MismatchedNumberOfProofs(path.length, proofLength);
        }

        // verify the absence of key in lowest subtree
        IMembershipMsgs.CommitmentProof memory firstProof = proof.proofs[0];
        IMembershipMsgs.ProofSpec memory firstSpec = proofSpecs[0];

        // keys are represented from root-to-leaf        
        bytes memory key = path[proofLength-1];

        if (firstProof.proofType != IMembershipMsgs.ProofType.NON_EXIST) {
            revert InvalidMerkleProof();
        }

        IMembershipMsgs.NonExistenceProof memory nonExistenceProof = firstProof.nonExistenceProof;
        bytes32 subroot = calculateNonExistenceRoot(nonExistenceProof);

        if (!verifyNonExistenceProof(nonExistenceProof, firstSpec, subroot, key)) {
            revert FailedToVerifyMembership();
        }
        // verify membership proofs starting from index 1 with value = subroot
        verifyMembership(
            proofSpecs,
            root,
            path,
            subroot,
            1,
            proof
        );
    }

    function calculateExistenceRoot(IMembershipMsgs.ExistenceProof memory proof) 
        internal 
        pure 
        returns (bytes32) 
    {
        if (proof.key.length == 0 || proof.value.length == 0) {
            revert InvalidExistenceProof();
        }

        IMembershipMsgs.LeafOp memory leafOp = proof.leaf;
        bytes32 current = applyLeaf(leafOp, proof.key, proof.value);
        for (uint256 i = 0; i < proof.path.length; i++) {
            current = applyInner(proof.path[i], current);
            // TODO
        }
        
        return current;
    }

    function calculateNonExistenceRoot(IMembershipMsgs.NonExistenceProof memory proof) 
        internal 
        pure 
        returns (bytes32) 
    {
        if (proof.hasLeft) {
            return calculateExistenceRoot(proof.left);
        } else if (proof.hasRight) {
            return calculateExistenceRoot(proof.right);
        } else {
            revert InvalidMerkleProof();
        }
    }

    function checkExistenceProof(
        IMembershipMsgs.ExistenceProof memory proof,
        IMembershipMsgs.ProofSpec memory spec
    ) internal view {
        if (!spec.hasLeafSpec) {
            revert MissingLeafSpec();
        }

        bytes memory leafPrefix = proof.leaf.prefix;
        // ensure leaf prefix matches the spec
        if (spec.specType == IMembershipMsgs.SpecType.IAVL) {
            uint256 offset = 0;
            
            uint256 remainingLength = ensureIavlPrefix(leafPrefix, 0);
            if (remainingLength != 0) {
                revert("bad prefix in leaf");
            }
        }

        //  ensure leaf hash matches the spec
        IMembershipMsgs.LeafOp memory leaf = proof.leaf;
        if (spec.leafOp.hashOp != leaf.hashOp) {
            revert("Unexpected leaf hash operation");
        }
        if (spec.leafOp.prehashKey != leaf.prehashKey) {
            revert("Unexpected leaf prehash key operation");
        }
        if (spec.leafOp.prehashValue != leaf.prehashValue) {
            revert("Unexpected leaf prehash value operation");
        }
        bytes memory leafSpecPrefix = spec.leafOp.prefix;
        if (leafSpecPrefix.length > leafPrefix.length ||
            !(keccak256(abi.encode(leafSpecPrefix)) == keccak256(abi.encode(getSlice(leafPrefix, 0, leafSpecPrefix.length))))) {
            revert("Incorrect prefix on leaf");
        }

        // ensure min/max depths
        if (spec.minDepth != 0) {
            if (proof.path.length < uint256(spec.minDepth)) {
                revert("Too few InnerOps");
            }
            if (proof.path.length > uint256(spec.maxDepth)) {
                revert("Too many InnerOps");
            }
        }

        if (spec.specType == IMembershipMsgs.SpecType.IAVL) {
            uint256 stepLength = proof.path.length;
            for (uint256 i = 0; i < stepLength; i++) {
                uint256 remainingLength = ensureIavlPrefix(proof.path[i].prefix, 0);
                IMembershipMsgs.InnerOp memory step = proof.path[i];
                if (remainingLength != 0) {
                    // 1 byte due to containing length prefix for left hash.
                    // 33 bytes due to IAVL length prefix + left hash + next IAVL legnth prefix
                    if (remainingLength != 1 && remainingLength != 34) {
                        revert("bad prefix in inner op");
                    }
                    if (step.hashOp != IMembershipMsgs.HashOp.SHA256) {
                        revert("bad hash operation");
                    }
                }

                if (!spec.hasInnerSpec) {
                    revert("InnerSpec is required");
                }
                IMembershipMsgs.InnerOp memory innerOp = proof.path[i];
                if (spec.innerSpec.hashOp != innerOp.hashOp) {
                    revert("Unexpected inner hash operation");
                }

                // NOTE: removed - inner ops should NOT be required to start with leaf spec prefix
                // if (leafSpecPrefix.length > innerOp.prefix.length || !(keccak256(abi.encode(leafSpecPrefix)) == keccak256(abi.encode(getSlice(innerOp.prefix, 0, leafSpecPrefix.length))))) {
                //     revert("Incorrect prefix on leaf");
                // }

                if (innerOp.prefix.length < spec.innerSpec.minPrefixLength) {
                    revert("Inner prefix too short");
                }

                uint32 maxLeftChild = uint32(spec.innerSpec.childOrder.length - 1) * (spec.innerSpec.childSize);
                if (innerOp.prefix.length > maxLeftChild + spec.innerSpec.maxPrefixLength) {
                    revert("Inner prefix too long");
                }

                if (spec.innerSpec.childSize == 0) {
                    revert("Inner child size must >=1");
                }

                if (innerOp.suffix.length % spec.innerSpec.childSize != 0) {
                    revert("Inner suffix malformed");
                }
            }
        }
    }

    function calculateExistenceRootForSpec(
        IMembershipMsgs.ExistenceProof memory proof,
        IMembershipMsgs.ProofSpec memory spec
    ) internal pure returns (bytes32) {
        if (proof.key.length == 0 || proof.value.length == 0) {
            revert InvalidExistenceProof();
        }
        IMembershipMsgs.LeafOp memory leafOp = proof.leaf;
        bytes32 leafHash = applyLeaf(leafOp, proof.key, proof.value);

        for (uint256 i = 0; i < proof.path.length; i++) {
            leafHash = applyInner(proof.path[i], leafHash);

            if (spec.hasInnerSpec) {
                if (leafHash.length > uint256(spec.innerSpec.childSize) && 
                    spec.innerSpec.childSize >= 32) {
                    revert("Invalid inner operation (child_size)");
                }
            }
        }

        return leafHash;
    }



    function verifyExistenceProof(
        IMembershipMsgs.ExistenceProof memory proof,
        IMembershipMsgs.ProofSpec memory spec,
        bytes32 subroot,
        bytes memory key,
        bytes32 value
    ) internal view returns (bool) {
        checkExistenceProof(proof, spec);
        if (keccak256(abi.encode(proof.key)) != keccak256(abi.encode(key)) || proof.value != value) {
            revert ProvidedKeyValueMismatch();
        }

        bytes32 calculateRoot = calculateExistenceRootForSpec(proof, spec);
        if (keccak256(abi.encode(calculateRoot)) != keccak256(abi.encode(subroot))) {
            revert RootMismatch();
        }
        return true;
    }
    
    function verifyNonExistenceProof(
        IMembershipMsgs.NonExistenceProof memory proof,
        IMembershipMsgs.ProofSpec memory spec,
        bytes32 root,
        bytes memory key
    ) internal view returns (bool) {

        bool preHash = spec.prehashKeyBeforeComparison;
        IMembershipMsgs.HashOp prehashOp = spec.leafOp.prehashKey;
        if (proof.hasLeft) {
            verifyExistenceProof(proof.left, spec, root, proof.left.key, proof.left.value);
            require(compareBytes(
                keyForComparison(key, preHash, prehashOp),
                keyForComparison(proof.left.key, preHash, prehashOp)
            ) == 1, "left key isn't before key");
        }

        if (proof.hasRight) {
            verifyExistenceProof(proof.right, spec, root, proof.right.key, proof.right.value);
            require(compareBytes(
                keyForComparison(key, preHash, prehashOp),
                keyForComparison(proof.right.key, preHash, prehashOp)
            ) == -1, "right key isn't after key");
        }

        if (!spec.hasInnerSpec) {
            revert("InnerSpec is missing");
        }
        IMembershipMsgs.InnerSpec memory innerSpec = spec.innerSpec;

        if (!proof.hasLeft && !proof.hasRight) {
            revert("neither left nor right proof defined");
        } else if (!proof.hasLeft && proof.hasRight) {
            ensureLeftMost(innerSpec, proof.left.path);
        } else if (proof.hasLeft && !proof.hasRight) {
            ensureRightMost(innerSpec, proof.right.path);
        } else if (proof.hasLeft && proof.hasRight) {

            uint256 leftIndex = proof.left.path.length -1;
            uint256 rightIndex = proof.right.path.length - 1;

            IMembershipMsgs.InnerOp memory topLeft = proof.left.path[leftIndex];
            IMembershipMsgs.InnerOp memory topRight = proof.right.path[rightIndex];
            while (leftIndex > 0 && rightIndex > 0
                && keccak256(abi.encode(topLeft.prefix)) == keccak256(abi.encode(topRight.prefix))
                && keccak256(abi.encode(topLeft.suffix)) == keccak256(abi.encode(topRight.suffix))
            ) {
                leftIndex--;
                rightIndex--;

                if (leftIndex < 0 || rightIndex < 0) {
                    revert("Invalid non-existence proof");
                }

                topLeft = proof.left.path[leftIndex];
                topRight = proof.right.path[rightIndex];
            }

            uint256 leftPaddingIdx = orderFromPadding(innerSpec, topLeft);
            uint256 rightPaddingIdx = orderFromPadding(innerSpec, topRight); 

            if (!(leftPaddingIdx + 1 == rightPaddingIdx)) {
                revert("Not left neighbor at first divergent step");
            }

            ensureLeftMost(innerSpec, proof.left.path);
            ensureRightMost(innerSpec, proof.right.path);
        } else {
            revert("Invalid non-existence proof");
        }
        return true;
    }
    /**
     * @dev Generic Merkle proof verification
     * @param root The Merkle root
     * @param leaf The leaf hash to verify
     * @param proof Array of sibling hashes for the proof path
     * @param index The index of the leaf in the tree
     * @return True if the proof is valid
     */
    function verifyProof(
        bytes32 root,
        bytes32 leaf,
        bytes[] memory proof,
        uint256 index
    ) internal pure returns (bool) {
        
        bytes32 computedHash = leaf;
        uint256 currentIndex = index;
        
        for (uint256 i = 0; i < proof.length; i++) {
            bytes32 proofElement = bytesToBytes32(proof[i]);
            
            if (currentIndex % 2 == 0) {
                // If current index is even, proof element is right sibling
                computedHash = keccak256(abi.encodePacked(computedHash, proofElement));
            } else {
                // If current index is odd, proof element is left sibling
                computedHash = keccak256(abi.encodePacked(proofElement, computedHash));
            }
            
            currentIndex = currentIndex / 2;
        }
        
        return computedHash == root;
    }

    function applyLeaf(
        IMembershipMsgs.LeafOp memory leafOp,
        bytes memory key,
        bytes32 value
    ) internal pure returns (bytes32) {
        bytes memory hashedData = leafOp.prefix;

        bytes memory prekey = prepareLeafData(leafOp.prehashKey, key);
        hashedData = abi.encodePacked(hashedData, prekey);

        bytes memory valueBytes = bytes32ToBytes(value);
        bytes memory preval = prepareLeafData(leafOp.prehashValue, valueBytes);
        hashedData = abi.encodePacked(hashedData, preval);

        return hashData(hashedData, leafOp.hashOp);
    }

    function applyInner(
        IMembershipMsgs.InnerOp memory inner,
        bytes32 child
    ) internal pure returns (bytes32) {
        if (child == bytes32(0)) {
            revert("missing child hash");
        }

        bytes memory image = abi.encodePacked(inner.prefix, child, inner.suffix);
        return hashData(image, inner.hashOp);
    }

    function prepareLeafData(
        IMembershipMsgs.HashOp prehashOp,
        bytes memory data
    ) internal pure returns (bytes memory) {
        if (data.length == 0) {
            revert InputDataMissing();
        }

        if (prehashOp == IMembershipMsgs.HashOp.NO_HASH) {
            bytes memory encodedLen = encodeVarint(data.length);
            return abi.encodePacked(encodedLen, data);
        }

        bytes32 hashedData = hashData(data, prehashOp);
        bytes memory encodedLength = encodeVarint(uint256(32));
        return abi.encodePacked(encodedLength, hashedData);
    }

    // true if this is the right-most path in the tree, excluding placeholder (empty child) nodes
    function ensureRightMost(
        IMembershipMsgs.InnerSpec memory innerSpec,
        IMembershipMsgs.InnerOp[] memory path
    ) internal view {
        IMembershipMsgs.Padding memory padding = getPadding(innerSpec, innerSpec.childOrder.length - 1);

        for (uint256 i = 0; i < path.length; i++) {
            IMembershipMsgs.InnerOp memory innerOp = path[i];
            bool rightHasPadding = hasPadding(innerOp, padding);
            uint256 rightBranches = innerSpec.childOrder.length - 1 - orderFromPadding(innerSpec, innerOp);

            bool isEmpty = true;
            uint256 childSize = uint256(innerSpec.childSize);
            if (rightBranches == 0 || innerOp.suffix.length != childSize) {
                isEmpty = false;
            } else {
                // compare prefix with the expected number of empty branches
                for (uint256 j = 0; j < rightBranches; j++) {
                    bool found = false;
                    uint256 idx;
                    for (uint256 k = 0; k < innerSpec.childOrder.length; k++) {
                        if (innerSpec.childOrder[k] == j) {
                            found = true;
                            idx = k;
                            break;
                        }
                    }

                    if (!found) {
                        isEmpty = false;
                        break;
                    }

                    uint256 from = idx * childSize;
                    if (keccak256(abi.encode(innerSpec.emptyChild)) != keccak256(abi.encode(getSlice(innerOp.suffix, from, from + childSize)))) {
                        isEmpty = false;
                        break;
                    }
                }
            }
            if (!rightHasPadding && !isEmpty) {
                revert("right padding mismatch");
            }
        }
    }

    function ensureLeftMost(
        IMembershipMsgs.InnerSpec memory innerSpec,
        IMembershipMsgs.InnerOp[] memory path
    ) internal view {
        // fails unless this is the left-most path in the tree, excluding placeholder (empty child) nodes
        IMembershipMsgs.Padding memory padding = getPadding(innerSpec, 0);
        for (uint256 i = 0; i < path.length; i++) {
            IMembershipMsgs.InnerOp memory innerOp = path[i];
            bool leftHasPadding = hasPadding(innerOp, padding);
            uint256 leftBranches = orderFromPadding(innerSpec, innerOp);

            bool isEmpty = true;
            if (leftBranches == 0) {
                isEmpty = false;
            } else {
                // compare prefix with the expected number of empty branches
                uint256 childSize = uint256(innerSpec.childSize); 
                (bool subSuccess, uint256 actualPrefix) = Math.trySub(innerOp.prefix.length, childSize * leftBranches);
                if (!subSuccess) {
                    isEmpty = false;
                } else {
                    for (uint256 j = 0; j < leftBranches; j++) {
                        bool found = false;
                        uint256 idx;
                        for (uint256 k = 0; k < innerSpec.childOrder.length; k++) {
                            if (innerSpec.childOrder[k] == j) {
                                found = true;
                                idx = k;
                                break;
                            }
                        }

                        if (!found) {
                            isEmpty = false;
                            break;
                        }

                        uint256 from = actualPrefix + idx * childSize;
                        if (keccak256(abi.encode(innerSpec.emptyChild)) != keccak256(abi.encode(getSlice(innerOp.prefix, from, from + childSize)))) {
                            isEmpty = false;
                            break;
                        }
                    }
                }
            }
            if (!leftHasPadding && !isEmpty) {
                revert("left padding mismatch");
            }
        }
    }

    function orderFromPadding(
        IMembershipMsgs.InnerSpec memory innerSpec,
        IMembershipMsgs.InnerOp memory innerOp
    ) internal pure returns (uint256) {
        uint256 childOrderLength = innerSpec.childOrder.length;
        for (uint256 branch = 0; branch < childOrderLength; branch++) {
            IMembershipMsgs.Padding memory padding = getPadding(innerSpec, uint32(branch));
            if (hasPadding(innerOp, padding)) {
                return branch;
            }
        }
        revert("padding doesn't match any branch");
    }
    
    function getPadding(
        IMembershipMsgs.InnerSpec memory innerSpec,
        uint256 branch
    ) internal pure returns (IMembershipMsgs.Padding memory) {
        uint256 foundIdx = 0;
        bool found = false;
        for (uint256 i = 0; i < innerSpec.childOrder.length; i++) {
            if (uint32(innerSpec.childOrder[i]) == branch) {
                foundIdx = i;
                found = true;
                break;
            }
        }
        
        // If branch not found, revert with error
        if (!found) {
            revert BranchNotFound(branch);
        }
    
        uint32 idx = uint32(foundIdx);
        uint32 prefix = idx * innerSpec.childSize;
        uint256 suffix = uint256(innerSpec.childSize) * (innerSpec.childOrder.length - 1 - idx);
        return IMembershipMsgs.Padding({
            minPrefix: prefix + innerSpec.minPrefixLength,
            maxPrefix: prefix + innerSpec.maxPrefixLength,
            suffix: suffix
        });
    }

    function hasPadding(
        IMembershipMsgs.InnerOp memory inner,
        IMembershipMsgs.Padding memory padding
    ) internal pure returns (bool) {
        return (
            inner.prefix.length >= padding.minPrefix &&
            inner.prefix.length <= padding.maxPrefix &&
            inner.suffix.length == padding.suffix
        );
    }

     function compareBytes(bytes memory a, bytes memory b) public pure returns (int8) {
        if (a.length != b.length) {
            // For different lengths, we still need to compare byte by byte
            // up to the shorter length, then compare lengths
            uint256 minLength = a.length < b.length ? a.length : b.length;
            
            for (uint256 i = 0; i < minLength; i++) {
                if (a[i] < b[i]) return -1;
                if (a[i] > b[i]) return 1;
            }
            
            return a.length < b.length ? int8(-1) : int8(1);
        }
        
        // Same length - we can use keccak256 for equality check first
        if (keccak256(a) == keccak256(b)) {
            return 0;
        }
        
        // Different content, same length - compare byte by byte
        for (uint256 i = 0; i < a.length; i++) {
            if (a[i] < b[i]) return -1;
            if (a[i] > b[i]) return 1;
        }
        
        return 0; // Should never reach here
    }

    /**
     * @dev Convert bytes to bytes32
     * @param data The bytes to convert
     * @return result The bytes32 result
     */
    function bytesToBytes32(bytes memory data) internal pure returns (bytes32 result) {
        if (data.length >= 32) {
            assembly ("memory-safe") {
                result := mload(add(data, 32))
            }
        } else {
            // Pad with zeros if data is shorter than 32 bytes
            bytes32 temp;
            assembly ("memory-safe") {
                temp := mload(add(data, 32))
            }
            result = temp >> (8 * (32 - data.length));
        }
    }

    function bytes32ToBytes(bytes32 data) internal pure returns (bytes memory) {
        bytes memory result = new bytes(32);
        assembly ("memory-safe") {
            mstore(add(result, 32), data)
        }
        return result;
    }

    function getSlice(bytes memory array, uint256 from, uint256 to) 
        public 
        view 
        returns (bytes memory) 
    {
        require(from <= to, "Invalid range: from > to");
        require(to <= array.length, "Range exceeds array length");
        
        uint256 length = to - from;
        bytes memory result = new bytes(length);
        
        assembly ("memory-safe") {
            let src := add(add(array, 0x20), from)
            let dest := add(result, 0x20)
            
            // Use identity precompile for efficient copying
            let success := staticcall(gas(), 0x04, src, length, dest, length)
        }
        
        return result;
    }

    function hashData(bytes memory data, IMembershipMsgs.HashOp hashOp) internal pure returns (bytes32) {
        if (hashOp == IMembershipMsgs.HashOp.NO_HASH) {
            return bytesToBytes32(data);
        } else if (hashOp == IMembershipMsgs.HashOp.SHA256) {
            return sha256(data);
        } else if (hashOp == IMembershipMsgs.HashOp.KECCAK256) {
            return keccak256(data);
        } else {
            revert("Unsupported hash operation");
        }
    }

    function keyForComparison(
        bytes memory key,
        bool prehash, 
        IMembershipMsgs.HashOp prehashOp
    ) internal pure returns (bytes memory) {
        if (prehash) {
            return bytes32ToBytes(hashData(key, prehashOp));
        } else {
            return key;
        }
    }

    function readVarint(bytes memory data, uint256 offset) internal pure returns (int64, uint256) {
        (uint64 ux, uint256 newOffset) = decodeVarint(data, offset);
        int64 x = int64(ux >> 1);
        if (ux & 1 != 0) {
            x = ~x;
        }
        return (x, newOffset);
    }
    
    function encodeVarint(uint256 value) internal pure returns (bytes memory) {
        if (value < 128) {
            return abi.encodePacked(uint8(value));
        }
        
        bytes memory result;
        while (value >= 128) {
            result = abi.encodePacked(result, uint8((value & 0x7F) | 0x80));
            value >>= 7;
        }
        result = abi.encodePacked(result, uint8(value));
        return result;
    }

    function decodeVarint(bytes memory data, uint256 offset) internal pure returns (uint64, uint256) {
        if (offset >= data.length) {
            revert InvalidOffset();
        }

        if (data.length - offset == 0) {
            revert InvalidVarint();
        }

        uint8 firstByte = uint8(data[offset]);
        if (firstByte < 0x80) {
            return (uint64(firstByte), offset + 1);
        } else if (data.length - offset > 10 || data[data.length - 1] < 0x80) {
            uint32 part0 = uint32(firstByte) - 0x80;
            uint8 b = uint8(data[offset + 1]);
            part0 += (uint32(b) << 7);
            if (b < 0x80) {
                return (part0, offset + 2);
            }
            part0 -= 0x80 << 7;

            b = uint8(data[offset + 2]);
            part0 += (uint32(b) << 14);
            if (b < 0x80) {
                return (part0, offset + 3);
            }
            part0 -= 0x80 << 14;

            b = uint8(data[offset + 3]);
            part0 += (uint32(b) << 21);
            if (b < 0x80) {
                return (part0, offset + 4);
            }
            part0 -= 0x80 << 21;

            uint64 value = uint64(part0);
            
            b = uint8(data[offset + 4]);
            uint32 part1 = uint32(b);
            if (b < 0x80) {
                return (value + (uint64(part1) << 28), offset + 5);
            }
            part1 -= 0x80;

            b = uint8(data[offset + 5]);
            part1 += (uint32(b) << 7);
            if (b < 0x80) {
                return (value + (uint64(part1) << 28), offset + 6);
            }
            part1 -= 0x80 << 7;

            b = uint8(data[offset + 6]);
            part1 += (uint32(b) << 14);
            if (b < 0x80) {
                return (value + (uint64(part1) << 28), offset + 7);
            }
            part1 -= 0x80 << 14;

            b = uint8(data[offset + 7]);
            part1 += (uint32(b) << 21);
            if (b < 0x80) {
                return (value + (uint64(part1) << 28), offset + 8);
            }
            part1 -= 0x80 << 21;

            value = value + (uint64(part1) << 28);

            b = uint8(data[offset + 8]);
            uint32 part2 = uint32(b);
            if (b < 0x80) {
                return (value + (uint64(part2) << 56), offset + 9);
            }
            part2 -= 0x80;

            b = uint8(data[offset + 9]);
            part2 += (uint32(b) << 7);

            if (b < 0x02) {
                return (value + (uint64(part2) << 56), offset + 10);
            }
            revert InvalidVarint();
        } else {
            uint64 value = 0;
            uint256 remaining = data.length - offset;
            uint256 maxCount = remaining < 10 ? remaining : 10;
            for (uint256 count = 0; count < maxCount; count++) {
                uint8 b = uint8(data[offset + count]);
                value |= uint64(b & 0x7F) << (count * 7);
                
                if (b <= 0x7F) {
                    // Check for overflow on the final byte
                    if (count == 9 && b >= 0x02) {
                        revert InvalidVarint();
                    }
                    return (value, count + 1);
                }
            }
            revert InvalidVarint();
        } 
    }
    
    function tendermintSpec() public pure returns (IMembershipMsgs.ProofSpec memory) {
        IMembershipMsgs.LeafOp memory leaf = IMembershipMsgs.LeafOp({
            hashOp: IMembershipMsgs.HashOp.SHA256,
            prehashKey: IMembershipMsgs.HashOp.NO_HASH,
            prehashValue: IMembershipMsgs.HashOp.SHA256,
            prefix: hex"00"
        });

        uint32[] memory childOrder = new uint32[](2);
        childOrder[0] = 0;
        childOrder[1] = 1;

        IMembershipMsgs.InnerSpec memory inner = IMembershipMsgs.InnerSpec({
            childOrder: childOrder,
            minPrefixLength: 1,
            maxPrefixLength: 1,
            childSize: 32,
            emptyChild: "",
            hashOp: IMembershipMsgs.HashOp.SHA256
        });

        return IMembershipMsgs.ProofSpec({
            specType: IMembershipMsgs.SpecType.TENDERMINT,
            hasLeafSpec: true,
            leafOp: leaf,
            hasInnerSpec: true,
            innerSpec: inner,
            minDepth: 0,
            maxDepth: 0,
            prehashKeyBeforeComparison: false
        });
    }

    function iavlSpec() public pure returns (IMembershipMsgs.ProofSpec memory) {
        IMembershipMsgs.LeafOp memory leaf = IMembershipMsgs.LeafOp({
            hashOp: IMembershipMsgs.HashOp.SHA256,
            prehashKey: IMembershipMsgs.HashOp.NO_HASH,
            prehashValue: IMembershipMsgs.HashOp.SHA256,
            prefix: hex"00"
        });

        uint32[] memory childOrder = new uint32[](2);
        childOrder[0] = 0;
        childOrder[1] = 1;

        IMembershipMsgs.InnerSpec memory inner = IMembershipMsgs.InnerSpec({
            childOrder: childOrder,
            minPrefixLength: 4,
            maxPrefixLength: 12,
            childSize: 33,
            emptyChild: "",
            hashOp: IMembershipMsgs.HashOp.SHA256
        });

        // Return proof spec
        return IMembershipMsgs.ProofSpec({
            specType: IMembershipMsgs.SpecType.IAVL,
            hasLeafSpec: true,
            leafOp: leaf,
            hasInnerSpec: true,
            innerSpec: inner,
            minDepth: 0,
            maxDepth: 0,
            prehashKeyBeforeComparison: false
        });
    }

    function ensureIavlPrefix(
        bytes memory leafPrefix,
        int64 minHeight
    ) internal pure returns (uint256) {
        uint256 offset = 0;
        (int64 height, uint256 newOffset) = readVarint(leafPrefix, offset);
        if (height < minHeight) {
            revert("Invalid height in leaf prefix");
        }
        (int64 size, uint256 newOffset2) = readVarint(leafPrefix, newOffset);
        if (size < 0) {
            revert("Invalid size in leaf prefix");
        }
        (int64 version, uint256 newOffset3) = readVarint(leafPrefix, newOffset2);
        if (version < 0) {
            revert("Invalid version in leaf prefix");
        }

        return leafPrefix.length - newOffset3;
    }
}