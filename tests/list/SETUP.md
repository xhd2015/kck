# Scenario

**Feature**: kck lists agent-run sessions with injectable live probe and iTerm

```
# isolated home + seeded meta.json
Caller -> run.MainWith(Options{Home, Args, Stdout, Stderr, Probe, ListITerm})
  -> session store under Home
  -> Probe(session) serial
  -> ListITerm (unless --no-iterm)
  -> human table | JSON | help
```

## Preconditions

- Package `kck/run` exposes `MainWith(Options)`, `ProbeResult`, `SessionMetaView`,
  `ITermSession`, and related types locked in root `DOCTEST.md`.
- Classic TDD: leaves **RED** until list behavior and injectables land.
- Parallel-safe: every leaf gets its own `t.TempDir()` home; no `t.Setenv`,
  `os.Setenv`, or `os.Chdir` for home or cwd.
- Fixtures write `sessions/<id>/meta.json` only (no agentstorage import required
  in harness).

## Steps

1. Root `Setup` sets `req.Home` to an empty isolated directory.
2. Leaf `Setup` sets `Args`, `Sessions`, probe fields, and/or iTerm fixtures.
3. Root `Run` seeds metas, builds injectables, calls `run.MainWith`, captures
   writers and error text.

## Context

- D1 `needs_attention` = Live && !Sendable && !Exited.
- Footer wording locked: `N sessions · M needs attention · K sendable`.
- Buffer writers are non-TTY → human list without requiring ANSI (product may
  still color if forced; asserts allow plain text).

```go
import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Home == "" {
		req.Home = filepath.Join(t.TempDir(), ".agent-run")
	}
	if req.Now.IsZero() {
		req.Now = time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)
	}
	if req.ProbeByID == nil {
		req.ProbeByID = map[string]run.ProbeResult{}
	}
	return nil
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
}

func assertSuccess(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode != 0 || resp.ErrText != "" {
		t.Fatalf("want success exit 0, got exit=%d err=%q stderr=%q stdout=%q",
			resp.ExitCode, resp.ErrText, resp.Stderr, resp.Stdout)
	}
}

func assertFailure(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode == 0 && resp.ErrText == "" {
		t.Fatalf("want non-zero failure, got exit=0 empty err; stdout=%q stderr=%q",
			resp.Stdout, resp.Stderr)
	}
}

func assertContains(t *testing.T, got, want, label string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s must contain %q; got:\n%s", label, want, got)
	}
}

func assertNotContains(t *testing.T, got, ban, label string) {
	t.Helper()
	if strings.Contains(got, ban) {
		t.Fatalf("%s must not contain %q; got:\n%s", label, ban, got)
	}
}

func assertNoANSI(t *testing.T, s, label string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("%s must not contain ANSI escapes; got:\n%q", label, s)
	}
}

func assertTrailingNewline(t *testing.T, s, label string) {
	t.Helper()
	if s == "" || !strings.HasSuffix(s, "\n") {
		t.Fatalf("%s must end with trailing newline; got %q", label, s)
	}
}

// twoSessionFixture: older s-old + newer s-new for order / multi-row leaves.
func twoSessionFixture(req *Request) {
	req.Sessions = []SessionSeed{
		{
			SessionID: "s-old",
			Runner:    "grok",
			Status:    "running",
			Workspace: "/ws/old",
			UpdatedAt: "2026-08-01T10:00:00Z",
			CreatedAt: "2026-08-01T09:00:00Z",
			Live:      true,
			Sendable:  true,
			State:     "idle",
			TTY:       "/dev/ttys010",
		},
		{
			SessionID: "s-new",
			Runner:    "codex",
			Status:    "running",
			Workspace: "/ws/new",
			UpdatedAt: "2026-08-02T10:00:00Z",
			CreatedAt: "2026-08-02T09:00:00Z",
			Live:      true,
			Sendable:  false,
			State:     "running",
			Reason:    "awaiting confirmation",
			Exited:    false,
			TTY:       "/dev/ttys011",
		},
	}
}

// threeSessionFixture for limit / multi-filter cases.
// IDs are distinctive substrings so asserts do not false-match footer text.
func threeSessionFixture(req *Request) {
	req.Sessions = []SessionSeed{
		{
			SessionID: "sess-sendable",
			Runner:    "grok",
			Workspace: "/wa",
			UpdatedAt: "2026-08-03T10:00:00Z",
			Live:      true,
			Sendable:  true,
			State:     "idle",
			TTY:       "/dev/ttys001",
		},
		{
			SessionID: "sess-attention",
			Runner:    "grok",
			Workspace: "/wb",
			UpdatedAt: "2026-08-02T10:00:00Z",
			Live:      true,
			Sendable:  false,
			State:     "running",
			Reason:    "tool running",
			TTY:       "/dev/ttys002",
		},
		{
			SessionID: "sess-exited",
			Runner:    "codex",
			Workspace: "/wc",
			UpdatedAt: "2026-08-01T10:00:00Z",
			Live:      false,
			Sendable:  false,
			State:     "exited",
			Exited:    true,
		},
	}
}
```
