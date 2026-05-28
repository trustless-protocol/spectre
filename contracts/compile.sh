#!/bin/bash
set -euo pipefail

ENCODE_LIB_PLACEHOLDER='__$e6bc332d3f714b58adb39753770f09750e$__'
ENCODE_LIB_ADDR="b4b46bdaa835f8e4b4d8e208b6559cd267851051"
HEADER_LIB_PLACEHOLDER='__$7440a880b7578767f72184d998805816e4$__'
HEADER_LIB_ADDR="a3c616dd54F6BB35a736cD6968c8EF7176faCACc"

function link_library {
    local solc_bin=$1
    local placeholder=$2
    local addr=$3
    local name=$4

    if [[ "${solc_bin}" != *"${placeholder}"* ]]; then
        echo "${solc_bin}"
        return
    fi

    addr="${addr#0x}"
    if [[ ! "${addr}" =~ ^[0-9a-fA-F]{40}$ ]]; then
        echo "missing/invalid hardcoded ${name} library address for placeholder ${placeholder}" >&2
        return 1
    fi

    echo "${solc_bin//${placeholder}/${addr}}"
}

function create_binding {
    contract_dir=$1
    contract=$2
    binding_dir=$3
    echo $contract
    mkdir -p $binding_dir/${contract}
    contract_json="../out/${contract}.sol/${contract}.json"
    solc_abi=$(jq -r '.abi' ${contract_json})
    solc_bin=$(jq -r '.bytecode.object' ${contract_json})

    solc_bin=$(link_library "${solc_bin}" "${ENCODE_LIB_PLACEHOLDER}" "${ENCODE_LIB_ADDR}" "ENCODE")
    solc_bin=$(link_library "${solc_bin}" "${HEADER_LIB_PLACEHOLDER}" "${HEADER_LIB_ADDR}" "HEADER")
    if [[ "${solc_bin}" == *'__$'* ]]; then
        echo "unlinked library placeholder remains in ${contract}: ${solc_bin}" >&2
        return 1
    fi

    mkdir -p data
    echo "${solc_abi}" > data/tmp.abi
    echo "${solc_bin}" > data/tmp.bin

    rm -f $binding_dir/${contract}/binding.go
    abigen --bin=data/tmp.bin --abi=data/tmp.abi --pkg=contract${contract} --out=$binding_dir/${contract}/binding.go
}

forge clean
forge build

create_binding ./programs/ "Membership" ../relayer/bindings
create_binding ./programs/ "Misbehaviour" ../relayer/bindings
create_binding ./programs/ "UpdateClient" ../relayer/bindings
create_binding ./light-clients/ "Groth16ICS07Tendermint" ../relayer/bindings
create_binding ./ "ICS20Transfer" ../relayer/bindings
create_binding ./ "ICS26Router" ../relayer/bindings

# ./compile.sh ./ ICS26Router ./bindings 
