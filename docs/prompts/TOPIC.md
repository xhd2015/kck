---
name: kck/prompts
description: >-
  List user prompts from Grok/Codex sessions (compact lines), with --first,
  --main, repeatable --grep (AND), and live host scopes (--this-window /
  --this-space / --tab).
---

# prompts

```text
kck grok prompts [OPTIONS]
kck codex prompts [OPTIONS]
```

## Behavior

- Prints **user prompts only** as `[YYYY-MM-DD HH:MM:SS] text…` lines.
- Thin wrapper over agent-pro `RunPrompts` (Grok: `updates.jsonl`; Codex: rollout user messages).
- Bare multi: last 10 sessions with prompts (same matrix as agent-pro).
- `--first`: only the first surviving prompt per session.
- `--main` (alias `--main-agent`): Grok only — keep main-agent class sessions
  (skip `subagent` / `subagent_resume` / `subagent_fork`, and parent-linked empty kind).
- `--grep P` repeatable: AND, case-insensitive literal.
- `--this-tab` / `--tab current`: current iTerm tab’s session.
- `--this-window`: live hosts in this iTerm window (no default session cap).
- `--this-space`: live hosts on this macOS Mission Control desktop.
- Tab parse/select SSOT: `dot-pkgs/.../iterm2/tabselect`.

## Useful flags

- `--session-id` / positional id / `--tab` / `--tab-index` / `--this-tab`
- `--this-window` / `--this-space`
- `--first`, `--main`, `--head` / `--tail`, `--exclude`, `--max-body`
- `--recent Nd|Nh|Nm`, `--limit N`
- `--color` / `--no-color`

## Agent tips

1. Orient: `kck grok prompts --first --main --max-body 80`.
2. Scope to this window: `kck grok prompts --this-window --first --main`.
3. Narrow: `kck grok prompts --grep foo --grep bar`.
