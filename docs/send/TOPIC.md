---
name: kck/send
description: >-
  Type text and/or keys into the live iTerm pane hosting a Grok or Codex session; optional --open.
---

# send

```text
kck grok send [text] (--session-id <id> | --tab SEL | --tab-index N) [OPTIONS]
kck codex send [text] (--session-id <id> | --tab SEL | --tab-index N) [OPTIONS]
```

Same write path as `kool iterm2 session <uuid> send …` (`iterm2.SendText`).
`--cron` is peeled in kck for both Grok and Codex (agent-pro send unchanged).

## Session source (exactly one)

- `--session-id ID` (payload may be positional or `--text` / key flags, so id is a flag)
- `--tab SEL` / `--tab-index N`

## Key options

| Flag | Effect |
|------|--------|
| `--open` | If no host: resume in a new window, wait up to 120s, then send (bare grok / unmanaged only) |
| `--no-agent-run` | Force iTerm path for `--session-id` (skip agent-run prefer) |
| `--focus` | Focus the tab before typing |
| `--no-ctrl-u` | Do not clear the line with Ctrl-U before typing |
| `--no-submit` | Type without submitting (no trailing Enter) |
| `--index N` | Disambiguate multi-host (`--session-id` only) |
| `--enter` | Append Enter (`\n`) in the send sequence |
| `--up` / `--down` / `--left` / `--right` | Append arrow (CSI) |
| `--esc` | Append Escape |
| `--ctrl-c` / `--ctrl-d` | Append Ctrl-C / Ctrl-D |
| `--text STR` | Append text in sequence order (interleaves with keys) |
| `--cron EXPR` | Foreground repeat send on an easy-cron schedule until done / Ctrl-C |

At least one of `[text]`, `--text`, or a key flag is required. Sequence flags
(`--enter` / `--up` / `--text` / …) keep CLI order and may repeat. Positional
`[text]`, when present, is always appended **last** after the sequence.

Key-only sends force no Ctrl-U and no AppleScript trailing newline. Any
`--enter` in the sequence also disables the trailing newline (Enter is already
in the payload).

`--open` cannot combine with `--tab` / `--tab-index`.

### `--cron` expressions

Examples: `every-1h`, `every-1h-at-4m`, `every-5m-until-19h00m`,
`every-5m-not-within-19h00m-to-06h30m` (see `easycron`).

- Runs in the foreground in this process (no daemon / job store).
- First tick can fire immediately; then sleeps until each next fire.
- Later-tick send failures print `warning:` and continue; the first tick still hard-fails.
- With `--dry-run`: print the next few fire times and one would-send; do not loop.

## Success stdout

- Default: `sent to session <id>` (or tab-oriented wording).
- With `--open` after bare grok resume: two lines — open/resume ack, then sent line.
- Agent-run-managed `--session-id`: one `sent to session …` line via agent-run queue/resume (no iTerm discovery or SendText). `--tab` / `--tab-index` still use iTerm.
- With `--cron`: each tick prints `sent…`, then `next <time> (<expr>)`; ends with `cron done: until reached` when Until expires.

## Agent tips

1. No idle gate — send even if the pane looks busy (user asked to type).
2. Prefer `--session-id` when you already know the Grok id from `kck` / `info`.
3. Use `--open` only when starting/resuming is acceptable for the user.
4. Use `--text` inside the sequence to interleave typed chunks with keys; positional text is always last.
