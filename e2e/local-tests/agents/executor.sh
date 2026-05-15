#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  echo "Usage: $0 <output-dir>" >&2
}

if [[ $# -ne 1 ]]; then
  usage
  exit 2
fi

OUTPUT_DIR="$(cd "$1" && pwd)"
QUEUE_DIR="$OUTPUT_DIR/queue"
DONE_DIR="$QUEUE_DIR/.done"
RESULTS_DIR="$OUTPUT_DIR/results"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

if [[ ! -d "$QUEUE_DIR" || ! -d "$DONE_DIR" || ! -d "$RESULTS_DIR" ]]; then
  log_err "output directory is missing queue/results structure: $OUTPUT_DIR"
  exit 2
fi

while true; do
  next_cmd="$(find "$QUEUE_DIR" -maxdepth 1 -type f -name '*.json' | sort | head -n1)"

  if [[ -z "$next_cmd" ]]; then
    exit 0
  fi

  cmd_name="$(basename "$next_cmd")"
  id="${cmd_name%%-*}"
  done_cmd="$DONE_DIR/$cmd_name"
  mv "$next_cmd" "$done_cmd"

  scenario="$(jq -r '.scenario' "$done_cmd")"
  scenario_dir="$(jq -r '.scenario_dir' "$done_cmd")"
  timeout_sec="$(jq -r '.timeout_seconds' "$done_cmd")"

  stdout_file="$RESULTS_DIR/${id}-stdout.log"
  stderr_file="$RESULTS_DIR/${id}-stderr.log"
  started_at="$(iso_timestamp)"
  start_epoch="$(date +%s)"

  set +e
  timeout "$timeout_sec" bash -c '
    set -euo pipefail
    cd "$1"
    cmd_json="$2"
    scenario_path="$3/$4"
    while IFS= read -r item; do
      key="$(jq -r ".key" <<<"$item")"
      val="$(jq -r ".value" <<<"$item")"
      export "$key=$val"
    done < <(jq -c ".env // {} | to_entries[]" "$cmd_json")
    bash "$scenario_path"
  ' bash "$REPO_ROOT" "$done_cmd" "$scenario_dir" "$scenario" >"$stdout_file" 2>"$stderr_file"
  exit_code=$?
  set -e

  finished_at="$(iso_timestamp)"
  end_epoch="$(date +%s)"
  duration=$((end_epoch - start_epoch))

  if [[ "$exit_code" -eq 0 ]]; then
    status="PASS"
  elif [[ "$exit_code" -eq 124 ]]; then
    status="TIMEOUT"
  elif [[ "$exit_code" -eq 125 || "$exit_code" -eq 126 || "$exit_code" -eq 127 || "$exit_code" -gt 128 ]]; then
    status="ERROR"
  else
    status="FAIL"
  fi

  jq -n \
    --arg id "$id" \
    --arg scenario "$scenario" \
    --arg status "$status" \
    --arg started_at "$started_at" \
    --arg finished_at "$finished_at" \
    --arg stdout_file "$stdout_file" \
    --arg stderr_file "$stderr_file" \
    --rawfile stdout_tail <(tail -n 50 "$stdout_file") \
    --rawfile stderr_tail <(tail -n 50 "$stderr_file") \
    --argjson exit_code "$exit_code" \
    --argjson duration "$duration" \
    '{
      id: $id,
      scenario: $scenario,
      status: $status,
      exit_code: $exit_code,
      duration_seconds: $duration,
      started_at: $started_at,
      finished_at: $finished_at,
      stdout_file: $stdout_file,
      stderr_file: $stderr_file,
      stdout_tail: $stdout_tail,
      stderr_tail: $stderr_tail
    }' > "$RESULTS_DIR/${id}.json.tmp"
  mv "$RESULTS_DIR/${id}.json.tmp" "$RESULTS_DIR/${id}.json"

  if [[ -f "$OUTPUT_DIR/stop_on_fail" && "$status" != "PASS" ]]; then
    touch "$QUEUE_DIR/.stop"
    exit 0
  fi

  if [[ -f "$QUEUE_DIR/.stop" ]]; then
    exit 0
  fi
done
