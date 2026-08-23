#!/usr/bin/env python3
"""Validate Solidity ownership, import direction, and stale source paths."""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
MANIFEST = ROOT / "scripts/solidity-refactor/ownership-manifest.json"
IMPORT = re.compile(r"""(?:from\s+|import\s+)["']([^"']+)["']""")
DECLARATION = re.compile(
    r"^\s*(?:abstract\s+)?(?:contract|interface|library)\s+([A-Za-z_][A-Za-z0-9_]*)",
    re.MULTILINE,
)


def owner(path: str) -> str | None:
    if path.startswith("contracts/shared/"):
        return "shared"
    if path.startswith("contracts/light-clients/interfaces/") or path.startswith(
        "contracts/light-clients/messages/"
    ):
        return "compatibility"
    if path.startswith("contracts/light-clients/spectre/"):
        return "spectre"
    if path.startswith("contracts/core/"):
        return "core"
    if path.startswith("contracts/apps/ics20/"):
        return "ics20"
    if path.startswith("contracts/periphery/"):
        return "periphery"
    if path.startswith("contracts/verifiers/"):
        return "generated"
    return None


def import_allowed(source: str, target: str) -> bool:
    source_owner = owner(source)
    target_owner = owner(target)
    if target_owner is None or source_owner == target_owner:
        return True
    allowed = {
        "shared": {"shared"},
        "compatibility": {"compatibility", "core", "spectre"},
        "core": {"core", "shared", "compatibility"},
        "ics20": {"ics20", "shared", "core"},
        "spectre": {"spectre", "shared", "compatibility", "core", "generated"},
        "periphery": {
            "periphery",
            "shared",
            "core",
            "ics20",
            "spectre",
            "compatibility",
        },
        "generated": {"generated"},
    }
    if target_owner not in allowed.get(source_owner or "", set()):
        return False
    if source_owner == "ics20" and target_owner == "core":
        return (
            "/interfaces/" in target
            or "/messages/" in target
            or "/compatibility/" in target
        )
    if source_owner == "spectre" and target_owner == "core":
        return "/messages/" in target or "/client-registry/migration/" in target
    if source_owner == "compatibility" and target_owner in {"core", "spectre"}:
        return "/messages/" in target
    return True


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", type=Path, default=MANIFEST)
    args = parser.parse_args()
    failures: list[str] = []
    manifest = json.loads(args.manifest.read_text())
    for item in manifest["entries"]:
        target = ROOT / item["target"]
        if not target.is_file():
            failures.append(f"missing owned target: {item['target']}")
        current = item.get("current")
        if current and current != item["target"] and (ROOT / current).exists():
            failures.append(f"stale source path still exists: {current}")

    for stale_dir in (
        "contracts/utils",
        "contracts/interfaces",
        "contracts/msgs",
        "contracts/errors",
        "contracts/core/client",
        "contracts/shared/bytes",
        "contracts/shared/encoding",
        "contracts/periphery/access",
    ):
        path = ROOT / stale_dir
        if path.exists() and any(path.rglob("*.sol")):
            failures.append(f"stale global Solidity directory is not empty: {stale_dir}")

    for path in sorted((ROOT / "contracts").rglob("*.sol")):
        relative = path.relative_to(ROOT).as_posix()
        if relative.startswith("contracts/verifiers/"):
            continue
        text = path.read_text()
        declarations = DECLARATION.findall(text)
        if path.stem not in declarations:
            failures.append(
                f"file/primary-symbol mismatch: {relative} declares "
                f"{', '.join(declarations) or 'none'}"
            )
        for imported in IMPORT.findall(text):
            if imported.startswith("."):
                failures.append(
                    f"relative repository import is forbidden: {relative} -> {imported}"
                )
                continue
            if imported.startswith(("test/", "scripts/")):
                failures.append(
                    f"production source imports non-production code: {relative} -> {imported}"
                )
            if imported.startswith("contracts/"):
                if not (ROOT / imported).is_file():
                    failures.append(f"unresolved contract import: {relative} -> {imported}")
                elif not import_allowed(relative, imported):
                    failures.append(f"disallowed ownership edge: {relative} -> {imported}")

    old_paths = {
        item["current"]
        for item in manifest["entries"]
        if item.get("current") and item["current"] != item["target"]
    }
    legacy_symbols = {
        "IEscrowErrors",
        "IIBCERC20Errors",
        "IICS20Errors",
        "IRateLimitErrors",
        "IICS02ClientErrors",
        "IICS24HostErrors",
        "IICS26RouterErrors",
        "ISpectreClientErrors",
        "IICS20TransferMsgs",
        "IICS02ClientMsgs",
        "IICS26RouterMsgs",
        "ILightClientMsgs",
        "IGroth16Msgs",
        "IICS07TendermintMsgs",
        "IMembershipMsgs",
        "ISpectreClientMsgs",
    }
    source_roots = (
        ROOT / "contracts",
        ROOT / "test",
        ROOT / "scripts",
        ROOT / "packages",
        ROOT / "e2e",
        ROOT / "relayer",
    )
    for source_root in source_roots:
        for path in source_root.rglob("*"):
            if not path.is_file() or path.suffix not in {".sol", ".sh", ".go", ".rs"}:
                continue
            try:
                text = path.read_text()
            except UnicodeDecodeError:
                continue
            for stale in old_paths:
                if stale in text:
                    failures.append(
                        f"stale source path in {path.relative_to(ROOT)}: {stale}"
                    )
            for symbol in legacy_symbols:
                if re.search(rf"\b{re.escape(symbol)}\b", text):
                    failures.append(
                        f"stale canonical symbol in {path.relative_to(ROOT)}: {symbol}"
                    )

    if failures:
        print("Solidity architecture check failed:", file=sys.stderr)
        for failure in sorted(set(failures)):
            print(f"- {failure}", file=sys.stderr)
        return 1
    print(
        "Solidity architecture check passed: ownership targets, primary symbols, "
        "imports, and stale paths are clean."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
