---
name: kck/resolve
description: >-
  Resolve a Grok session id via ancestor process walk or sibling iTerm tab.
---

# resolve

```text
kck grok resolve [OPTIONS]
```

## Behavior

- Default: walk ancestors from the current process (or `--pid`) to the nearest
  grok runner; session id comes from open-file paths only (never cmdline flags).
- `--tab SEL` / `--tab-index N`: resolve from a sibling iTerm2 tab in this window
  (1-based / `next|left|right`, or 0-based index). No wrap at edges.
- Default success stdout: bare `<session-id>` (script-friendly).
- `--dry-run`: same discovery; `[dry-run]` plan lines instead of bare id.
- `-v`: bare id on stdout; detail fields on stderr.
- `--json`: indented JSON of detail fields (no bare-id line).

## Options

- `--pid PID` — start pid for ancestor walk (default: current process)
- `--tab SEL` — 1-based tab index, or `next|left|right` (`right` ≡ `next`)
- `--tab-index N` — 0-based tab index in this iTerm window
- `--dry-run` / `-v` / `--json` — plan / verbose / machine-readable

Exactly one session source: ancestor walk (default / `--pid`), or `--tab`, or
`--tab-index`. `--pid` cannot combine with tab flags.

## Agent tips

1. Use when you need the current (or sibling-tab) Grok id for `send` / `open` /
   `messages` without listing all live hosts.
2. Prefer bare output for scripting: `id=$(kck grok resolve)`.
3. Misses (no ancestor, no grok on tab, flag conflicts) → hard `Error:`.
