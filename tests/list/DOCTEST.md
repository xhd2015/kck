# kck list — default session table (slice 1)

Classic TDD for **`kck` list mode**: with no mode flags, list all agent-run
sessions (newest first) with live status columns. L2 in-process only —
injectable home, writers, probe, and iTerm list. No real TTY, osascript, or
process env mutation. (Send lives under `kck grok send`, not list flags.)

```text
kck [OPTIONS]
  (default) list sessions under agent-run home
```

# DSN (Domain Specific Notion)

**Participants**

- **Caller** — user or harness; passes argv without program name; injects
  `Home`, `Stdout`/`Stderr`, `Probe`, `ListITerm` via `run.MainWith`.
- **`kck/run`** — CLI parse (less-flags) + list / help dispatch.
  Production `Main` wires `os.Stdout`/`os.Stderr` and real probe/iTerm;
  tests use **`MainWith(Options)`**.
- **Session store** — flat `$home/sessions/<session_id>/meta.json` (agent-run
  FileStore layout). List reads all sessions; newest first by `updated_at`.
- **Probe** — serial per-session live status: live / sendable / state / reason /
  exited / TTY. Injectable; production talks to TTY registry + process.
- **iTerm list** — injectable `ListITerm`; match session TTY →
  `w=<WindowID> t=<TabIndex>` with multi-match `(+N)` or `-` when none;
  `--no-iterm` skips; list failure → `warning:` + `-` for all.

**Behaviors**

- Default (no mode flags): list all sessions newest-first; human columns
  `SESSION_ID RUNNER LIVE SENDABLE STATE REASON ITERM UPDATED WORKSPACE`;
  footer `N sessions · M needs attention · K sendable`.
- `needs_attention` (D1) = live + not sendable + not exited; show REASON;
  filterable with `--needs-confirm`.
- `--sendable` filters idle writable rows; `--json` machine list without ANSI.
- `--home PATH` isolates store (no `os.Setenv` / `t.Setenv` for home).
- Unreachable TTY: mark LIVE not alive, optional `warning:`, continue (exit 0).
- Help documents list flags; unknown flag / extra positionals → Error.

### Public API (Classic TDD — locked for implementer)

```go
// Package: kck/run

// ProbeResult is one session's live classification for the list row.
type ProbeResult struct {
    Live     bool   // terminal reachable / process alive
    Sendable bool   // idle writable (can accept send)
    State    string // e.g. idle, running, exited, unknown
    Reason   string // needs_attention reason when applicable; else empty
    Exited   bool   // clearly exited → not needs_attention
    TTY      string // device path for iTerm match; empty → no match
}

// SessionMetaView is the meta subset list uses (from store).
type SessionMetaView struct {
    SessionID string
    Runner    string
    Status    string
    Workspace string
    UpdatedAt string
    CreatedAt string
}

// ProbeFunc is invoked once per session, serially. nil → production probe.
type ProbeFunc func(sessionID string, meta SessionMetaView) (ProbeResult, error)

// ITermSession is one iTerm window/tab for column resolution.
type ITermSession struct {
    WindowID string
    TabIndex int // 1-based
    TTY      string
}

// ListITermFunc lists iTerm sessions. nil → production. error → soft-fail.
type ListITermFunc func() ([]ITermSession, error)

// Options drives MainWith (L2 injectable entry).
type Options struct {
    Args      []string // argv after program name
    Home      string   // explicit home; --home in Args wins when both set
    Stdout    io.Writer
    Stderr    io.Writer
    Probe     ProbeFunc
    ListITerm ListITermFunc
    Now       time.Time // zero → time.Now(); for UPDATED ages if relative
}

// Main is production CLI: MainWith(Options{Args, Stdout:os.Stdout, Stderr:os.Stderr}).
func Main(args []string) error

// MainWith runs kck with injectable IO / home / probe / iTerm (L2).
// Success (list, help) → nil. Failures (bad flags) → non-nil error after
// writing user-facing "Error: …" lines to Stderr when appropriate.
func MainWith(opts Options) error
```

**Stable messages (assertable):**

| Case | Message |
|------|---------|
| Footer | `{N} sessions · {M} needs attention · {K} sendable` |
| ITERM single | `w=<WindowID> t=<TabIndex>` |
| ITERM multi | `w=… t=…(+N)` where N = extra matches beyond first |
| ITERM none / `--no-iterm` / soft-fail | `-` |
| LIVE yes/no | `yes` / `no` |
| SENDABLE yes/no | `yes` / `no` |

