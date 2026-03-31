#!/bin/sh

set -euxo pipefail

# This script expects an environment variable PROGRAM_VERSIONS

PROGRAMS="groth16-ics07-tendermint-membership groth16-ics07-tendermint-update-client groth16-ics07-tendermint-uc-and-membership groth16-ics07-tendermint-misbehaviour"

# Check if the environment variable is set
if [ -z "$PROGRAM_VERSIONS" ]; then
  echo "Error: PROGRAM_VERSIONS environment variable is not set or is empty." >&2
  exit 1
fi

# Loop through each version provided
for version in $PROGRAM_VERSIONS; do
  target_dir="/usr/local/bin/groth16-programs/$version"
  mkdir -p "$target_dir"

  # Download each program for the current version
  for program in $PROGRAMS; do
      wget --no-check-certificate https://github.com/cosmos/solidity-ibc-eureka/releases/download/groth16-programs-$version/$program -O $target_dir/$program
  done
done

/usr/local/bin/relayer $@
