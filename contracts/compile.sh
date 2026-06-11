#!/bin/bash
set -euo pipefail

# Header and Encode are internal libraries (inlined), so contract bytecode no
# longer contains library placeholders — no link step is needed. The guard in
# create_binding still fails loudly if an unlinked placeholder ever reappears.

function create_binding {
    contract_dir=$1
    contract=$2
    binding_dir=$3
    echo $contract
    mkdir -p $binding_dir/${contract}
    contract_json="../out/${contract}.sol/${contract}.json"
    solc_abi=$(jq -r '.abi' ${contract_json})
    solc_bin=$(jq -r '.bytecode.object' ${contract_json})

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
