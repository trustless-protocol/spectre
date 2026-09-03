#!/usr/bin/env python3
"""Fail closed when the Solidity refactor drifts from the frozen baseline."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
import sys
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[2]


def canonical_hash(value: Any) -> str:
    encoded = json.dumps(value, sort_keys=True, separators=(",", ":")).encode()
    return hashlib.sha256(encoded).hexdigest()


def file_hash(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def dotted(value: dict[str, Any], key: str) -> Any:
    current: Any = value
    for part in key.split("."):
        current = current.get(part) if isinstance(current, dict) else None
    return current


def load(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text())


def artifact(contract: str) -> dict[str, Any]:
    path = ROOT / "out" / f"{contract}.sol" / f"{contract}.json"
    if not path.is_file():
        raise FileNotFoundError(f"missing Foundry artifact: {path.relative_to(ROOT)}")
    return load(path)


def normalize_tooling_text(value: str, renames: list[dict[str, str]]) -> str:
    for rename in sorted(renames, key=lambda item: len(item["target"]), reverse=True):
        value = value.replace(rename["target"], rename["current"])
    return value


def normalize_tooling_value(value: Any, renames: list[dict[str, str]]) -> Any:
    if isinstance(value, str):
        return normalize_tooling_text(value, renames)
    if isinstance(value, list):
        return [normalize_tooling_value(item, renames) for item in value]
    if isinstance(value, dict):
        return {key: normalize_tooling_value(item, renames) for key, item in value.items()}
    return value


def artifact_record(data: dict[str, Any], renames: list[dict[str, str]]) -> dict[str, Any]:
    runtime = dotted(data, "deployedBytecode.object") or ""
    initcode = dotted(data, "bytecode.object") or ""
    return {
        "abi_sha256": canonical_hash(normalize_tooling_value(data.get("abi"), renames)),
        "runtime_sha256": hashlib.sha256(runtime.encode()).hexdigest(),
        "initcode_sha256": hashlib.sha256(initcode.encode()).hexdigest(),
        "runtime_bytes": max(0, (len(runtime) - 2) // 2),
        "initcode_bytes": max(0, (len(initcode) - 2) // 2),
    }


def normalized_storage_record(
    root: Path,
    entry: dict[str, Any],
    renames: list[dict[str, str]],
    source_key: str = "source",
) -> dict[str, Any]:
    source_path = root / entry[source_key]
    source = source_path.read_text()
    structs: dict[str, str] = {}
    for struct_name in entry["structs"]:
        match = re.search(
            rf"\bstruct\s+{re.escape(struct_name)}\s*\{{(.*?)\n\s*\}}",
            source,
            flags=re.DOTALL,
        )
        if match is None:
            raise ValueError(f"missing storage struct {struct_name} in {entry[source_key]}")
        body = re.sub(r"/\*.*?\*/", " ", match.group(1), flags=re.DOTALL)
        body = re.sub(r"//[^\n]*", " ", body)
        normalized_body = re.sub(r"\s+", " ", body).strip()
        structs[struct_name] = normalize_tooling_text(normalized_body, renames)

    slots: dict[str, str] = {}
    for slot_name in entry["slots"]:
        match = re.search(
            rf"bytes32\s+(?:(?:private|internal|public)\s+)?constant\s+"
            rf"{re.escape(slot_name)}\s*=\s*(0x[0-9a-fA-F]{{64}})\s*;",
            source,
            flags=re.DOTALL,
        )
        if match is None:
            raise ValueError(f"missing storage slot {slot_name} in {entry[source_key]}")
        slots[slot_name] = match.group(1).lower()

    namespaces = sorted(set(re.findall(r"@custom:storage-location\s+([^\s*]+)", source)))
    return {"structs": structs, "slots": slots, "namespaces": namespaces}


def gas_results(contract: str) -> dict[str, int]:
    completed = subprocess.run(
        ["forge", "test", "--match-contract", contract, "--json"],
        cwd=ROOT,
        check=True,
        text=True,
        capture_output=True,
    )
    suites = json.loads(completed.stdout)
    matching = [
        suite
        for suite_name, suite in suites.items()
        if suite_name.endswith(f":{contract}")
    ]
    if len(matching) != 1:
        raise ValueError(f"expected one gas suite for {contract}, found {len(matching)}")
    return {
        name: result["kind"]["Unit"]["gas"]
        for name, result in matching[0]["test_results"].items()
    }


def compare(args: argparse.Namespace) -> list[str]:
    failures: list[str] = []
    baseline = load(args.baseline)
    contract_manifest = load(args.contract_manifest)
    verifier_manifest = load(args.verifier_manifest)
    storage_manifest = load(args.storage_manifest)
    test_manifest = load(args.test_manifest)
    tooling_rename_manifest = load(args.tooling_rename_manifest)
    toolchain_manifest = load(args.toolchain_manifest)
    tooling_renames = tooling_rename_manifest["renames"]
    approved_generated_digests = tooling_rename_manifest["generated_digests"]
    approved_consumer_digests = tooling_rename_manifest.get("consumer_digests", {})

    source_by_contract = {
        item["name"]: item["source"] for item in contract_manifest["contracts"]
    }
    for contract, expected in baseline["artifact_contracts"].items():
        try:
            data = artifact(contract)
        except FileNotFoundError as error:
            failures.append(str(error))
            continue
        actual = artifact_record(data, tooling_renames)
        for key, expected_value in expected.items():
            if actual[key] != expected_value:
                failures.append(
                    f"{contract} {key}: expected {expected_value}, got {actual[key]}"
                )
        if actual["runtime_bytes"] > 24_576:
            failures.append(f"{contract} exceeds EIP-170: {actual['runtime_bytes']} bytes")
        source = source_by_contract.get(contract)
        if source and source not in data.get("metadata", {}).get("sources", {}):
            failures.append(f"{contract} artifact does not contain manifest source {source}")

    for entry in storage_manifest["entries"]:
        try:
            actual_storage = normalized_storage_record(ROOT, entry, tooling_renames)
        except (OSError, ValueError) as error:
            failures.append(str(error))
            continue
        expected_storage = baseline["storage_namespaces"].get(entry["name"])
        if actual_storage != expected_storage:
            failures.append(f"{entry['name']} ERC-7201 namespace or struct layout changed")

    for relative, expected_hash in baseline["digests"].items():
        path = ROOT / relative
        if not path.is_file():
            failures.append(f"missing locked fixture or binding: {relative}")
        else:
            actual_hash = file_hash(path)
            approved_hash = approved_generated_digests.get(relative)
            locked_hash = approved_hash if approved_hash is not None else expected_hash
            if actual_hash != locked_hash:
                label = "approved tooling artifact" if approved_hash is not None else "locked fixture or binding"
                failures.append(f"{label} changed: {relative}")

    unknown_approved_paths = sorted(set(approved_generated_digests) - set(baseline["digests"]))
    for relative in unknown_approved_paths:
        failures.append(f"tooling rename manifest approves an unlocked artifact: {relative}")

    for relative, expected_hash in approved_consumer_digests.items():
        path = ROOT / relative
        if not path.is_file():
            failures.append(f"missing coordinated tooling consumer: {relative}")
        elif file_hash(path) != expected_hash:
            failures.append(f"coordinated tooling consumer changed: {relative}")

    for rename in tooling_rename_manifest.get("source_module_renames", []):
        current = ROOT / rename["current"]
        target = ROOT / rename["target"]
        if current.exists():
            failures.append(f"stale tooling source module still exists: {rename['current']}")
        if not target.is_file():
            failures.append(f"missing renamed tooling source module: {rename['target']}")
        for consumer in rename.get("consumers", []):
            if not (ROOT / consumer).is_file():
                failures.append(f"missing tooling source consumer: {consumer}")

    for artifact_rename in tooling_rename_manifest.get("retired_generated_artifacts", []):
        current = ROOT / artifact_rename["current"]
        replacement = ROOT / artifact_rename["replacement"]
        if current.exists():
            failures.append(f"stale generated tooling artifact still exists: {artifact_rename['current']}")
        if not replacement.is_file():
            failures.append(f"missing generated tooling replacement: {artifact_rename['replacement']}")

    for relative, expected_hash in toolchain_manifest.get("locks", {}).items():
        path = ROOT / relative
        if not path.is_file():
            failures.append(f"missing locked toolchain input: {relative}")
        elif file_hash(path) != expected_hash:
            failures.append(f"locked toolchain input changed: {relative}")

    for relative, expected in baseline.get("semantic_artifact_digests", {}).items():
        path = ROOT / relative
        if not path.is_file():
            failures.append(f"missing semantic artifact: {relative}")
            continue
        data = load(path)
        semantic_fields = {
            "abi": normalize_tooling_value(data.get("abi"), tooling_renames),
            "initcode": data.get("bytecode", {}).get("object"),
            "runtime": data.get("deployedBytecode", {}).get("object"),
        }
        for field, value in semantic_fields.items():
            if canonical_hash(value) != expected[f"{field}_sha256"]:
                failures.append(f"semantic artifact field changed: {relative} {field}")

    tests = sorted(
        path.relative_to(ROOT).as_posix() for path in (ROOT / "test").rglob("*.t.sol")
    )
    expected_tests = sorted(item["path"] for item in test_manifest["tests"])
    if tests != expected_tests:
        failures.append("Foundry test inventory differs from test-manifest.json")
    if len(tests) != baseline["test_count"]:
        failures.append(f"expected {baseline['test_count']} test contracts, found {len(tests)}")

    if verifier_manifest.get("scope") != "checked-local-generator":
        failures.append("verifier manifest must identify the checked-local-generator scope")
    if verifier_manifest.get("supported_buckets") != [4]:
        failures.append("checked local generator buckets must be exactly [4]")
    for item in verifier_manifest.get("artifacts", {}).values():
        path = ROOT / item["path"]
        if not path.is_file():
            failures.append(f"missing paired verifier artifact: {item['path']}")
            continue
        expected_hash = item.get("sha256")
        if expected_hash is not None and file_hash(path) != expected_hash:
            failures.append(f"deterministic verifier artifact changed: {item['path']}")

    if not failures:
        gas_baseline = baseline.get("gas_baseline", {})
        max_regression_percent = gas_baseline.get("max_regression_percent", 5)
        for contract, expected_tests in gas_baseline.get("contracts", {}).items():
            try:
                actual_tests = gas_results(contract)
            except (json.JSONDecodeError, subprocess.CalledProcessError, ValueError) as error:
                failures.append(f"unable to measure {contract} gas: {error}")
                continue
            for test_name, expected_gas in expected_tests.items():
                actual_gas = actual_tests.get(test_name)
                if actual_gas is None:
                    failures.append(f"missing locked gas test: {contract}.{test_name}")
                elif actual_gas * 100 > expected_gas * (100 + max_regression_percent):
                    failures.append(
                        f"{contract}.{test_name} gas regressed more than "
                        f"{max_regression_percent}%: expected {expected_gas}, got {actual_gas}"
                    )

    baseline_commit = baseline.get("baseline_commit")
    protected_commands = [
        ["git", "diff", "--name-only", "--", "docs/refactor"],
        ["git", "diff", "--cached", "--name-only", "--", "docs/refactor"],
    ]
    if baseline_commit:
        protected_commands.append(
            ["git", "diff", "--name-only", f"{baseline_commit}..HEAD", "--", "docs/refactor"]
        )
    protected = set()
    for command in protected_commands:
        protected.update(
            subprocess.run(
                command,
                cwd=ROOT,
                check=True,
                text=True,
                capture_output=True,
            ).stdout.splitlines()
        )
    if protected:
        failures.append(f"protected docs/refactor files changed: {', '.join(sorted(protected))}")
    return failures


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--baseline",
        type=Path,
        default=ROOT / "scripts/solidity/baseline.json",
    )
    parser.add_argument(
        "--contract-manifest",
        type=Path,
        default=ROOT / "scripts/solidity/contracts.json",
    )
    parser.add_argument(
        "--verifier-manifest",
        type=Path,
        default=ROOT / "scripts/solidity/verifier-manifest.json",
    )
    parser.add_argument(
        "--storage-manifest",
        type=Path,
        default=ROOT / "scripts/solidity/storage-manifest.json",
    )
    parser.add_argument(
        "--test-manifest",
        type=Path,
        default=ROOT / "scripts/solidity/test-manifest.json",
    )
    parser.add_argument(
        "--tooling-rename-manifest",
        type=Path,
        default=ROOT / "scripts/solidity/tooling-rename-manifest.json",
    )
    parser.add_argument(
        "--toolchain-manifest",
        type=Path,
        default=ROOT / "scripts/solidity/toolchain-manifest.json",
    )
    args = parser.parse_args()
    failures = compare(args)
    if failures:
        print("Solidity compatibility check failed:", file=sys.stderr)
        for failure in failures:
            print(f"- {failure}", file=sys.stderr)
        return 1
    print(
        "Solidity compatibility check passed: ABI, storage, bytecode, gas, "
        "fixtures, bindings, tests, verifier artifacts, and protected docs match."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
