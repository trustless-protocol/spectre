#!/usr/bin/env python3
"""Generates docs/assets/src/architecture.html and docs/assets/architecture.svg.

Run from the repository root:  python3 docs/assets/src/architecture-gen.py
"""
from pathlib import Path

INK, GREY, LINE, ACCENT, TINT = "#0a0a0a", "#737373", "#d4d4d4", "#ea580c", "#fff7ed"
SANS = "Inter, -apple-system, 'Segoe UI', Helvetica, Arial, sans-serif"
MONO = "'Geist Mono', ui-monospace, SFMono-Regular, Menlo, Consolas, monospace"

W, H = 1100, 500
FRAME_Y, FRAME_H, FRAME_W = 70, 366, 310
BW, BH = 270, 76
ROWS = [152, 244, 336]

out = []


def text(x, y, s, size=13, fill=INK, weight=400, font=SANS, anchor="start", spacing=None):
    extra = f' letter-spacing="{spacing}"' if spacing else ""
    out.append(f'<text x="{x}" y="{y}" font-family="{font}" font-size="{size}" font-weight="{weight}" '
               f'fill="{fill}" text-anchor="{anchor}"{extra}>{s}</text>')


def frame(x, label, sub):
    out.append(f'<rect x="{x}" y="{FRAME_Y}" width="{FRAME_W}" height="{FRAME_H}" fill="none" stroke="{LINE}"/>')
    out.append(f'<rect x="{x + 20}" y="{FRAME_Y + 21}" width="7" height="7" fill="{ACCENT}"/>')
    text(x + 35, FRAME_Y + 29, label, size=11.5, fill=INK, font=MONO, spacing="0.08em", weight=600)
    text(x + 20, FRAME_Y + 56, sub, size=12, fill=GREY)


def box(x, row, title, sub, highlight=False):
    y = ROWS[row]
    stroke, width = (ACCENT, 1.6) if highlight else (INK, 1.2)
    out.append(f'<rect x="{x}" y="{y}" width="{BW}" height="{BH}" fill="#ffffff" stroke="{stroke}" '
               f'stroke-width="{width}"/>')
    text(x + 18, y + 32, title, size=15, weight=600)
    text(x + 18, y + 54, sub, size=12, fill=GREY)


def arrow(x1, x2, y, label):
    out.append(f'<line x1="{x1}" y1="{y}" x2="{x2}" y2="{y}" stroke="{INK}" stroke-width="1.4" '
               f'marker-end="url(#head)"/>')
    text((x1 + x2) / 2, y - 12, label, size=12, fill=ACCENT, font=MONO, anchor="middle", weight=600)


text(30, 40, "ARCHITECTURE · SPECTRE", size=11.5, fill=GREY, font=MONO, spacing="0.14em")

# ---- Cosmos ------------------------------------------------------------------------------------------
frame(30, "COSMOS · ON-CHAIN", "standard ibc-go")
box(50, 0, "ibc-go router", "ICS-26 · packet lifecycle")
box(50, 1, "08-wasm client", "ICS-08 · Ethereum finality")
box(50, 2, "transfer (ICS-20)", "ICS-20 · token transfer")

# ---- Ethereum ----------------------------------------------------------------------------------------
frame(760, "ETHEREUM · ON-CHAIN", "Solidity contracts")
box(780, 0, "ICS26Router", "ICS-26 · packet lifecycle")
box(780, 1, "Spectre light client", "modified ICS-07 · Groth16 signature check", highlight=True)
box(780, 2, "ICS20Transfer", "ICS-20 · Cosmos tokens ↔ ERC-20")

# ---- relayer, between the two chains -----------------------------------------------------------------
rx, ry, rw, rh = 400, 240, 300, 84
out.append(f'<rect x="{rx}" y="{ry}" width="{rw}" height="{rh}" fill="{TINT}" stroke="{ACCENT}" stroke-width="1.6"/>')
text(rx + rw / 2, ry + 36, "Relayer + Groth16 prover", size=15, weight=600, anchor="middle")
text(rx + rw / 2, ry + 58, "off-chain · packets + proofs", size=12, fill=GREY, anchor="middle")

arrow(340 + 6, 760 - 8, 204, "Cosmos → Ethereum")
arrow(760 - 6, 340 + 8, 366, "Ethereum → Cosmos")

# ---- note --------------------------------------------------------------------------------------------
out.append(f'<line x1="30" y1="456" x2="{W - 30}" y2="456" stroke="{LINE}" stroke-dasharray="5 4"/>')
out.append(f'<rect x="30" y="471" width="12" height="12" fill="#ffffff" stroke="{ACCENT}" stroke-width="1.6"/>')
text(52, 482, "Only the light client is re-designed (highlighted): it proves Cosmos → Ethereum consensus in ZK.",
     size=12.5, fill=GREY)

svg = (f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {W} {H}" width="{W}" height="{H}" role="img" '
       f'aria-label="Spectre architecture">\n'
       f'<defs><marker id="head" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="8" markerHeight="8" '
       f'orient="auto-start-reverse"><path d="M0 0L10 5L0 10z" fill="{INK}"/></marker></defs>\n'
       f'<rect width="{W}" height="{H}" fill="#ffffff"/>\n' + "\n".join(out) + "\n</svg>\n")

here = Path(__file__).resolve().parent
(here.parent / "architecture.svg").write_text(svg)
(here / "architecture.html").write_text(
    '<!doctype html><meta charset="utf-8"><title>Spectre architecture</title>'
    '<body style="margin:0;background:#fff">\n' + svg + "</body>\n")
print("wrote", here.parent / "architecture.svg")
