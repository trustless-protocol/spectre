#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"

REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
DEFAULT_OUTPUT_DIR="$SCRIPT_DIR/output"
SOURCE_DIR="e2e/local-tests/scenarios/eth-to-cosmos"

PLAN_FILE=""
ONLY_PATTERN="*"
TIMEOUT_SECONDS=1800
OUTPUT_DIR="$DEFAULT_OUTPUT_DIR"
BEFORE_ALL=""
AFTER_ALL=""
STOP_ON_FAIL=false
VERBOSE=false
declare -a EXCLUDE_PATTERNS=()
declare -a GLOBAL_ENV=()

usage() {
  cat <<'EOF'
Usage: ./agents/manager.sh [OPTIONS]

Options:
  --plan=FILE           Pre-written plan.json (default: auto-generate)
  --only=PATTERN        Include only matching scenarios (bash glob)
  --exclude=PATTERN     Exclude matching scenarios; repeatable
  --timeout=SECONDS     Per-scenario timeout (default: 1800)
  --output-dir=DIR      Workspace (default: e2e/local-tests/agents/output/)
  --env=KEY=VAL         Global env var for every scenario; repeatable
  --before-all=CMD      Shell command to run before first scenario
  --after-all=CMD       Shell command to run after last scenario
  --stop-on-fail        Halt on first non-PASS result
  --verbose             Embed full stdout/stderr in results.md
  --help                Show usage and exit
EOF
}

die_usage() {
  log_err "$1"
  usage >&2
  exit 2
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --plan=*) PLAN_FILE="${1#*=}" ;;
    --plan) shift; [[ $# -gt 0 ]] || die_usage "--plan requires a file"; PLAN_FILE="$1" ;;
    --only=*) ONLY_PATTERN="${1#*=}" ;;
    --only) shift; [[ $# -gt 0 ]] || die_usage "--only requires a pattern"; ONLY_PATTERN="$1" ;;
    --exclude=*) EXCLUDE_PATTERNS+=("${1#*=}") ;;
    --exclude) shift; [[ $# -gt 0 ]] || die_usage "--exclude requires a pattern"; EXCLUDE_PATTERNS+=("$1") ;;
    --timeout=*) TIMEOUT_SECONDS="${1#*=}" ;;
    --timeout) shift; [[ $# -gt 0 ]] || die_usage "--timeout requires seconds"; TIMEOUT_SECONDS="$1" ;;
    --output-dir=*) OUTPUT_DIR="${1#*=}" ;;
    --output-dir) shift; [[ $# -gt 0 ]] || die_usage "--output-dir requires a directory"; OUTPUT_DIR="$1" ;;
    --env=*) GLOBAL_ENV+=("${1#*=}") ;;
    --env) shift; [[ $# -gt 0 ]] || die_usage "--env requires KEY=VAL"; GLOBAL_ENV+=("$1") ;;
    --before-all=*) BEFORE_ALL="${1#*=}" ;;
    --before-all) shift; [[ $# -gt 0 ]] || die_usage "--before-all requires a command"; BEFORE_ALL="$1" ;;
    --after-all=*) AFTER_ALL="${1#*=}" ;;
    --after-all) shift; [[ $# -gt 0 ]] || die_usage "--after-all requires a command"; AFTER_ALL="$1" ;;
    --stop-on-fail) STOP_ON_FAIL=true ;;
    --verbose) VERBOSE=true ;;
    --help) usage; exit 0 ;;
    *) die_usage "unknown option: $1" ;;
  esac
  shift
done

[[ "$TIMEOUT_SECONDS" =~ ^[0-9]+$ && "$TIMEOUT_SECONDS" -gt 0 ]] || die_usage "--timeout must be a positive integer"
for pair in "${GLOBAL_ENV[@]}"; do
  parse_env_pair "$pair" || die_usage "invalid --env pair: $pair"
done

for cmd in jq timeout bash find sort; do
  command -v "$cmd" >/dev/null || {
    log_err "required command not found: $cmd"
    exit 2
  }
done

mkdir -p "$OUTPUT_DIR"
OUTPUT_DIR="$(cd "$OUTPUT_DIR" && pwd)"
QUEUE_DIR="$OUTPUT_DIR/queue"
DONE_DIR="$QUEUE_DIR/.done"
RESULTS_DIR="$OUTPUT_DIR/results"
PLAN_OUT="$OUTPUT_DIR/plan.json"
PROGRESS_MD="$OUTPUT_DIR/progress.md"
RESULTS_MD="$OUTPUT_DIR/results.md"
RUN_STARTED="$(iso_timestamp)"
RUN_ID="$(date -u +"%Y%m%d-%H%M%S")"

