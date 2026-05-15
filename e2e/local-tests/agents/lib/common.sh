#!/usr/bin/env bash

set -euo pipefail

C_RESET=$'\033[0m'
C_GREEN=$'\033[32m'
C_RED=$'\033[31m'
C_YELLOW=$'\033[33m'
C_BLUE=$'\033[34m'

generate_id() {
  local index="$1"
  printf "%03d" "$index"
}

json_escape() {
  jq -Rs .
}

write_json_safe() {
  local path="$1"
  local content="$2"
  local tmp="${path}.tmp"
  printf "%s" "$content" > "$tmp"
  mv "$tmp" "$path"
}

iso_timestamp() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

log_info() {
  printf "%s->%s %s\n" "$C_BLUE" "$C_RESET" "$*"
}

log_ok() {
  printf "%sOK%s %s\n" "$C_GREEN" "$C_RESET" "$*"
}

log_warn() {
  printf "%sWARN%s %s\n" "$C_YELLOW" "$C_RESET" "$*" >&2
}

log_err() {
  printf "%sERR%s %s\n" "$C_RED" "$C_RESET" "$*" >&2
}

parse_env_pair() {
  local pair="$1"
  if [[ "$pair" != *=* || "$pair" == "="* ]]; then
    return 1
  fi
  env_key="${pair%%=*}"
  env_val="${pair#*=}"
  [[ "$env_key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]
}

compact_status_line() {
  local passed="$1"
  local failed="$2"
  local skipped="$3"
  local total="$4"
  local current="$5"
  local completed=$((passed + failed + skipped))
  printf "\r[%s/%s] %s | Pass: %s  Fail: %s  Skip: %s" \
    "$completed" "$total" "$current" "$passed" "$failed" "$skipped"
}