## Version

0.0.2

## Decision Tree

```
tests/list/
├── DOCTEST.md
├── SETUP.md
├── help/
│   └── documents-flags/              -h/--help lists home,json,filters
├── usage-errors/
│   ├── unknown-flag/                 --not-a-real-flag → Error
│   └── extra-positional/             bare arg → Error
└── list/                             default mode
    ├── empty/
    │   └── no-sessions/              empty home → header + 0-count footer
    ├── table/                        multi-session human list
    │   ├── newest-first/             order by updated_at desc
    │   ├── columns-and-footer/       all columns + summary line
    │   ├── needs-attention-reason/   D1: live+!sendable+!exited → REASON
    │   ├── sendable-yes/             SENDABLE yes row
    │   └── unreachable-soft/         LIVE no; exit 0; optional warning
    ├── filter/
    │   ├── needs-confirm/            --needs-confirm only attention rows
    │   └── sendable-only/            --sendable only sendable rows
    ├── iterm/
    │   ├── single-match/             w=42 t=3
    │   ├── multi-plus-n/             w=1 t=2(+1)
    │   ├── missing-dash/             no TTY match → -
    │   ├── no-iterm-opt-out/         --no-iterm → all -
    │   └── soft-fail-warning/        ListITerm err → warning + -
    ├── json/
    │   ├── valid-no-ansi/            valid JSON structure, no ESC
    │   └── includes-iterm-fields/    iterm string/fields present
    ├── home/
    │   └── flag-isolation/           --home only sees that store
    └── limit/
        └── max-rows/                 --limit 1 of 3
```

Parameter ranking (most → least significant):

1. **Mode** — help | usage-errors | list (default)
2. **List data / concern** — empty | table classify | filter | iterm | json | home | limit
3. **Within table** — order, columns, D1, sendable, unreachable
4. **Within iterm** — match cardinality, opt-out, soft-fail

## Test Index

| # | Leaf | Description | Expect |
|---|------|-------------|--------|
| 1 | `help/documents-flags` | `-h` documents home, json, needs-confirm, sendable, no-iterm | GREEN |
| 2 | `usage-errors/unknown-flag` | unknown flag → error | GREEN |
| 3 | `usage-errors/extra-positional` | positional arg → error | GREEN |
| 4 | `list/empty/no-sessions` | empty home → 0 footer; exit 0 | GREEN |
| 5 | `list/table/newest-first` | two metas → newer updated_at first | GREEN |
| 6 | `list/table/columns-and-footer` | header columns + `N sessions · M · K` | GREEN |
| 7 | `list/table/needs-attention-reason` | D1 row shows REASON; counted | GREEN |
| 10 | `list/table/sendable-yes` | sendable probe → SENDABLE yes | RED |
| 11 | `list/table/unreachable-soft` | !Live → no crash; exit 0 | RED |
| 12 | `list/filter/needs-confirm` | only needs_attention rows | RED |
| 13 | `list/filter/sendable-only` | only sendable rows | RED |
| 14 | `list/iterm/single-match` | ITERM `w=42 t=3` | RED |
| 15 | `list/iterm/multi-plus-n` | ITERM `w=1 t=2(+1)` | RED |
| 16 | `list/iterm/missing-dash` | no match → `-` | RED |
| 17 | `list/iterm/no-iterm-opt-out` | `--no-iterm` → `-` | RED |
| 18 | `list/iterm/soft-fail-warning` | ListITerm err → `warning:` + `-` | RED |
| 19 | `list/json/valid-no-ansi` | JSON parse; no ANSI; summary counts | RED |
| 20 | `list/json/includes-iterm-fields` | JSON includes iterm field(s) | RED |
| 21 | `list/home/flag-isolation` | only sessions under `--home` | RED |
| 22 | `list/limit/max-rows` | `--limit 1` shows one of three | GREEN |

## How to Run

```sh
# from kck module root
doctest vet ./tests/list
doctest test ./tests/list
doctest test -v ./tests/list/list/table/newest-first
```

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"

	"kck/run"
)

