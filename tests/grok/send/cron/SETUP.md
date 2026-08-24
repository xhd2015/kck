# Scenario

**Feature**: `kck grok send --cron` foreground easy-cron loop

```
peel --cron → easycron.Parse → loop Next/sleep/RunSend (or dry-run preview)
```

## Preconditions

- Parent fixture home / SendFake apply.
- Leaves that exercise the loop set `req.CronClock` (UTC) so sleeps are virtual.
- No live iTerm; use `writeKckSendSession` + live host helpers when a send is expected.

## Steps

1. Leaf sets `req.Args` including `--cron EXPR` and session source.
2. Optional: `req.CronClock`, `req.FailSendOnTick`.
3. Root `Run` injects cron clock/sleep when `CronClock` is set.