rm -rf "$QUEUE_DIR" "$RESULTS_DIR"
mkdir -p "$QUEUE_DIR" "$DONE_DIR" "$RESULTS_DIR"
rm -f "$OUTPUT_DIR/stop_on_fail"
if [[ "$STOP_ON_FAIL" == "true" ]]; then
  touch "$OUTPUT_DIR/stop_on_fail"
fi

matches_filters() {
  local name="$1"
  [[ "$name" == $ONLY_PATTERN ]] || return 1
  local pattern
  for pattern in "${EXCLUDE_PATTERNS[@]}"; do
    [[ "$name" == $pattern ]] && return 1
  done
  return 0
}

env_object_json() {
  local json="{}"
  local pair env_key env_val
  for pair in "${GLOBAL_ENV[@]}"; do
    parse_env_pair "$pair"
    json="$(jq -c --arg k "$env_key" --arg v "$env_val" '. + {($k): $v}' <<<"$json")"
  done
  printf "%s" "$json"
}

if [[ -n "$PLAN_FILE" ]]; then
  [[ -r "$PLAN_FILE" ]] || die_usage "plan file is not readable: $PLAN_FILE"
  jq -e '.source_dir and (.scenarios | type == "array")' "$PLAN_FILE" >/dev/null || die_usage "invalid plan file schema: $PLAN_FILE"
  cp "$PLAN_FILE" "$PLAN_OUT"
  PLAN_SOURCE="$PLAN_FILE"
else
  mapfile -t scenario_names < <(
    find "$REPO_ROOT/$SOURCE_DIR" -maxdepth 1 -type f -name 'run-*.sh' -printf '%f\n' | sort
  )
  filtered=()
  for name in "${scenario_names[@]}"; do
    case "$name" in
      transfer.sh|check-*.sh) continue ;;
    esac
    if matches_filters "$name"; then
      filtered+=("$name")
    fi
  done

  env_json="$(env_object_json)"
  scenarios_json="[]"
  for name in "${filtered[@]}"; do
    scenarios_json="$(jq -c \
      --arg name "$name" \
      --arg path "$SOURCE_DIR/$name" \
      --argjson timeout "$TIMEOUT_SECONDS" \
      --argjson env "$env_json" \
      '. + [{name: $name, path: $path, timeout_seconds: $timeout, env: $env}]' \
      <<<"$scenarios_json")"
  done

  jq -n \
    --arg generated_at "$(iso_timestamp)" \
    --arg source_dir "$SOURCE_DIR" \
    --argjson scenarios "$scenarios_json" \
    '{generated_at: $generated_at, source_dir: $source_dir, scenarios: $scenarios}' > "$PLAN_OUT"
  PLAN_SOURCE="auto-generated from $SOURCE_DIR/run-*.sh"
fi

total="$(jq '.scenarios | length' "$PLAN_OUT")"
if [[ "$total" -eq 0 ]]; then
  log_warn "no scenarios matched"
fi

global_env_json="$(env_object_json)"
for ((i = 0; i < total; i++)); do
  id="$(generate_id "$((i + 1))")"
  scenario="$(jq -r ".scenarios[$i].name" "$PLAN_OUT")"
  scenario_path="$(jq -r ".scenarios[$i].path" "$PLAN_OUT")"
  scenario_dir="$(dirname "$scenario_path")"
  scenario_timeout="$(jq -r ".scenarios[$i].timeout_seconds // $TIMEOUT_SECONDS" "$PLAN_OUT")"
  scenario_env="$(jq -c ".scenarios[$i].env // {}" "$PLAN_OUT")"
  merged_env="$(jq -c --argjson a "$global_env_json" --argjson b "$scenario_env" -n '$a + $b')"

  jq -n \
    --arg id "$id" \
    --arg scenario "$scenario" \
    --arg scenario_dir "$scenario_dir" \
    --arg queued_at "$(iso_timestamp)" \
    --argjson timeout "$scenario_timeout" \
    --argjson env "$merged_env" \
    '{id: $id, scenario: $scenario, scenario_dir: $scenario_dir, timeout_seconds: $timeout, env: $env, queued_at: $queued_at}' \
    > "$QUEUE_DIR/${id}-${scenario%.sh}.json"
