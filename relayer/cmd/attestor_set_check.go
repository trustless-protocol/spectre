package main

import (
	"context"
	"fmt"
	"log"

	tendermintClient "relayer/client"
	"relayer/services"

	"attestor/types/attestation"
)

// validateL2AttestorSet compares the attestor set this config collects
// signatures for with the set the L2 wasm client on Cosmos trusts (#575).
//
// The client stores its set in its client state and can rotate it after
// creation. A config still holding the old set starts cleanly, collects
// signatures the client refuses, and fails every L2 header from then on --
// while the log only says that headers are refused.
//
// A mismatch is a finding. A client that cannot be read is not: availability is
// not evidence the config is wrong, the same doctrine as the other probes.
func validateL2AttestorSet(stdCtx context.Context, cfg *appConfig) []configFinding {
	if cfg == nil || !cfg.HasL2ToCosmos {
		return nil
	}
	src := cfg.L2ToCosmosConfig
	if src.L2WasmClientID == "" || src.TmRpcUrl == "" {
		return nil
	}
	want, err := src.Attestors.SetHash()
	if err != nil {
		return nil // an invalid configured set is reported by the config checks
	}
	onChain, err := probeL2ClientAttestors(stdCtx, src.TmRpcUrl, src.L2WasmClientID)
	if err != nil {
		log.Printf("WARN [start] could not read the attestor set of L2 wasm client %s: %v; "+
			"a mismatch would show up as refused headers", src.L2WasmClientID, err)
		return nil
	}
	got, err := onChain.SetHash()
	if err != nil {
		log.Printf("WARN [start] L2 wasm client %s holds an attestor set this build cannot hash: %v", src.L2WasmClientID, err)
		return nil
	}
	if got != want {
		return []configFinding{{
			key: "l2_to_cosmos.attestors",
			detail: fmt.Sprintf("set hash %x (%d keys, threshold %d) differs from the set L2 wasm client %s trusts "+
				"(%x, %d keys, threshold %d). Every header this relayer collects would be refused; copy the "+
				"client's current set into the config",
				want, len(src.Attestors.PublicKeys), src.Attestors.Threshold, src.L2WasmClientID,
				got, len(onChain.PublicKeys), onChain.Threshold),
		}}
	}
	log.Printf("[start] attestor set verified against L2 wasm client %s: %d keys, threshold %d",
		src.L2WasmClientID, len(onChain.PublicKeys), onChain.Threshold)
	return nil
}

// probeL2ClientAttestors reads the attestor set of one L2 wasm client, bounded
// by configProbeTimeout.
func probeL2ClientAttestors(stdCtx context.Context, url, clientID string) (attestation.AttestorConfig, error) {
	ctx, cancel := context.WithTimeout(stdCtx, configProbeTimeout)
	defer cancel()
	c, err := tendermintClient.DialCosmosRPC(url, "/websocket", configProbeTimeout)
	if err != nil {
		return attestation.AttestorConfig{}, err
	}
	return services.L2ClientAttestors(ctx, c, clientID)
}
