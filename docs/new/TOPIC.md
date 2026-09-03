---
name: kck/new
description: >-
  Open a new empty Grok/Codex session via agent-run (wrk create agent-launch
  shape): brainstorm-prefixed prompt, grok-tty/codex-tty, default new terminal.
---

# new

```text
kck grok new "prompt..." [OPTIONS]
kck codex new "prompt..." [OPTIONS]
```

Opens a **new empty** session by wrapping `agent-run run --open …`.

Prompt wrap (wrk create defaults):

- Grok: `/brainstorm <msg>` with `--agent-runner grok-tty`
- Codex: `$brainstorm <msg>` with `--agent-runner codex-tty`

Default launch stages the prompt (`--no-submit`). Pass `--submit` to auto-submit.

## Placement

| Flags | Behavior |
|-------|----------|
| (default) / `--new-terminal` | New iTerm2 window; wait for **provider** (Grok/Codex) session id via registry PID → open-files; print `opened:` + `session-id:`; exit |
| `--here` / `--no-new-terminal` | Current terminal; **silent** (no kck stdout/stderr — protects TUI) |

`--here`/`--no-new-terminal` and `--new-terminal` cannot be combined.

## Key options

| Flag | Effect |
|------|--------|
| `--dir DIR` | Workspace (default: cwd) |
| `--here` / `--no-new-terminal` | Current terminal (silent) |
| `--new-terminal` | New iTerm2 window (default) |
| `--submit` | Auto-submit (omit `--no-submit`) |
| `--dry-run` | Print plan; do not launch |

## Success stdout

New terminal (`session-id` is the **Grok/Codex** provider id, not the agent-run slug):

```text
opened: new terminal; new grok session
session-id: 01a064d2-70ec-7162-b36b-8a50ba323569
```

Wait uses `agentrunapi.WaitProviderSessionID`: agent-run registry/`tty.json` →
`command_pid`/`pid` → `procresolve.ResolveFromPID`.

Here: no kck output.

## Agent tips

1. Prefer default new-terminal when calling from inside a live agent pane.
2. Use `--here` only when you intentionally replace the current TTY.
3. Distinct from `pickup`: `new` has no base session / skill draft.
