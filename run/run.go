package run

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/agent-pro/pkgs/itermsnapshot"
	lessflags "github.com/xhd2015/less-flags"
)

// ProbeResult is one session's live classification for the list row.
type ProbeResult struct {
	Live     bool   // terminal reachable / process alive
	Sendable bool   // idle writable (can accept send)
	State    string // e.g. idle, running, exited, unknown
	Reason   string // needs_attention reason when applicable; else empty
	Exited   bool   // clearly exited → not needs_attention
	TTY      string // device path
}

// SessionMetaView is the meta subset passed to Probe.
type SessionMetaView struct {
	SessionID      string
	Runner         string
	Status         string
	Workspace      string
	UpdatedAt      string
	CreatedAt      string
	AgentSessionID string // runner-native id (store: meta.runner_session_id)
}

// ProbeFunc classifies a session for LIVE/SENDABLE/STATE/REASON/ITERM matching.
type ProbeFunc func(sessionID string, meta SessionMetaView) (ProbeResult, error)

// ITermSession is one iTerm2 window/tab reference with its TTY.
type ITermSession struct {
	WindowID string
	TabIndex int // 1-based
	TTY      string
}

// ListITermFunc lists iTerm sessions for TTY → window/tab matching.
type ListITermFunc func() ([]ITermSession, error)

// Options configures an injectable MainWith invocation (L2 harness / production Main).
type Options struct {
	Args      []string // argv after program name
	Home      string   // explicit home; --home in Args wins when both set
	Stdout    io.Writer
	Stderr    io.Writer
	Probe     ProbeFunc
	ListITerm ListITermFunc
	Now       time.Time // zero → time.Now(); for UPDATED ages if relative

	// LiveCapture injects itermsnapshot result for L2 live path.
	// nil → production: itermsnapshot.Capture(CaptureOpts{}).
	// Invoked only when resolved home is empty (live mode).
	// Signature: (result, warnings, error). Hard error is soft-failed by kck.
	LiveCapture func() (*itermsnapshot.Result, []string, error)

	// GrokHome for `kck grok …`. Empty → agenttty.GrokHome().
	GrokHome string
	// GrokOpenOpts injects open probes/launchers for L2. nil → production.
	GrokOpenOpts *sessions.OpenOpts
	// GrokSnapshotOpts injects snapshot probes/Contents for L2. nil → production.
	GrokSnapshotOpts *sessions.SnapshotOpts
	// GrokSendOpts injects send probes/SendText/open for L2. nil → production.
	GrokSendOpts *sessions.SendOpts
	// GrokMessagesOpts injects messages tab-resolve probes for L2. nil → production.
	GrokMessagesOpts *sessions.MessagesOpts
	// GrokListLiveOpts injects list-live probes for L2. nil → production.
	GrokListLiveOpts *sessions.ListLiveOpts
	// GrokLiveOpts injects Status/Info live PID probes for L2. nil → production.
	GrokLiveOpts *sessions.LiveOptions
	// GrokResolveOpts injects resolve ancestor/tab probes for L2. nil → production.
	GrokResolveOpts *sessions.ResolveOpts
	// GrokNow for info relative timestamps. Zero → time.Now().
	GrokNow time.Time

	// GrokCronNow injects the clock for --cron loops. nil → time.Now.
	GrokCronNow func() time.Time
	// GrokCronSleep injects wait between cron ticks. nil → time.Sleep.
	// Returning a non-nil error aborts the loop (e.g. tests / cancellation).
	GrokCronSleep func(time.Duration) error
	// GrokCronLoc is the location for easy-cron TOD math. nil → time.Local.
	GrokCronLoc *time.Location
}

const helpText = `Usage: kck [OPTIONS]
       kck grok list [OPTIONS]
       kck grok open <session-id> [OPTIONS]
       kck grok snapshot <session-id> [OPTIONS]
       kck grok send <text> --session-id <id> [OPTIONS]
       kck grok messages <session-id> [OPTIONS]
       kck grok info <session-id> [OPTIONS]
       kck grok status <session-id> [OPTIONS]
       kck grok resolve [OPTIONS]
       kck skill --show|--list|--install …

Default mode: list live iTerm agent panes (streams rows as windows are scanned).
With --home: list agent-run store sessions under that home.

Commands:
  grok list …              list Grok ids hosted in iTerm tabs
  grok open …              focus hosting tab or resume (--tab / --tab-index / <id>)
  grok snapshot …          capture visible pane text (--tab / --tab-index / <id>)
  grok send …              type text into hosting pane (--session-id / --tab / --open)
  grok messages …          print recent chat messages (--limit / --offset-from-end)
  grok info …              show session detail + Active block
  grok status …            dual-signal liveness + session path
  grok resolve …           resolve Grok session id (ancestor walk or --tab)
  skill                    show/install embedded skill docs

Options:
  -h, --help            show help message
  --home PATH           agent-run home (sessions under PATH/sessions)
  --json                machine-readable JSON list (buffered; no ANSI)
  --needs-confirm       only sessions that need attention (live, not sendable)
  --sendable            only sendable (idle writable) sessions
  --no-iterm            skip iTerm window/tab resolution (ITERM shows -)
  --fast                skip lsof; AGENT_SID stays - (faster; live only)
  --enrich              also resolve grok title/model after sid hit (live only)
  --limit N             show at most N newest sessions

Columns AGENT_RUN (process under agent-run) and AGENT_SID (runner session id)
are always present. Live: AGENT_RUN from process tree; AGENT_SID via lsof by
default (use --fast to skip). Store: AGENT_RUN=yes, AGENT_SID from
meta.runner_session_id.

Run kck skill --help and kck skill --install --help for skill docs.
`

