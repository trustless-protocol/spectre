#!/bin/bash

ENCODE_LIB_PLACEHOLDER='__$e6bc332d3f714b58adb39753770f09750e$__'
ENCODE_LIB_ADDR="b4b46bdaa835f8e4b4d8e208b6559cd267851051"

function create_binding {
    contract_dir=$1
    contract=$2
    binding_dir=$3
    echo $contract
    mkdir -p $binding_dir/${contract}
    contract_json="../out/${contract}.sol/${contract}.json"
    solc_abi=$(cat ${contract_json} | jq -r '.abi')
    solc_bin=$(cat ${contract_json} | jq -r '.bytecode.object')

    # Link Encode library placeholder with deployed address
    solc_bin="${solc_bin//$ENCODE_LIB_PLACEHOLDER/$ENCODE_LIB_ADDR}"

    mkdir -p data
    echo ${solc_abi} > data/tmp.abi
    echo ${solc_bin} > data/tmp.bin

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