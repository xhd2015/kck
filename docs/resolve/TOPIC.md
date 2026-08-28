---
name: kck/resolve
description: >-
  Resolve a Grok or Codex session id via ancestor process walk or sibling iTerm tab.
---

# resolve

```text
kck grok resolve [OPTIONS]
kck codex resolve [OPTIONS]
```

## Behavior

- Default: walk ancestors from the current process (or `--pid`) to the nearest
  runner (`grok` or `codex` matching the command); session id comes from
  open-file paths only (never cmdline flags).
- `--tab SEL` / `--tab-index N`: resolve from a sibling iTerm2 tab in this window
  (1-based / `next|left|right`, or 0-based index). No wrap at edges.
  - Grok: when a parent session and its child subagent share the tab, the parent
    id wins; unrelated multi-session tabs still refuse.
  - Codex: multiple unrelated Codex sessions on the same tab refuse (no
    parent/child collapse yet).
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

1. Use when you need the current (or sibling-tab) runner id for follow-up ops
   without listing all live hosts.
2. Prefer bare output for scripting: `id=$(kck grok resolve)` or
   `id=$(kck codex resolve --tab 1)`.
3. Misses (no ancestor, no matching runner on tab, flag conflicts) → hard `Error:`.