// Main is production CLI: MainWith with os.Stdout/os.Stderr.
func Main(args []string) error {
	return MainWith(Options{
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

// MainWith runs kck with injectable IO / home / probe / iTerm (L2).
// Success (list, help) → nil. Failures (bad flags) → non-nil error after writing
// user-facing "Error: …" lines to Stderr when appropriate.
func MainWith(opts Options) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	if len(opts.Args) > 0 {
		switch opts.Args[0] {
		case "grok":
			return runGrok(opts)
		case "skill":
			return runSkill(opts)
		}
	}

	var (
		home         string
		jsonOut      bool
		needsConfirm bool
		sendableOnly bool
		noITerm      bool
		fast         bool
		enrich       bool
		limit        int
	)

	remain, err := lessflags.
		String("--home", &home).
		Bool("--json", &jsonOut).
		Bool("--needs-confirm", &needsConfirm).
		Bool("--sendable", &sendableOnly).
		Bool("--no-iterm", &noITerm).
		Bool("--fast", &fast).
		Bool("--enrich", &enrich).
		Int("--limit", &limit).
		HelpFunc("-h,--help", func() {
			txt := strings.TrimPrefix(helpText, "\n")
			if !strings.HasSuffix(txt, "\n") {
				txt += "\n"
			}
			fmt.Fprint(stdout, txt)
		}).
		HelpNoExit().
		Parse(opts.Args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return writeError(stderr, err.Error())
	}
	if len(remain) > 0 {
		return writeError(stderr, fmt.Sprintf("unrecognized extra args: %s", strings.Join(remain, " ")))
	}

	// Home: --home flag wins over Options.Home.
	resolvedHome := opts.Home
	if home != "" {
		resolvedHome = home
	}

	// Mode route: non-empty home → store list; empty → live iTerm inventory.
	if resolvedHome == "" {
		// LiveCapture inject (tests) → buffered. nil → production streams
		// human rows via CaptureProgressive (unless --json).
		return runLiveList(stdout, stderr, liveConfig{
			JSON:         jsonOut,
			NeedsConfirm: needsConfirm,
			SendableOnly: sendableOnly,
			NoITerm:      noITerm,
			Limit:        limit,
			Fast:         fast,
			Enrich:       enrich,
		}, opts.LiveCapture)
	}

	return runList(stdout, stderr, listConfig{
		Home:         resolvedHome,
		JSON:         jsonOut,
		NeedsConfirm: needsConfirm,
		SendableOnly: sendableOnly,
		NoITerm:      noITerm,
		Limit:        limit,
		Probe:        opts.Probe,
		ListITerm:    opts.ListITerm,
		Now:          opts.Now,
	})
}

func writeError(stderr io.Writer, msg string) error {
	// Avoid double "Error: " if caller already included it.
	if strings.HasPrefix(msg, "Error:") {
		fmt.Fprintln(stderr, msg)
		return errors.New(msg)
	}
	line := "Error: " + msg
	fmt.Fprintln(stderr, line)
	return errors.New(line)
}


type listConfig struct {
	Home         string
	JSON         bool
	NeedsConfirm bool
	SendableOnly bool
	NoITerm      bool
	Limit        int
	Probe        ProbeFunc
	ListITerm    ListITermFunc
	Now          time.Time
}

type listRow struct {
	Meta           SessionMetaView
	Probe          ProbeResult
	ITerm          string
	NeedsAtt       bool
	AgentRun       string // "yes" | "no" | "" (unknown → "-")
	AgentSessionID string // runner-native id (grok/codex); empty → "-"
}

func runList(stdout, stderr io.Writer, cfg listConfig) error {
	metas, err := loadSessionMetas(cfg.Home)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	sortMetasNewestFirst(metas)

	probe := cfg.Probe
	if probe == nil {
		probe = defaultProbe
	}

	// Resolve iTerm list once (unless opted out).
	var itermByTTY map[string][]ITermSession
	itermOK := true
	if !cfg.NoITerm {
		listFn := cfg.ListITerm
		if listFn != nil {
			sessions, err := listFn()
			if err != nil {
				itermOK = false
				fmt.Fprintf(stderr, "warning: iTerm list failed: %v\n", err)
			} else {
				itermByTTY = indexITermByTTY(sessions)
			}
		}
	}

	rows := make([]listRow, 0, len(metas))
	for _, m := range metas {
		pr, err := probe(m.SessionID, m)
		if err != nil {
			// Soft-fail probe errors: mark not live, continue.
			fmt.Fprintf(stderr, "warning: probe %s: %v\n", m.SessionID, err)
			pr = ProbeResult{State: "unknown"}
		}
		iterm := "-"
		if !cfg.NoITerm && itermOK {
			iterm = formatITerm(pr.TTY, itermByTTY)
		}
		needs := pr.Live && !pr.Sendable && !pr.Exited
		// Store rows are agent-run owned by definition.
		rows = append(rows, listRow{
			Meta:           m,
			Probe:          pr,
			ITerm:          iterm,
			NeedsAtt:       needs,
			AgentRun:       "yes",
			AgentSessionID: m.AgentSessionID,
		})
	}

	// Filters apply before limit; counts match displayed set.
	filtered := make([]listRow, 0, len(rows))
	for _, r := range rows {
		if cfg.NeedsConfirm && !r.NeedsAtt {
			continue
		}
		if cfg.SendableOnly && !r.Probe.Sendable {
			continue
		}
		filtered = append(filtered, r)
	}
	if cfg.Limit > 0 && len(filtered) > cfg.Limit {
		filtered = filtered[:cfg.Limit]
	}

	if cfg.JSON {
		return writeJSON(stdout, filtered)
	}
	return writeHumanTable(stdout, filtered)
}

func defaultProbe(sessionID string, meta SessionMetaView) (ProbeResult, error) {
	_ = sessionID
	_ = meta
	return ProbeResult{State: "unknown"}, nil
}

type diskMeta struct {
	Runner          string `json:"runner"`
	SessionID       string `json:"session_id"`
	RunnerSessionID string `json:"runner_session_id"`
	Status          string `json:"status"`
	Workspace       string `json:"workspace"`
	UpdatedAt       string `json:"updated_at"`
	CreatedAt       string `json:"created_at"`
}

func loadSessionMetas(home string) ([]SessionMetaView, error) {
	if home == "" {
		return nil, nil
	}
	dir := filepath.Join(home, "sessions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]SessionMetaView, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		id := e.Name()
		data, err := os.ReadFile(filepath.Join(dir, id, "meta.json"))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		var m diskMeta
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("session %s meta.json: %w", id, err)
		}
		if m.SessionID == "" {
			m.SessionID = id
		}
		out = append(out, SessionMetaView{
			SessionID:      m.SessionID,
			Runner:         m.Runner,
			Status:         m.Status,
			Workspace:      m.Workspace,
			UpdatedAt:      m.UpdatedAt,
			CreatedAt:      m.CreatedAt,
			AgentSessionID: strings.TrimSpace(m.RunnerSessionID),
		})
	}
	return out, nil
}

func sortMetasNewestFirst(list []SessionMetaView) {
	sort.SliceStable(list, func(i, j int) bool {
		ui := parseTime(list[i].UpdatedAt)
		uj := parseTime(list[j].UpdatedAt)
		if !ui.Equal(uj) {
			return ui.After(uj)
		}
		ci := parseTime(list[i].CreatedAt)
		cj := parseTime(list[j].CreatedAt)
		if !ci.Equal(cj) {
			return ci.After(cj)
		}
		return list[i].SessionID < list[j].SessionID
	})
}

func parseTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Time{}
}

