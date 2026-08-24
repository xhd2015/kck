# `kck grok status`

Thin CLI over `sessions.Status` + formatters. Dual-signal liveness plus
session `summary.json` path (`pathfmt.TildeHome` in text; absolute in JSON).

L2 only — injectable `GrokHome` + `GrokLiveOpts`. No live ps/lsof.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — `run.MainWith` with `Args: ["grok","status",…]`.
- **`kck/run`** — parses flags, `Error:` prefix, kck-flavored help.
- **`sessions.Status`** — Find + file-active + optional live PIDs.

## Behaviors

- Root / grok / status help document status + path.
- Missing / unknown id → `Error:`.
- Running / inactive / `--no-pid` states.
- Text includes `Path:` with `summary.json`.
- `--json` includes absolute `path`, no ANSI.

## Version

0.0.1

## Decision Tree

```text
grok/status/
├── help/
│   ├── root-lists-status/
│   ├── grok-usage/
│   └── status-usage/
├── missing-session-id/
├── unknown-session/
├── running/
├── inactive/
├── no-pid/
└── json-path/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-status/` | `kck -h` mentions `grok status`. |
| `help/grok-usage/` | `kck grok --help` lists `status`. |
| `help/status-usage/` | status help documents path + `--json` + `--no-pid`. |
| `missing-session-id/` | `Error:` usage. |
| `unknown-session/` | `Error: grok session not found`. |
| `running/` | State running; Path line; pid present. |
| `inactive/` | State inactive; Path line; PIDs none. |
| `no-pid/` | PIDs skipped; state from file. |
| `json-path/` | absolute `path` ends with summary.json; no ANSI. |

## How to Run

```sh
doctest vet ./tests/grok/status
doctest test ./tests/grok/status
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
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
	GrokHome  string
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
		Args:         append([]string(nil), req.Args...),
		Stdout:       &stdout,
		Stderr:       &stderr,
		GrokHome:     req.GrokHome,
		GrokLiveOpts: live,
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
