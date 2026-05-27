#ifndef ZK_ECIP_H
#define ZK_ECIP_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Result structure for zk_ecip_hint
 *
 * On success: q_ptr, a_num_ptr, a_den_ptr, b_num_ptr, b_den_ptr contain serialized BigUints
 * On error: error_msg contains the error message
 *
 * Each field (q, a_num, etc.) is serialized as:
 * [count: u32][len1: u32][bytes1...][len2: u32][bytes2...]...
 */
typedef struct {
    uint8_t* q_ptr;
    uint32_t q_len;
    uint8_t* a_num_ptr;
    uint32_t a_num_len;
    uint8_t* a_den_ptr;
    uint32_t a_den_len;
    uint8_t* b_num_ptr;
    uint32_t b_num_len;
    uint8_t* b_den_ptr;
    uint32_t b_den_len;
    uint8_t* error_msg;
    uint32_t error_len;
} ZkEcipResult;

/**
 * Compute ECIP hint for elliptic curve scalar multiplication
 *
 * @param points_ptr Pointer to serialized points (BigUints)
 * @param points_len Length of points data
 * @param scalars_ptr Pointer to serialized scalars (BigUints)
 * @param scalars_len Length of scalars data
 * @param curve_id Curve identifier:
 *                 0 = BN254
 *                 1 = BLS12381
 *                 2 = SECP256K1
 *                 3 = SECP256R1
 *                 4 = X25519
 *                 5 = Grumpkin
 *
 * @return ZkEcipResult containing result or error
 *
 * Note: Caller must free result using zk_ecip_free_result()
 */
ZkEcipResult zk_ecip_hint_c(
    const uint8_t* points_ptr,
    uint32_t points_len,
    const uint8_t* scalars_ptr,
    uint32_t scalars_len,
    uint32_t curve_id
);

/**
 * Free memory allocated by zk_ecip_hint_c
 *
 * @param result Result structure returned by zk_ecip_hint_c
 */
void zk_ecip_free_result(ZkEcipResult result);

#ifdef __cplusplus
}
#endif

#endif /* ZK_ECIP_H */
