# `kck codex wait`

Thin CLI over `codex/sessions.Wait`. Requires live running state, then classifies
turn open/closed from the rollout JSONL (not screen idle).

L2 only — injectable `CodexHome`, `CodexLiveOpts`, `CodexWaitOpts`.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — `run.MainWith` with `Args: ["codex","wait",…]`.
- **`kck/run`** — parses flags, `Error:` prefix, kck-flavored help.
- **`sessions.Wait`** — Status running gate + ScanLinesFromTail + WatchLine.

## Behaviors

- Root / codex / wait help document wait.
- Missing / unknown / not-running → `Error:`.
- Outside turn (last task_complete) while running → immediate success.
- Mid-turn → wait until injected WatchLine delivers task_complete or turn_aborted.
- Timeout while mid-turn → `Error: timeout`.

## Version

0.0.1

## Decision Tree

```text
codex/wait/
├── help/
│   ├── root-lists-wait/
│   ├── codex-usage/
│   └── wait-usage/
├── missing-session-id/
├── unknown-session/
├── not-running/
├── outside-turn/
├── mid-turn-completes/
├── mid-turn-aborted/
└── timeout/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-wait/` | `kck -h` mentions `codex wait`. |
| `help/codex-usage/` | `kck codex --help` lists `wait`. |
| `help/wait-usage/` | wait help documents timeout + rollout semantics. |
| `missing-session-id/` | `Error:` usage. |
| `unknown-session/` | `Error: codex session not found`. |
| `not-running/` | `Error: session not running`. |
| `outside-turn/` | `reason: turn_completed` + session-id. |
| `mid-turn-completes/` | blocks then `reason: turn_completed` on task_complete. |
| `mid-turn-aborted/` | blocks then `reason: turn_completed` on turn_aborted. |
| `timeout/` | `Error: timeout` while mid-turn. |

## How to Run

```sh
doctest vet ./tests/codex/wait
doctest test ./tests/codex/wait
```

```go
import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/agent-pro/pkgs/procresolve"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/logs"

	"kck/run"
)

type FixtureProc struct {
	PID  int
	PPID int
	Cmd  string
}

type Request struct {
	Args      []string
	CodexHome string
	TempDir   string
	SessionID string
	Procs     []FixtureProc
	OpenFiles map[int][]string
	// WatchLine injects Phase B for mid-turn / timeout leaves. nil → production.
	WatchLine func(ctx context.Context, path string, opts logs.WatchLineOptions, fn func(line string) error) error
	// WaitTimeout overrides sessions.Wait timeout when set.
	WaitTimeout time.Duration
}

type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	files := req.OpenFiles
	if files == nil {
		files = map[int][]string{}
	}
	live := &sessions.LiveOptions{
		ListProcs: func() []procresolve.Proc {
			out := make([]procresolve.Proc, 0, len(req.Procs))
			for _, p := range req.Procs {
				out = append(out, procresolve.Proc{PID: p.PID, PPID: p.PPID, Cmd: p.Cmd})
			}
			return out
		},
		Lsof: func(pid int) []string {
			return files[pid]
		},
	}
	waitOpts := &sessions.WaitOpts{
		Live:           live,
		Timeout:        req.WaitTimeout,
		StatusInterval: 20 * time.Millisecond,
		WatchLine:      req.WatchLine,
	}
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:          append([]string(nil), req.Args...),
		Stdout:        &stdout,
		Stderr:        &stderr,
		CodexHome:     req.CodexHome,
		CodexLiveOpts: live,
		CodexWaitOpts: waitOpts,
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