done

declare -A STATUS_BY_ID=()
declare -A RESULT_BY_ID=()

status_counts() {
  local passed=0 failed=0 skipped=0 timed_out=0 errors=0
  local id status
  for id in "${!STATUS_BY_ID[@]}"; do
    status="${STATUS_BY_ID[$id]}"
    case "$status" in
      PASS) passed=$((passed + 1)) ;;
      FAIL) failed=$((failed + 1)) ;;
      TIMEOUT) timed_out=$((timed_out + 1)); failed=$((failed + 1)) ;;
      ERROR) errors=$((errors + 1)); failed=$((failed + 1)) ;;
      SKIPPED) skipped=$((skipped + 1)) ;;
    esac
  done
  printf "%s %s %s %s %s\n" "$passed" "$failed" "$skipped" "$timed_out" "$errors"
}

write_progress_md() {
  local i id scenario result status duration started finished passed failed skipped _timed_out _errors
  {
    echo "# Test Progress"
    echo
    echo "**Run ID:** $RUN_ID  "
    echo "**Started:** $RUN_STARTED  "
    echo "**Updated:** $(iso_timestamp)"
    echo
    echo "| # | Scenario | Status | Duration | Started | Finished |"
    echo "|---|----------|--------|----------|---------|----------|"
    for ((i = 0; i < total; i++)); do
      id="$(generate_id "$((i + 1))")"
      scenario="$(jq -r ".scenarios[$i].name" "$PLAN_OUT")"
      result="${RESULT_BY_ID[$id]:-}"
      if [[ -n "$result" && -f "$result" ]]; then
        status="$(jq -r '.status' "$result")"
        duration="$(jq -r '.duration_seconds | tostring + "s"' "$result")"
        started="$(jq -r 'if .started_at then (.started_at | split("T")[1] | sub("Z$"; "")) else "--" end' "$result")"
        finished="$(jq -r 'if .finished_at then (.finished_at | split("T")[1] | sub("Z$"; "")) else "--" end' "$result")"
      else
        status="${STATUS_BY_ID[$id]:-QUEUED}"
        duration="--"
        started="--"
        finished="--"
      fi
      echo "| $((i + 1)) | $scenario | $status | $duration | $started | $finished |"
    done
    read -r passed failed skipped _timed_out _errors < <(status_counts)
    echo
    echo "**Summary:** $passed passed, $failed failed, $skipped skipped | $((passed + failed + skipped))/$total completed"
  } > "$PROGRESS_MD.tmp"
  mv "$PROGRESS_MD.tmp" "$PROGRESS_MD"
}

write_results_header() {
  {
    echo "# Test Results"
    echo
    echo "**Run started:** $RUN_STARTED  "
    echo "**Plan source:** $PLAN_SOURCE  "
    echo "**Timeout:** ${TIMEOUT_SECONDS}s per scenario  "
    echo "**Stop on fail:** $STOP_ON_FAIL  "
    echo "**Verbose:** $VERBOSE"
  } > "$RESULTS_MD"
}

append_log_details() {
  local title="$1"
  local file="$2"
  local mode="$3"
  echo "<details><summary>$title</summary>" >> "$RESULTS_MD"
  echo >> "$RESULTS_MD"
  echo '```' >> "$RESULTS_MD"
  if [[ ! -s "$file" ]]; then
    echo "(empty)" >> "$RESULTS_MD"
  elif [[ "$mode" == "full" ]]; then
    cat "$file" >> "$RESULTS_MD"
  else
    tail -n 50 "$file" >> "$RESULTS_MD"
  fi
  echo '```' >> "$RESULTS_MD"
  echo "</details>" >> "$RESULTS_MD"
}