func indexITermByTTY(sessions []ITermSession) map[string][]ITermSession {
	m := make(map[string][]ITermSession)
	for _, s := range sessions {
		tty := strings.TrimSpace(s.TTY)
		if tty == "" {
			continue
		}
		m[tty] = append(m[tty], s)
	}
	return m
}

func formatITerm(tty string, byTTY map[string][]ITermSession) string {
	tty = strings.TrimSpace(tty)
	if tty == "" || byTTY == nil {
		return "-"
	}
	matches := byTTY[tty]
	if len(matches) == 0 {
		return "-"
	}
	first := matches[0]
	base := fmt.Sprintf("w=%s t=%d", first.WindowID, first.TabIndex)
	if len(matches) > 1 {
		base = fmt.Sprintf("%s(+%d)", base, len(matches)-1)
	}
	return base
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// Fixed-width human table so streaming rows stay column-aligned from the first
// line (tabwriter needs a final Flush and would reflow). ITERM may contain a
// space ("w=42 t=3") but is one cell via width, not field-splitting.
const (
	colSession   = 36
	colRunner    = 10
	colLive      = 4
	colSendable  = 8
	colState     = 7
	colReason    = 16
	colAgentRun  = 9
	colAgentSID  = 36
	colITerm     = 16
	colUpdated   = 10
)

func writeHumanHeader(stdout io.Writer) error {
	_, err := fmt.Fprintf(stdout, "%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %s\n",
		colSession, "SESSION_ID",
		colRunner, "RUNNER",
		colLive, "LIVE",
		colSendable, "SENDABLE",
		colState, "STATE",
		colReason, "REASON",
		colAgentRun, "AGENT_RUN",
		colAgentSID, "AGENT_SID",
		colITerm, "ITERM",
		colUpdated, "UPDATED",
		"WORKSPACE",
	)
	return err
}

func writeHumanRow(stdout io.Writer, r listRow) error {
	// Each field is one %-width cell (spaces inside ITERM like "w=42 t=3" stay in-cell).
	// Values longer than the pad width still print fully (may nudge later cols on that row only).
	_, err := fmt.Fprintf(stdout, "%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %s\n",
		colSession, r.Meta.SessionID,
		colRunner, emptyDash(r.Meta.Runner),
		colLive, yesNo(r.Probe.Live),
		colSendable, yesNo(r.Probe.Sendable),
		colState, emptyDash(r.Probe.State),
		colReason, emptyField(r.Probe.Reason),
		colAgentRun, emptyDash(r.AgentRun),
		colAgentSID, emptyDash(r.AgentSessionID),
		colITerm, emptyDash(r.ITerm),
		colUpdated, emptyDash(r.Meta.UpdatedAt),
		emptyDash(r.Meta.Workspace),
	)
	return err
}

func writeHumanTable(stdout io.Writer, rows []listRow) error {
	if err := writeHumanHeader(stdout); err != nil {
		return err
	}
	n, needs, sendable := 0, 0, 0
	for _, r := range rows {
		n++
		if r.NeedsAtt {
			needs++
		}
		if r.Probe.Sendable {
			sendable++
		}
		if err := writeHumanRow(stdout, r); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(stdout, "%d sessions · %d needs attention · %d sendable\n", n, needs, sendable)
	return err
}

func emptyDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func emptyField(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

// agentRunJSON maps "yes"/"no"/"" → true/false/null for JSON.
func agentRunJSON(s string) *bool {
	switch strings.TrimSpace(s) {
	case "yes":
		t := true
		return &t
	case "no":
		f := false
		return &f
	default:
		return nil
	}
}

type jsonSession struct {
	SessionID      string `json:"session_id"`
	Runner         string `json:"runner"`
	Live           bool   `json:"live"`
	Sendable       bool   `json:"sendable"`
	State          string `json:"state"`
	Reason         string `json:"reason,omitempty"`
	AgentRun       *bool  `json:"agent_run"` // null when unknown
	AgentSessionID string `json:"agent_session_id"`
	ITerm          string `json:"iterm"`
	UpdatedAt      string `json:"updated_at"`
	Workspace      string `json:"workspace"`
	NeedsAtt       bool   `json:"needs_attention"`
}

type jsonList struct {
	Sessions []jsonSession `json:"sessions"`
	Summary  jsonSummary   `json:"summary"`
}

type jsonSummary struct {
	Total          int `json:"total"`
	NeedsAttention int `json:"needs_attention"`
	Sendable       int `json:"sendable"`
}

func writeJSON(stdout io.Writer, rows []listRow) error {
	out := jsonList{
		Sessions: make([]jsonSession, 0, len(rows)),
	}
	for _, r := range rows {
		if r.NeedsAtt {
			out.Summary.NeedsAttention++
		}
		if r.Probe.Sendable {
			out.Summary.Sendable++
		}
		out.Sessions = append(out.Sessions, jsonSession{
			SessionID:      r.Meta.SessionID,
			Runner:         r.Meta.Runner,
			Live:           r.Probe.Live,
			Sendable:       r.Probe.Sendable,
			State:          r.Probe.State,
			Reason:         r.Probe.Reason,
			AgentRun:       agentRunJSON(r.AgentRun),
			AgentSessionID: r.AgentSessionID,
			ITerm:          r.ITerm,
			UpdatedAt:      r.Meta.UpdatedAt,
			Workspace:      r.Meta.Workspace,
			NeedsAtt:       r.NeedsAtt,
		})
	}
	out.Summary.Total = len(out.Sessions)
	enc := json.NewEncoder(stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(out); err != nil {
		return err
	}
	// Encoder already adds trailing newline.
	return nil
}
