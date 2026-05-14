# Manager & Executor Agent System — Implementation Guide

> **Status:** Ready to implement  
> **Scope:** Orchestrate ETH→Cosmos E2E scenario scripts in `e2e/local-tests/scenarios/eth-to-cosmos/`  
> **Assumption:** Local environment (Kurtosis ETH + Cosmos + relayer) is already running. Agents only orchestrate; they do not set up infrastructure.

---

## Table of Contents

1. [Overview](#1-overview)
2. [File Layout](#2-file-layout)
3. [Agent Behavior Summary](#3-agent-behavior-summary)
4. [Communication Protocol](#4-communication-protocol)
5. [Shared Library (`lib/common.sh`)](#5-shared-library-libcommonsh)
6. [Manager (`manager.sh`)](#6-manager-managersh)
7. [Executor (`executor.sh`)](#7-executor-executorsh)
8. [Plan JSON](#8-plan-json)
9. [Command JSON](#9-command-json)
10. [Result JSON](#10-result-json)
11. [Progress Markdown (`progress.md`)](#11-progress-markdown-progressmd)
12. [Results Markdown (`results.md`)](#12-results-markdown-resultsmd)
13. [CLI Reference](#13-cli-reference)
14. [Exit Codes](#14-exit-codes)
15. [Example Usage](#15-example-usage)
16. [Implementation Checklist](#16-implementation-checklist)

---

## 1. Overview

Two bash scripts coordinate running E2E test scenarios:

- **Manager** generates an execution plan, optionally runs pre/post hooks, spawns the executor, and monitors results.
- **Executor** reads command JSONs from a queue, runs scenarios sequentially in isolated subshells, and writes structured result JSONs.

All state lives in a single `output/` directory. Both agents communicate via JSON files on disk.

---

## 2. File Layout

```
e2e/local-tests/agents/
├── manager.sh              # Entry point: plan, dispatch, monitor, report
├── executor.sh             # Worker: sequential queue processor
├── lib/
│   └── common.sh           # JSON helpers, timestamps, shared logging
└── output/                 # Created at runtime by manager
    ├── plan.json           # Auto-generated scenario manifest
    ├── queue/              # Command JSONs dropped by manager
    │   └── .done/          # Executor moves processed commands here (audit trail)
    ├── results/            # Result JSONs + captured stdout/stderr .log files
    ├── progress.md         # Live status table (rewritten on each result)
    └── results.md          # Incremental report + final summary
```

> No existing files in the repo are modified.

---

## 3. Agent Behavior Summary

| Aspect | Decision |
|--------|----------|
| **Concurrency** | Sequential (single executor, FIFO queue). |
| **Environment** | Assumed already running. Agents never call `setup/*.sh`. |
| **Retry** | None. First result is final. |
| **Hook failure** | `--before-all` failure aborts immediately: no executor, no scenarios, no `--after-all`. |
| **Default timeout** | **1800 seconds** (30 minutes) per scenario. |
| **Output directory** | `e2e/local-tests/agents/output/` (default). |
| **Plan generation** | Auto-scan `scenarios/eth-to-cosmos/run-*.sh`; exclude `transfer.sh` and `check-*.sh`. |
| **Failure handling** | Continue by default. `--stop-on-fail` halts executor and marks remaining as `SKIPPED`. |
| **JSON dependency** | `jq` (already used throughout the repo). |
| **Live stdout** | Compact single-line status overwritten in terminal; full summary printed at end. |
| **Skipped items** | Included in `results.md` with a SKIPPED note. |

---

## 4. Communication Protocol

All coordination is **file-based JSON**.

```
Manager                        Executor
  │                              │
  ├── writes queue/001-*.json ──→│
  ├── writes queue/002-*.json ──→│
  │                              │
  │                       ┌──────┘
  │                       │ reads oldest command
  │                       │ runs scenario with timeout
  │                       │ captures stdout/stderr
  │                       │ writes results/001.json
  │                       │ moves command to queue/.done/
  │←── detects new result ─┘
  │   rewrites progress.md
  │   appends to results.md
```

**Stop signal:** Manager touches `queue/.stop`; executor checks it between iterations and exits gracefully after finishing the current scenario.

---

## 5. Shared Library (`lib/common.sh`)

Both agents `source` this file.

### Required Functions

| Function | Signature | Purpose |
|----------|-----------|---------|
| `generate_id(index)` | `index` (int) → `"001"` | Zero-padded 3-digit ID from 0-based index. |
| `json_escape(string)` | raw string → escaped string | Escape `"`, `\`, `\b`, `\f`, `\n`, `\r`, `\t` for JSON embedding. |
| `write_json_safe(path, content)` | path, string → void | Atomic write: write to `.tmp`, then `mv`. |
| `iso_timestamp()` | void → `"2026-05-14T12:00:00Z"` | UTC ISO-8601 timestamp. |
| `log_info(msg)` | string → stdout | Colored `→` prefix. |
| `log_ok(msg)` | string → stdout | Colored `✓` prefix. |
| `log_warn(msg)` | string → stderr | Colored `⚠` prefix. |
| `log_err(msg)` | string → stderr | Colored `✗` prefix. |
| `parse_env_pair(pair)` | `"KEY=VAL"` → `key` and `val` vars | Split at first `=`. |
| `compact_status_line(passed, failed, skipped, total, current_scenario)` | ints, string → stdout | Single overwriteable terminal line (use `\r`). |

### Design Note

Do **not** source `../lib/common.sh` to avoid pulling in contract-discovery logic. Define minimal colored logging using the same ANSI codes for visual consistency.

---

## 6. Manager (`manager.sh`)

### 6.1 CLI Interface

```bash
./e2e/local-tests/agents/manager.sh [OPTIONS]

Options:
  --plan=FILE           Pre-written plan.json (default: auto-generate)
  --only=PATTERN        Include only matching scenarios (bash glob)
  --exclude=PATTERN     Exclude matching scenarios
  --timeout=SECONDS     Per-scenario timeout (default: 1800)
  --output-dir=DIR      Workspace (default: e2e/local-tests/agents/output/)
  --env=KEY=VAL         Global env var for every scenario (repeatable)
  --before-all=CMD      Shell command to run before first scenario
  --after-all=CMD       Shell command to run after last scenario
  --stop-on-fail        Halt on first failure; mark rest as SKIPPED
  --verbose             Embed full stdout/stderr in results.md
  --help                Show usage and exit
```

### 6.2 Auto-Plan Generation

If `--plan` is omitted:

```bash
# 1. Scan target directory
for f in "$REPO_ROOT/e2e/local-tests/scenarios/eth-to-cosmos"/run-*.sh; do
  name="$(basename "$f")"
  # 2. Exclude helpers
  case "$name" in
    transfer.sh) continue ;;
    check-*.sh) continue ;;
  esac
  SCENARIOS+=("$name")
done

# 3. Apply --only and --exclude bash pattern filters

# 4. Write output/plan.json
```

### 6.3 Initialization Sequence

```
1. Parse CLI args into variables
2. If --plan not given → auto-generate plan.json
3. mkdir -p output/{queue,results,queue/.done}
4. Write command JSONs to output/queue/ with 3-digit IDs (001, 002, ...)
5. Write initial progress.md (all rows = QUEUED)
6. IF --before-all:
     a. Run command with timeout 1800
     b. Capture exit code, stdout, stderr, duration
     c. IF exit != 0:
          - Write ABORT header to results.md
          - Include hook stdout/stderr
          - EXIT 1 (no executor, no --after-all)
     d. Append "Before-All Hook" section to results.md
7. Spawn: ./executor.sh <output-dir> &
8. Record executor PID
9. Enter polling loop
10. Wait for executor to finish
11. IF --after-all (and not aborted):
      a. Run command
      b. Append "After-All Hook" section to results.md
12. Write final summary block to results.md
13. Print colored summary to stdout
14. Exit with appropriate code
```

### 6.4 Polling Loop

```bash
POLL_INTERVAL=1.5  # seconds
declare -A SEEN_RESULTS

while true; do
  for result_json in "$RESULTS_DIR"/*.json; do
    [ -f "$result_json" ] || continue
    id="$(basename "$result_json" .json)"
    [ -n "${SEEN_RESULTS[$id]}" ] && continue
    SEEN_RESULTS[$id]=1

    status="$(jq -r '.status' "$result_json")"

    # Rewrite progress.md
    update_progress_md

    # Append to results.md
    append_result_to_results_md "$result_json"

    # Compact terminal status
    print_compact_status

    # Handle --stop-on-fail
    if [ "$STOP_ON_FAIL" = "true" ] && [ "$status" != "PASS" ]; then
      touch "$QUEUE_DIR/.stop"
      wait "$EXECUTOR_PID" 2>/dev/null || true
      mark_remaining_skipped
      break 2
    fi
  done

  # Check if all done
  total_cmds="$(find "$QUEUE_DIR" -maxdepth 1 -name '*.json' | wc -l)"
  # Actually count original plan items vs results
  if all_have_results; then break; fi

  sleep "$POLL_INTERVAL"
done
```

### 6.5 Hook Failure (Abort)

If `--before-all` exits non-zero or times out:
- Write `results.md` header with **ABORT** status.
- Include hook stdout/stderr.
- **Do not** spawn executor.
- **Do not** run `--after-all`.
- Exit code: `1`.

### 6.6 Compact Terminal Status

A single overwriteable line using `\r` (no newline):

```
[3/10] run-error-ack.sh RUNNING | Pass: 2  Fail: 0  Skip: 0
```

At the end, print a newline then the full summary.

---

## 7. Executor (`executor.sh`)

### 7.1 Interface

```bash
./executor.sh <output-dir>
```

### 7.2 Main Loop

```bash
QUEUE_DIR="$1/queue"
DONE_DIR="$1/queue/.done"
RESULTS_DIR="$1/results"
REPO_ROOT="$(cd "$1/../../.." && pwd)"

while true; do
  # Check stop sentinel
  [ -f "$QUEUE_DIR/.stop" ] && exit 0

  # Find oldest command (lowest numeric prefix)
  next_cmd="$(find "$QUEUE_DIR" -maxdepth 1 -name '*.json' -not -path '*/.done/*' | sort | head -n1)"

  if [ -z "$next_cmd" ]; then
    # Nothing queued. Check if we might be truly done.
    sleep 1
    continue
  fi

  cmd_name="$(basename "$next_cmd")"
  id="${cmd_name%%-*}"  # "001" from "001-run-success.json"

  # Move to .done/ for audit trail
  mv "$next_cmd" "$DONE_DIR/$cmd_name"

  # Parse command JSON
  scenario="$(jq -r '.scenario' "$DONE_DIR/$cmd_name")"
  scenario_dir="$(jq -r '.scenario_dir' "$DONE_DIR/$cmd_name")"
  timeout_sec="$(jq -r '.timeout_seconds' "$DONE_DIR/$cmd_name")"

  # Extract env vars from JSON into shell variables
  # Then export them in the subshell

  stdout_file="$RESULTS_DIR/${id}-stdout.log"
  stderr_file="$RESULTS_DIR/${id}-stderr.log"

  started_at="$(iso_timestamp)"
  start_epoch="$(date +%s)"

  # Execute with timeout
  exit_code=0
  (
    cd "$REPO_ROOT"
    # Export per-scenario env vars + global env vars
    # (Implementation: loop over jq '.env | to_entries[]')
    bash "$scenario_dir/$scenario"
  ) > "$stdout_file" 2> "$stderr_file"
  exit_code=$?

  # Note: timeout sets exit_code=124 on timeout
  finished_at="$(iso_timestamp)"
  end_epoch="$(date +%s)"
  duration=$((end_epoch - start_epoch))

  # Determine status
  if [ "$exit_code" -eq 0 ]; then
    status="PASS"
  elif [ "$exit_code" -eq 124 ]; then
    status="TIMEOUT"
  elif [ "$exit_code" -gt 128 ]; then
    status="ERROR"
  else
    status="FAIL"
  fi

  # Build tail of stdout (last 50 lines, or full if <50)
  stdout_tail="$(tail -n 50 "$stdout_file" | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')"
  stderr_tail="$(tail -n 50 "$stderr_file" | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')"

  # Write result JSON
  cat > "$RESULTS_DIR/${id}.json" <<EOF
{
  "id": "$id",
  "scenario": "$scenario",
  "status": "$status",
  "exit_code": $exit_code,
  "duration_seconds": $duration,
  "started_at": "$started_at",
  "finished_at": "$finished_at",
  "stdout_file": "$stdout_file",
  "stderr_file": "$stderr_file",
  "stdout_tail": "$stdout_tail",
  "stderr_tail": "$stderr_tail"
}
EOF
done
```

### 7.3 Status Determination

| `timeout` exit code | Scenario exit code | Result status |
|--------------------|--------------------|---------------|
| `0` | `0` | `PASS` |
| `0` | `1–127` | `FAIL` |
| `124` | any | `TIMEOUT` |
| `125` | any | `ERROR` (`timeout` command itself failed) |
| `126`, `127` | any | `ERROR` (command not executable / not found) |
| `>128` | any | `ERROR` (signal death) |

---

## 8. Plan JSON

**File:** `output/plan.json`

```json
{
  "generated_at": "2026-05-14T12:00:00Z",
  "source_dir": "e2e/local-tests/scenarios/eth-to-cosmos",
  "scenarios": [
    {
      "name": "run-success.sh",
      "path": "e2e/local-tests/scenarios/eth-to-cosmos/run-success.sh",
      "timeout_seconds": 1800,
      "env": {}
    },
    {
      "name": "run-timeout.sh",
      "path": "e2e/local-tests/scenarios/eth-to-cosmos/run-timeout.sh",
      "timeout_seconds": 1800,
      "env": {}
    }
  ]
}
```

**If user provides `--plan=FILE`:** The manager uses it as-is without auto-generation. It must conform to the schema above (at minimum: `source_dir` and `scenarios[].{name,path}`).

---

## 9. Command JSON

**File:** `output/queue/001-run-success.json`

```json
{
  "id": "001",
  "scenario": "run-success.sh",
  "scenario_dir": "e2e/local-tests/scenarios/eth-to-cosmos",
  "timeout_seconds": 1800,
  "env": {
    "AMOUNT": "1000000000"
  },
  "queued_at": "2026-05-14T12:00:00Z"
}
```

**ID format:** Zero-padded 3-digit integer, monotonically increasing from `001`.

**Env merging order:** `--env=KEY=VAL` (global) are defaults. Per-scenario `env` in `plan.json` overrides. Command JSON stores the merged final set.

---

## 10. Result JSON

**File:** `output/results/001-run-success.json`

```json
{
  "id": "001",
  "scenario": "run-success.sh",
  "status": "PASS",
  "exit_code": 0,
  "duration_seconds": 45,
  "started_at": "2026-05-14T12:00:05Z",
  "finished_at": "2026-05-14T12:00:50Z",
  "stdout_file": "/home/user/.../agents/output/results/001-stdout.log",
  "stderr_file": "/home/user/.../agents/output/results/001-stderr.log",
  "stdout_tail": "  ✓ Cosmos voucher received: 1000000000\n  ...",
  "stderr_tail": ""
}
```

**Valid statuses:** `PASS`, `FAIL`, `TIMEOUT`, `SKIPPED`, `ERROR`.

**Tail extraction:** Last 50 lines. If stdout/stderr is under 50 lines, store the full content. Newlines in JSON strings are encoded as `\n`.

---

## 11. Progress Markdown (`progress.md`)

Rewritten atomically on every new result.

```markdown
# Test Progress

**Run ID:** 20260514-120000  
**Started:** 2026-05-14T12:00:00Z  
**Updated:** 2026-05-14T12:05:00Z

| # | Scenario | Status | Duration | Started | Finished |
|---|----------|--------|----------|---------|----------|
| 1 | run-success.sh | PASS | 45s | 12:00:05 | 12:00:50 |
| 2 | run-timeout.sh | FAIL | 1800s | 12:01:00 | 12:31:00 |
| 3 | run-error-ack.sh | SKIPPED | — | — | — |

**Summary:** 1 passed, 1 failed, 1 skipped | 3/10 completed
```

**Implementation:** Write to `progress.md.tmp`, then `mv` for atomicity.

---

## 12. Results Markdown (`results.md`)

### 12.1 Header (written once at start)

```markdown
# Test Results

**Run started:** 2026-05-14T12:00:00Z  
**Plan source:** auto-generated from scenarios/eth-to-cosmos/run-*.sh  
**Timeout:** 1800s per scenario  
**Stop on fail:** yes  
**Verbose:** no
```

### 12.2 Before-All Hook Section (if `--before-all` provided)

```markdown
---

## Before-All Hook

**Status:** PASS | **Duration:** 5s | **Exit code:** 0

<details><summary>stdout</summary>

```
[content or tail]
```
</details>

<details><summary>stderr</summary>

```
[content or "empty"]
```
</details>
```

If the hook fails, replace the section header with:

```markdown
---

## Before-All Hook — ABORT

**Status:** FAIL | **Duration:** 3s | **Exit code:** 1

> Run aborted. No scenarios were executed.
```

### 12.3 Per-Scenario Section (appended incrementally)

```markdown
---

## 001. run-success.sh — PASS (45s)

- **Path:** `e2e/local-tests/scenarios/eth-to-cosmos/run-success.sh`
- **Exit code:** 0
- **Started:** 12:00:05 | **Finished:** 12:00:50

<details><summary>stdout (last 50 lines)</summary>

```
[last 50 lines]
```
</details>

<details><summary>stderr</summary>

```
[content or "(empty)"]
```
</details>
```

For `SKIPPED` items:

```markdown
---

## 003. run-error-ack.sh — SKIPPED

> Skipped because a previous scenario failed and `--stop-on-fail` was enabled.
```

### 12.4 After-All Hook Section (if `--after-all` provided)

Same format as Before-All, but labeled "After-All Hook".

### 12.5 Final Summary Block (appended at end)

```markdown
---

## Final Summary

| Metric | Count |
|--------|-------|
| Total scenarios | 10 |
| Passed | 8 |
| Failed | 1 |
| Timed out | 0 |
| Skipped | 1 |
| **Total duration** | **2h 15m** |

**Overall:** ❌ FAIL (1 scenario failed)
```

---

## 13. CLI Reference

### 13.1 Full Option Table

| Option | Argument | Default | Description |
|--------|----------|---------|-------------|
| `--plan` | `FILE` | *(auto)* | Read scenario list from a JSON file instead of auto-scanning. |
| `--only` | `PATTERN` | `*` | Bash glob. Only scenarios whose basename matches are included. |
| `--exclude` | `PATTERN` | *(none)* | Bash glob. Scenarios matching are excluded. Applied after `--only`. |
| `--timeout` | `SECONDS` | `1800` | Default timeout for each scenario (and hooks). |
| `--output-dir` | `DIR` | `agents/output/` | Runtime workspace. Created if missing. |
| `--env` | `KEY=VAL` | *(none)* | Repeatable. Injected into every scenario's environment. |
| `--before-all` | `CMD` | *(none)* | Shell command run before the first scenario. Failure aborts. |
| `--after-all` | `CMD` | *(none)* | Shell command run after the last scenario (unless aborted). |
| `--stop-on-fail` | flag | off | Touch `queue/.stop` on first non-PASS result. |
| `--verbose` | flag | off | Embed **full** stdout/stderr in `results.md` instead of tail. |
| `--help` | flag | — | Print usage and exit `0`. |

### 13.2 Argument Parsing Rules

- `--env` may appear multiple times. Later occurrences with the same `KEY` override earlier ones.
- `--only` and `--exclude` support bash glob patterns (e.g., `run-*timeout*`).
- `--before-all` and `--after-all` are passed verbatim to `bash -c "<cmd>"`.
- Unknown flags or malformed `--env` pairs → print usage to stderr, exit `2`.

---

## 14. Exit Codes

| Code | Meaning |
|------|---------|
| `0` | All scenarios and hooks passed. |
| `1` | One or more scenarios failed / timed out / errored; OR `--before-all` aborted the run; OR a hook failed. |
| `2` | CLI usage error (invalid flag, missing argument, unreadable plan file, etc.). |

---

## 15. Example Usage

### 15.1 Full default run

```bash
cd e2e/local-tests
./agents/manager.sh
```

### 15.2 Only timeout scenarios, stop on first failure

```bash
./agents/manager.sh --only="run-*timeout*" --stop-on-fail
```

### 15.3 With env vars, hooks, and verbose logging

```bash
./agents/manager.sh \
  --env="AMOUNT=500000000" \
  --before-all="./scenarios/eth-to-cosmos/check-eth.sh" \
  --after-all="./scenarios/eth-to-cosmos/check-cosmos.sh" \
  --verbose
```

### 15.4 Use a custom plan, custom output dir

```bash
./agents/manager.sh \
  --plan=smoke-tests.json \
  --output-dir=/tmp/fast-ibc-run-001 \
  --timeout=600
```

### 15.5 Exclude specific scenarios

```bash
./agents/manager.sh --exclude="run-batch-*" --exclude="run-uint256*"
```

---

## 16. Implementation Checklist

Use this list when implementing the agents.

### 16.1 `agents/lib/common.sh`

- [ ] Define ANSI color constants (`C_RESET`, `C_GREEN`, etc.).
- [ ] Implement `generate_id(index)` → zero-padded 3-digit string.
- [ ] Implement `json_escape(string)` → escaped JSON string.
- [ ] Implement `write_json_safe(path, content)` → atomic write via `.tmp` + `mv`.
- [ ] Implement `iso_timestamp()` → UTC ISO-8601.
- [ ] Implement `log_info`, `log_ok`, `log_warn`, `log_err` with colored prefixes.
- [ ] Implement `parse_env_pair(pair)` → sets `env_key` and `env_val` variables.
- [ ] Implement `compact_status_line(passed, failed, skipped, total, current)` → single `\r` line.

### 16.2 `agents/manager.sh`

- [ ] Parse CLI with a `while` loop and `case` statement.
- [ ] Validate required tools: `jq`, `bash`, `timeout` (from coreutils).
- [ ] Auto-generate `plan.json` if `--plan` not given.
- [ ] Apply `--only` and `--exclude` filters to scenario list.
- [ ] Create output directory tree (`queue`, `results`, `queue/.done`).
- [ ] Write command JSONs to `queue/` with 3-digit IDs.
- [ ] Write initial `progress.md` (all QUEUED).
- [ ] Run `--before-all` with `timeout`; handle abort on failure.
- [ ] Spawn executor, record PID.
- [ ] Polling loop: detect new `results/*.json`, rewrite `progress.md`, append to `results.md`.
- [ ] Implement `--stop-on-fail`: touch `.stop`, wait for executor, mark remaining SKIPPED.
- [ ] Run `--after-all` (unless aborted).
- [ ] Write final summary to `results.md`.
- [ ] Print final colored summary to stdout.
- [ ] Exit with correct code (`0`, `1`, or `2`).

### 16.3 `agents/executor.sh`

- [ ] Accept `<output-dir>` as sole positional argument.
- [ ] Validate `output-dir` contains expected subdirectories.
- [ ] Main loop: check `.stop` sentinel, find oldest command, move to `.done/`.
- [ ] Parse command JSON with `jq`.
- [ ] Export env vars in subshell before running scenario.
- [ ] Capture stdout/stderr to `results/<id>-{stdout,stderr}.log`.
- [ ] Determine status correctly from `timeout` exit codes.
- [ ] Build `stdout_tail` and `stderr_tail` (last 50 lines, JSON-escaped).
- [ ] Write structured result JSON to `results/<id>.json`.
- [ ] Exit gracefully when queue is empty and `.stop` is present.

### 16.4 Integration Tests (manual)

- [ ] Run manager with no args against running environment → all PASS.
- [ ] Run with `--only="run-success*"` → only one scenario executed.
- [ ] Run with `--stop-on-fail` and a failing scenario → subsequent ones SKIPPED.
- [ ] Run with `--before-all="false"` → abort, no scenarios, exit `1`.
- [ ] Verify `progress.md` is updated live (open in another terminal with `watch -n1 cat`).
- [ ] Verify `results.md` contains all sections including final summary.
- [ ] Verify `queue/.done/` contains all processed command JSONs.
- [ ] Verify `results/` contains `.json`, `-stdout.log`, and `-stderr.log` for every scenario.

---

## Appendix A: Rejected Alternatives (for reference)

| Alternative | Why Rejected |
|-------------|--------------|
| Go implementation | User chose bash for consistency with existing test scripts. |
| Parallel execution | User explicitly requested sequential. |
| Retry logic | User explicitly requested no retry. |
| Retry on timeout | User: "just report". |
| Setup integration in agents | User: reuse existing environment, no setup logic in agents. |
| `inotifywait` for polling | Less portable than `sleep` polling; requires extra package. |

---

## Appendix B: Schema Quick Reference

### `plan.json`
```json
{
  "generated_at": "ISO-8601",
  "source_dir": "relative/path",
  "scenarios": [
    {"name": "...", "path": "...", "timeout_seconds": N, "env": {}}
  ]
}
```

### `queue/001-*.json`
```json
{
  "id": "001",
  "scenario": "...",
  "scenario_dir": "...",
  "timeout_seconds": 1800,
  "env": {},
  "queued_at": "ISO-8601"
}
```

### `results/001.json`
```json
{
  "id": "001",
  "scenario": "...",
  "status": "PASS|FAIL|TIMEOUT|SKIPPED|ERROR",
  "exit_code": 0,
  "duration_seconds": 45,
  "started_at": "ISO-8601",
  "finished_at": "ISO-8601",
  "stdout_file": "...",
  "stderr_file": "...",
  "stdout_tail": "...",
  "stderr_tail": "..."
}
```
