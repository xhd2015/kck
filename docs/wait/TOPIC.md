---
name: kck/wait
description: >-
  Block until the current Grok or Codex turn finishes (updates.jsonl / rollout),
  or error if the session is not running. Prefer over ad-hoc monitor scripts.
---

# wait

```text
kck grok wait <session-id> [OPTIONS]
kck codex wait <session-id> [OPTIONS]
```

## Behavior

- Requires status **running** (live PIDs). Not running → `Error:`.
- **Grok** turn state from `updates.jsonl`: open after `user_message_chunk`,
  closed on `turn_completed` (no `turn_started` on the wire).
- **Codex** turn state from the rollout JSONL (`Status.Path`): open after
  `event_msg` `task_started`, closed on `task_complete` or `turn_aborted`.
- **Codex readiness:** right after `new`, the id may exist via
  `thread-writer-locks/<id>.lock` before the rollout file. Wait races
  rollout create (fsnotify), lock release (flock), and process-tree exit
  (kqueue `NOTE_EXIT` on lock holders + parents) so closing the iTerm window
  aborts ASAP. Still no rollout → `Error: session never created`. No PID
  polling loops.
- Already outside a turn → returns immediately with `reason: turn_completed`
  or `reason: outside_turn`.
- Mid-turn → watches new lines until a close event or `--timeout`.
- Not screen/TTY idle. No `--accept-idle` / `--until` / `--json`.

## Options

- `--timeout DUR` — max wait (default `30m`)

## Agent tips

1. After `kck grok new` / `kck codex new`, poll with `wait <id>` then `messages <id>`.
2. Do not use provider `monitor` + message-hash bash for turn completion.