// SessionSeed is one on-disk meta.json under Home/sessions/<id>/.
type SessionSeed struct {
	SessionID       string
	Runner          string
	RunnerSessionID string // meta.runner_session_id → AGENT_SID
	Status          string
	Workspace       string
	UpdatedAt       string // RFC3339
	CreatedAt       string
	// Probe overrides for this session id (used when building Probe map).
	Live     bool
	Sendable bool
	State    string
	Reason   string
	Exited   bool
	TTY      string
}

// Request drives one MainWith invocation.
type Request struct {
	Args []string

	// Home is the isolated agent-run home (root Setup sets t.TempDir path).
	Home string

	// Sessions written under Home before Run (meta.json only).
	Sessions []SessionSeed

	// ProbeByID maps session_id → ProbeResult. Empty → default dead/not-sendable.
	ProbeByID map[string]run.ProbeResult

	// ITerm sessions for ListITerm inject. ListITermErr soft-fails when set.
	ITermSessions []run.ITermSession
	ListITermErr  string

	// SkipProbe: leave Probe nil (production path — avoid in list leaves).
	SkipProbe bool
	// SkipITerm: leave ListITerm nil unless --no-iterm / soft-fail leaves set inject.
	SkipITerm bool

	// Fixed Now for relative UPDATED if product uses ages.
	Now time.Time
}

// Response is the observable MainWith outcome.
type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string // error.Error() if non-nil
	ExitCode int    // 0 if err==nil else 1 (thin-main convention)
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req.Home == "" {
		t.Fatal("req.Home empty; root Setup must set isolated home")
	}
	if err := seedSessions(req.Home, req.Sessions); err != nil {
		return nil, fmt.Errorf("seed sessions: %w", err)
	}

	// Ensure probe map includes seed probe fields when ProbeByID not fully filled.
	probeMap := mergeProbeMap(req)

	var stdout, stderr bytes.Buffer
	opts := run.Options{
		Args:   append([]string(nil), req.Args...),
		Home:   req.Home,
		Stdout: &stdout,
		Stderr: &stderr,
		Now:    req.Now,
	}
	if !req.SkipProbe {
		opts.Probe = func(sessionID string, meta run.SessionMetaView) (run.ProbeResult, error) {
			if p, ok := probeMap[sessionID]; ok {
				return p, nil
			}
			// Default: not live, not sendable, unknown state.
			return run.ProbeResult{State: "unknown"}, nil
		}
	}
	if !req.SkipITerm {
		opts.ListITerm = func() ([]run.ITermSession, error) {
			if req.ListITermErr != "" {
				return nil, fmt.Errorf("%s", req.ListITermErr)
			}
			return append([]run.ITermSession(nil), req.ITermSessions...), nil
		}
	}

	err := run.MainWith(opts)
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	// Harness never fails the leaf on product error — Assert decides success/fail.
	return resp, nil
}

func seedSessions(home string, seeds []SessionSeed) error {
	for _, s := range seeds {
		id := strings.TrimSpace(s.SessionID)
		if id == "" {
			return fmt.Errorf("session seed missing SessionID")
		}
		dir := filepath.Join(home, "sessions", id)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		runner := s.Runner
		if runner == "" {
			runner = "grok"
		}
		status := s.Status
		if status == "" {
			status = "running"
		}
		updated := s.UpdatedAt
		if updated == "" {
			updated = "2026-08-01T12:00:00Z"
		}
		created := s.CreatedAt
		if created == "" {
			created = updated
		}
		meta := map[string]any{
			"runner":     runner,
			"session_id": id,
			"status":     status,
			"workspace":  s.Workspace,
			"updated_at": updated,
			"created_at": created,
		}
		if sid := strings.TrimSpace(s.RunnerSessionID); sid != "" {
			meta["runner_session_id"] = sid
		}
		data, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "meta.json"), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func mergeProbeMap(req *Request) map[string]run.ProbeResult {
	out := make(map[string]run.ProbeResult)
	for k, v := range req.ProbeByID {
		out[k] = v
	}
	for _, s := range req.Sessions {
		if _, ok := out[s.SessionID]; ok {
			continue
		}
		// Only auto-fill when leaf set any probe-related field on the seed.
		if s.Live || s.Sendable || s.Exited || s.State != "" || s.Reason != "" || s.TTY != "" {
			out[s.SessionID] = run.ProbeResult{
				Live:     s.Live,
				Sendable: s.Sendable,
				State:    s.State,
				Reason:   s.Reason,
				Exited:   s.Exited,
				TTY:      s.TTY,
			}
		}
	}
	return out
}
```
