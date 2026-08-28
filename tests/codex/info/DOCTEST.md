# `kck codex info`

Thin CLI over `sessions.Info` + `FormatInfoText` + `FormatActiveBlock`.

L2 only — injectable `CodexHome` + `CodexLiveOpts` + `CodexNow`.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — `run.MainWith` with `Args: ["codex","info",…]`.
- **`kck/run`** — parses `--no-pid`, kck help, `Error:` prefix.
- **`sessions.Info` / `Status`** — detail + Active block (File always no).

## Behaviors

- Codex / info help document info + Active / `--no-pid`.
- Missing / unknown id → `Error:`.
- Known session → Session/Title + Active block.
- `--no-pid` → Active PIDs skipped.

## Version

0.0.1

## Decision Tree

```text
codex/info/
├── help/
│   ├── root-lists-info/
│   ├── codex-usage/
│   └── info-usage/
├── missing-session-id/
├── unknown-session/
├── detail-with-active/
└── no-pid/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-info/` | `kck -h` mentions `codex info`. |
| `help/codex-usage/` | `kck codex --help` lists `info`. |
| `help/info-usage/` | info help documents Active / `--no-pid`. |
| `missing-session-id/` | `Error:`. |
| `unknown-session/` | `Error: codex session not found`. |
| `detail-with-active/` | Session + Active + pid. |
| `no-pid/` | Active PIDs skipped. |

## How to Run

```sh
doctest vet ./tests/codex/info
doctest test ./tests/codex/info
```

```go
import (
	"bytes"
	"testing"
	"time"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/agent-pro/pkgs/procresolve"
	"github.com/xhd2015/doctest/session"

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
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:          append([]string(nil), req.Args...),
		Stdout:        &stdout,
		Stderr:        &stderr,
		CodexHome:     req.CodexHome,
		CodexLiveOpts: live,
		CodexNow:      time.Date(2026, 8, 1, 13, 0, 0, 0, time.UTC),
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