run_hook() {
  local label="$1"
  local command="$2"
  local abort_on_fail="$3"
  local hook_slug
  hook_slug="$(tr '[:upper:] ' '[:lower:]-' <<<"$label")"
  local stdout_file="$RESULTS_DIR/${hook_slug}-stdout.log"
  local stderr_file="$RESULTS_DIR/${hook_slug}-stderr.log"
  local start_epoch end_epoch duration exit_code
  start_epoch="$(date +%s)"
  set +e
  timeout "$TIMEOUT_SECONDS" bash -c "cd '$REPO_ROOT/e2e/local-tests' && $command" >"$stdout_file" 2>"$stderr_file"
  exit_code=$?
  set -e
  end_epoch="$(date +%s)"
  duration=$((end_epoch - start_epoch))

  echo >> "$RESULTS_MD"
  echo "---" >> "$RESULTS_MD"
  echo >> "$RESULTS_MD"
  if [[ "$exit_code" -eq 0 ]]; then
    echo "## $label" >> "$RESULTS_MD"
    echo >> "$RESULTS_MD"
    echo "**Status:** PASS | **Duration:** ${duration}s | **Exit code:** 0" >> "$RESULTS_MD"
  elif [[ "$abort_on_fail" == "true" ]]; then
    echo "## $label - ABORT" >> "$RESULTS_MD"
    echo >> "$RESULTS_MD"
    echo "**Status:** FAIL | **Duration:** ${duration}s | **Exit code:** $exit_code" >> "$RESULTS_MD"
    echo >> "$RESULTS_MD"
    echo "> Run aborted. No scenarios were executed." >> "$RESULTS_MD"
  else
    echo "## $label" >> "$RESULTS_MD"
    echo >> "$RESULTS_MD"
    echo "**Status:** FAIL | **Duration:** ${duration}s | **Exit code:** $exit_code" >> "$RESULTS_MD"
  fi
  echo >> "$RESULTS_MD"
  append_log_details "stdout" "$stdout_file" "$(if [[ "$VERBOSE" == "true" ]]; then echo full; else echo tail; fi)"
  echo >> "$RESULTS_MD"
  append_log_details "stderr" "$stderr_file" "$(if [[ "$VERBOSE" == "true" ]]; then echo full; else echo tail; fi)"

  return "$exit_code"
}

append_result_md() {
  local result="$1"
  local id scenario status duration exit_code started finished stdout_file stderr_file mode
  id="$(jq -r '.id' "$result")"
  scenario="$(jq -r '.scenario' "$result")"
  status="$(jq -r '.status' "$result")"
  duration="$(jq -r '.duration_seconds' "$result")"
  exit_code="$(jq -r '.exit_code' "$result")"
  started="$(jq -r '.started_at | split("T")[1] | sub("Z$"; "")' "$result")"
  finished="$(jq -r '.finished_at | split("T")[1] | sub("Z$"; "")' "$result")"
  stdout_file="$(jq -r '.stdout_file' "$result")"
  stderr_file="$(jq -r '.stderr_file' "$result")"
  mode="$(if [[ "$VERBOSE" == "true" ]]; then echo full; else echo tail; fi)"

  echo >> "$RESULTS_MD"
  echo "---" >> "$RESULTS_MD"
  echo >> "$RESULTS_MD"
  echo "## $id. $scenario - $status (${duration}s)" >> "$RESULTS_MD"
  echo >> "$RESULTS_MD"
  echo "- **Path:** \`$SOURCE_DIR/$scenario\`" >> "$RESULTS_MD"
  echo "- **Exit code:** $exit_code" >> "$RESULTS_MD"
  echo "- **Started:** $started | **Finished:** $finished" >> "$RESULTS_MD"
  echo >> "$RESULTS_MD"
  append_log_details "stdout (last 50 lines)" "$stdout_file" "$mode"
  echo >> "$RESULTS_MD"
  append_log_details "stderr" "$stderr_file" "$mode"
}

mark_remaining_skipped() {
  local reason="$1"
  local i id scenario result existing_status
  for ((i = 0; i < total; i++)); do
    id="$(generate_id "$((i + 1))")"
    [[ -n "${STATUS_BY_ID[$id]:-}" ]] && continue
    if [[ -f "$RESULTS_DIR/${id}.json" ]]; then
      existing_status="$(jq -r '.status' "$RESULTS_DIR/${id}.json")"
      STATUS_BY_ID[$id]="$existing_status"
      RESULT_BY_ID[$id]="$RESULTS_DIR/${id}.json"
      append_result_md "$RESULTS_DIR/${id}.json"
      continue
    fi
    scenario="$(jq -r ".scenarios[$i].name" "$PLAN_OUT")"
    result="$RESULTS_DIR/${id}.json"
    jq -n \
      --arg id "$id" \
      --arg scenario "$scenario" \
      --arg reason "$reason" \
      '{id: $id, scenario: $scenario, status: "SKIPPED", exit_code: 0, duration_seconds: 0, started_at: null, finished_at: null, stdout_file: "", stderr_file: "", stdout_tail: $reason, stderr_tail: ""}' \
      > "$result"
    STATUS_BY_ID[$id]="SKIPPED"
    RESULT_BY_ID[$id]="$result"
    echo >> "$RESULTS_MD"
    echo "---" >> "$RESULTS_MD"
    echo >> "$RESULTS_MD"
    echo "## $id. $scenario - SKIPPED" >> "$RESULTS_MD"
    echo >> "$RESULTS_MD"
    echo "> $reason" >> "$RESULTS_MD"
  done
}

