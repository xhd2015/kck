---
name: kck/pickup
description: >-
  Open a new empty Grok/Codex session staged with a kck-pickup-a-session draft
  from a base session (read from the bottom + new instruction; not a native fork).
---

# pickup

```text
kck grok pickup "msg..." (--session-id ID | --tab SEL | --tab-index N) [OPTIONS]
kck codex pickup "msg..." (--session-id ID | --tab SEL | --tab-index N) [OPTIONS]
```

Opens a **new empty** session and stages (does **not** submit) a draft:

```text
read ~/.cache/kck-pickup-a-session/SKILL.md, session-id: <base-id>, <msg...>
```

By default runs in the **current terminal**. Pass `--new-window` to ForceNew an
iTerm window. `--here` forces the default (current terminal). `--here` and
`--new-window` cannot be combined.

On each run (including `--dry-run`), kck hydrates the embedded agent skill
`kck-pickup-a-session` to `~/.cache/kck-pickup-a-session/SKILL.md`, skipping
the write when the on-disk MD5 already matches the binary. The draft uses
`pathfmt.TildeHome` for the skill path.

Spirit: pick up from the base session (read from the bottom / where it left
off) and continue with the new instruction. Not a native
`grok --fork-session` / transcript clone — the base id is a reference the new
agent fetches via `kck grok|codex messages`.

## Session source (exactly one)

- `--session-id ID` (message is positional, so id is a flag — same as `send`)
- `--tab SEL` / `--tab-index N`

## Key options

| Flag | Effect |
|------|--------|
| `--dir DIR` | Workspace for the new session (default: base session cwd) |
| `--here` | Current terminal (default; explicit) |
| `-n` / `--new-window` / `--new-terminal` | New iTerm2 window |
| `--no-agent-run` | Bare `grok`/`codex` + stage draft without submit |
| `--dry-run` | Resolve + hydrate skill cache + print plan; do not launch |

Default launch prefers `agent-run run --open --no-submit --agent-runner=…`.

## Success stdout

```text
opened: here; pickup from session <id>
opened: new window; pickup from session <id>
```

## Agent tips

1. Use when you want a fresh session that still knows the base id + intent.
2. Prefer `--session-id` when you already have the base id from `kck` / `info`.
3. From inside a live agent pane, prefer `--new-window` so you do not replace that TTY.
4. Review the staged draft in the composer, then press Enter.
5. Agent skill body: `kck skill --show kck-pickup-a-session`.
