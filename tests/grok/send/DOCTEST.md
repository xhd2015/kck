# `kck grok send`

Thin CLI over `agent-pro` `sessions.Send`: type text into a live Grok iTerm
host via `iterm2.SendText`. Supports `--session-id` / `--tab` / `--tab-index`,
optional `--open` resume-then-send.

L2 only — injectable `GrokHome` + `GrokSendOpts` (SendFake). No live iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["grok","send",…]`.
- **`kck/run`** — dispatches `grok send`, prints kck-flavored help, prefixes
  errors with `Error:`.
- **`sessions.RunSend`** — core send (agent-pro).

## Behaviors

- Root help mentions `grok send`.
- `kck grok --help` / `kck grok send --help` document session source + flags (including `--cron`).
- Missing text-or-key / session source → `Error:`.
- Key/`--text` sequence flags compose in order; positional text is always last.
- Unknown session → `Error: grok session not found: …`.
- One hosting tab → `sent to session …`.
- No live host → `Error: no hosting iTerm tab …`.
- `--open` resume then send → two stdout lines.
- `--open --tab` → usage error.
- `--tab N` → send to resolved tab.
- Flag opts plumbed into SendText.
- `--cron` peels in kck, foreground loop via `easycron`; agent-pro unchanged.
- Bad `--cron` → `Error: invalid --cron:…` before any send.
- `--cron` + `--dry-run` → schedule preview + one would-send; no loop.

## Version

0.0.1

## Decision Tree

```text
grok/send/
├── help/
│   ├── root-lists-send/
│   ├── grok-usage/
│   └── send-usage/              documents --cron
├── missing-session-source/
├── missing-text/
├── unknown-session/
├── send-exactly-one/
├── no-live-fails/
├── open-resume-then-send/
├── open-with-tab-rejected/
├── tab-send/
├── opts-flags/
├── keys/
│   ├── ctrl-c-only/
│   └── interleave/
└── cron/
    ├── invalid-expr/
    ├── until-two-ticks/
    ├── quiet-skips/
    ├── mid-fail-continues/
    └── dry-run-preview/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-send/` | `kck -h` mentions `grok send`. |
| `help/grok-usage/` | `kck grok --help` lists `send`. |
| `help/send-usage/` | `kck grok send --help` documents `--cron` and key flags. |
| `missing-session-source/` | `Error:` usage. |
| `missing-text/` | `Error:` missing text or key. |
| `unknown-session/` | `Error: grok session not found`. |
| `send-exactly-one/` | `sent to session`; SendText called. |
| `no-live-fails/` | Hard error; SendText not called. |
| `open-resume-then-send/` | Open + sent lines; opener + SendText. |
| `open-with-tab-rejected/` | Usage error. |
| `tab-send/` | `--tab 2` sends. |
| `opts-flags/` | Focus/NoSubmit/NoCtrlU plumbed. |
| `keys/ctrl-c-only/` | `--ctrl-c` → `\x03`; NoCtrlU+NoSubmit. |
| `keys/interleave/` | `--up --text --enter` + positional last. |
| `cron/invalid-expr/` | `Error: invalid --cron:`; no SendText. |
| `cron/until-two-ticks/` | two sends then `cron done: until reached`. |
| `cron/quiet-skips/` | quiet band skipped; fires resume after End. |
| `cron/mid-fail-continues/` | 2nd tick `warning:`; 3rd still sends. |
| `cron/dry-run-preview/` | `cron …` + `next[i]` + would-send; no SendText. |

## How to Run

```sh
doctest vet ./tests/grok/send
doctest test ./tests/grok/send
```

```go
import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"

	"kck/run"
)

type Request struct {
	Args             []string
	GrokHome         string
	TempDir          string
	ProjectDir       string
	SessionID        string
	Procs            []sessions.FocusProc
	OpenFiles        map[int][]string
	ITerm            []iterm2.SessionRef
	CurrentSessionID string
	ControllingTTY   string
	AfterOpenHost    bool

	// CronClock, when non-zero, injects GrokCronNow/Sleep/Loc (UTC) for --cron.
	CronClock time.Time
	// FailSendOnTick makes the Nth SendText call fail (1-based). 0 = never.
	FailSendOnTick int
}

type Response struct {
	Stdout    string
	Stderr    string
	ErrText   string
	ExitCode  int
	SendCalls []sessions.SendCall
	Opened    []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	fake := &sessions.SendFake{
		FocusFake: sessions.FocusFake{
			Procs:     req.Procs,
			OpenFiles: req.OpenFiles,
			ITerm:     req.ITerm,
		},
		CurrentSessionID: req.CurrentSessionID,
		ControllingTTY:   req.ControllingTTY,
	}
	if req.AfterOpenHost {
		sid := req.SessionID
		fake.AfterOpen = func(f *sessions.SendFake) {
			f.Procs = []sessions.FocusProc{
				{PID: 9001, PPID: 1, TTY: "ttys148", Cmd: "/usr/local/bin/grok"},
			}
			f.OpenFiles = map[int][]string{
				9001: {"/Users/fixture/.grok/sessions/%2Ftmp%2Fproj/" + sid + "/events.jsonl"},
			}
			f.ITerm = []iterm2.SessionRef{
				{WindowID: "3", WindowName: "worktrees", TabIndex: 1, SessionID: "w2t1p0", TTY: "/dev/ttys148"},
			}
		}
	}
	var stdout, stderr bytes.Buffer
	sendOpts := fake.SendOpts()
	if req.FailSendOnTick > 0 {
		orig := sendOpts.SendText
		n := 0
		sendOpts.SendText = func(sessionID, text string, opts iterm2.SendTextOptions, cfg *iterm2.SendTextConfig) error {
			n++
			if n == req.FailSendOnTick {
				return fmt.Errorf("no hosting iTerm tab for session %s", sessionID)
			}
			return orig(sessionID, text, opts, cfg)
		}
	}
	runOpts := run.Options{
		Args:         append([]string(nil), req.Args...),
		Stdout:       &stdout,
		Stderr:       &stderr,
		GrokHome:     req.GrokHome,
		GrokSendOpts: sendOpts,
	}
	if !req.CronClock.IsZero() {
		clock := req.CronClock
		loc := time.FixedZone("UTC", 0)
		runOpts.GrokCronLoc = loc
		runOpts.GrokCronNow = func() time.Time { return clock }
		runOpts.GrokCronSleep = func(d time.Duration) error {
			clock = clock.Add(d)
			return nil
		}
	}
	err := run.MainWith(runOpts)
	resp := &Response{
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		SendCalls: append([]sessions.SendCall(nil), fake.SendCalls...),
		Opened:    append([]string(nil), fake.Opened...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