append_final_summary() {
  local passed failed skipped timed_out errors total_duration status_text result
  read -r passed failed skipped timed_out errors < <(status_counts)
  total_duration=0
  for result in "$RESULTS_DIR"/*.json; do
    [[ -f "$result" ]] || continue
    total_duration=$((total_duration + $(jq -r '.duration_seconds // 0' "$result")))
  done
  if [[ "$failed" -eq 0 ]]; then
    status_text="PASS"
  else
    status_text="FAIL"
  fi
  {
    echo
    echo "---"
    echo
    echo "## Final Summary"
    echo
    echo "| Metric | Count |"
    echo "|--------|-------|"
    echo "| Total scenarios | $total |"
    echo "| Passed | $passed |"
    echo "| Failed | $((failed - timed_out - errors)) |"
    echo "| Timed out | $timed_out |"
    echo "| Errors | $errors |"
    echo "| Skipped | $skipped |"
    echo "| **Total duration** | **${total_duration}s** |"
    echo
    echo "**Overall:** $status_text"
  } >> "$RESULTS_MD"
}

write_results_header
write_progress_md

if [[ -n "$BEFORE_ALL" ]]; then
  if ! run_hook "Before-All Hook" "$BEFORE_ALL" true; then
    mark_remaining_skipped "Skipped because before-all hook failed."
    write_progress_md
    append_final_summary
    log_err "before-all hook failed; no scenarios executed"
    exit 1
  fi
fi

log_info "running $total scenario(s); output: $OUTPUT_DIR"
"$SCRIPT_DIR/executor.sh" "$OUTPUT_DIR" &
EXECUTOR_PID=$!

while true; do
  mapfile -t result_files < <(find "$RESULTS_DIR" -maxdepth 1 -type f -name '*.json' | sort)
  for result in "${result_files[@]}"; do
    [[ -f "$result" ]] || continue
    jq -e '.id and .status' "$result" >/dev/null 2>&1 || continue
    id="$(jq -r '.id' "$result")"
    [[ -n "${STATUS_BY_ID[$id]:-}" ]] && continue
    status="$(jq -r '.status' "$result")"
    STATUS_BY_ID[$id]="$status"
    RESULT_BY_ID[$id]="$result"
    append_result_md "$result"
    write_progress_md

    read -r passed failed skipped _timed_out _errors < <(status_counts)
    current="$(jq -r '.scenario' "$result") $status"
    compact_status_line "$passed" "$failed" "$skipped" "$total" "$current"

    if [[ "$STOP_ON_FAIL" == "true" && "$status" != "PASS" ]]; then
      touch "$QUEUE_DIR/.stop"
      wait "$EXECUTOR_PID" 2>/dev/null || true
      mark_remaining_skipped "Skipped because a previous scenario failed and --stop-on-fail was enabled."
      write_progress_md
      break 2
    fi
  done

  if [[ $(( ${#STATUS_BY_ID[@]} )) -ge "$total" ]]; then
    break
  fi

  if ! kill -0 "$EXECUTOR_PID" 2>/dev/null; then
    wait "$EXECUTOR_PID" 2>/dev/null || true
    if [[ $(( ${#STATUS_BY_ID[@]} )) -lt "$total" ]]; then
      mark_remaining_skipped "Skipped because executor exited before all scenarios produced results."
      write_progress_md
    fi
    break
  fi

  sleep 1
done

wait "$EXECUTOR_PID" 2>/dev/null || true
printf "\n"

hook_failed=false
if [[ -n "$AFTER_ALL" ]]; then
  if ! run_hook "After-All Hook" "$AFTER_ALL" false; then
    hook_failed=true
  fi
fi

append_final_summary
write_progress_md

read -r passed failed skipped timed_out errors < <(status_counts)
if [[ "$failed" -eq 0 && "$hook_failed" == "false" ]]; then
  log_ok "agents completed: $passed passed, $failed failed, $skipped skipped"
  exit 0
fi

log_err "agents completed with failures: $passed passed, $failed failed, $skipped skipped, $timed_out timed out, $errors errors"
exit 1
