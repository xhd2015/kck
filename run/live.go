package run

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/xhd2015/agent-pro/pkgs/itermsnapshot"
	"github.com/xhd2015/agent-pro/pkgs/procresolve"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2/snapshot"
)

// agentLikeTokens are command basenames / whitespace tokens that mark a pane
// as agent-like for the default live agents-only filter.
var agentLikeTokens = []string{"grok", "codex", "mark", "agent-run"}

type liveConfig struct {
	JSON         bool
	NeedsConfirm bool
	SendableOnly bool
	NoITerm      bool
	Limit        int
	Fast         bool // skip lsof; AGENT_SID stays "-"
	Enrich       bool // EnrichInfo title/model only (after hard sid hit)
}

func runLiveList(stdout, stderr io.Writer, cfg liveConfig, capture func() (*itermsnapshot.Result, []string, error)) error {
	// Injected capture (L2 tests) or JSON (needs full document): buffered path.
	if capture != nil || cfg.JSON {
		return runLiveListBuffered(stdout, stderr, cfg, capture)
	}
	// Production human list: stream rows as each iTerm window is ready.
	return runLiveListStreaming(stdout, stderr, cfg)
}

func runLiveListBuffered(stdout, stderr io.Writer, cfg liveConfig, capture func() (*itermsnapshot.Result, []string, error)) error {
	if capture == nil {
		// Production buffered (JSON): base inventory, then soft attach with
		// process tree + default lsof for AGENT_SID (--fast skips lsof).
		capture = func() (*itermsnapshot.Result, []string, error) {
			res, warn, err := itermsnapshot.Capture(itermsnapshot.CaptureOpts{NoEnrich: true})
			if err != nil {
				return res, warn, err
			}
			softAttachAgents(res, cfg)
			return res, warn, nil
		}
	}

	result, warnings, err := capture()
	if err != nil {
		fmt.Fprintf(stderr, "warning: live capture failed: %v\n", err)
		if cfg.JSON {
			return writeJSON(stdout, nil)
		}
		return writeHumanTable(stdout, nil)
	}
	printLiveWarnings(stderr, warnings)

	rows := liveRowsFromResult(result, cfg.NoITerm)
	filtered := filterLiveRows(rows, cfg)
	if cfg.JSON {
		return writeJSON(stdout, filtered)
	}
	return writeHumanTable(stdout, filtered)
}

// makeLiveResolve builds a shared-ListProcs ResolveFromPID.
// Tree/agent-run roles always. Lsof (AGENT_SID) is on by default; --fast skips
// it. EnrichInfo (title/model) only when cfg.Enrich.
func makeLiveResolve(cfg liveConfig) func(int) (*procresolve.Result, error) {
	procs := procresolve.ListLiveProcs()
	var lsof func(int) []string
	if !cfg.Fast {
		lsof = procresolve.LiveLsof
	}
	return func(pid int) (*procresolve.Result, error) {
		return procresolve.ResolveFromPID(pid, procresolve.Options{
			ListProcs:  func() []procresolve.Proc { return procs },
			Lsof:       lsof,
			EnrichInfo: cfg.Enrich && !cfg.Fast,
		})
	}
}

// softAttachAgents fills Result.Agents for busy panes using soft tree attach
// (and default lsof for AGENT_SID unless --fast).
func softAttachAgents(result *itermsnapshot.Result, cfg liveConfig) {
	if result == nil || result.Snapshot == nil {
		return
	}
	resolve := makeLiveResolve(cfg)
	agents := result.Agents
	if agents == nil {
		agents = map[string]*itermsnapshot.SessionAgent{}
	}
	for wi := range result.Snapshot.Windows {
		win := &result.Snapshot.Windows[wi]
		for ti := range win.Tabs {
			tab := &win.Tabs[ti]
			for si := range tab.Sessions {
				sess := &tab.Sessions[si]
				if _, ok := agents[sess.ID]; ok {
					continue
				}
				if ag := attachAgentFromResolve(sess, resolve); ag != nil {
					agents[sess.ID] = ag
				}
			}
		}
	}
	if len(agents) > 0 {
		result.Agents = agents
	}
}

// runLiveListStreaming prints the header immediately, then one row at a time as
// each iTerm window finishes process enrich (CaptureProgressive / phased AS).
func runLiveListStreaming(stdout, stderr io.Writer, cfg liveConfig) error {
	// Header first so the user sees activity immediately.
	if err := writeHumanHeader(stdout); err != nil {
		return err
	}
	flushWriter(stdout)

	// Process tree always; lsof for AGENT_SID unless --fast.
	resolve := makeLiveResolve(cfg)
	var (
		n, needs, sendable int
		hitLimit           bool
	)

	c := snapshot.NewCollector()
	_, warnings, err := c.CaptureProgressive(func(win snapshot.SnapshotWindow) error {
		if hitLimit {
			return nil
		}
		agents := map[string]*itermsnapshot.SessionAgent{}
		for ti := range win.Tabs {
			for si := range win.Tabs[ti].Sessions {
				sess := &win.Tabs[ti].Sessions[si]
				if ag := attachAgentFromResolve(sess, resolve); ag != nil {
					agents[sess.ID] = ag
				}
			}
		}
		for ti := range win.Tabs {
			tab := &win.Tabs[ti]
			for si := range tab.Sessions {
				if hitLimit {
					return nil
				}
				sess := &tab.Sessions[si]
				if !includeLiveSession(sess, agents) {
					continue
				}
				row := mapLiveSession(sess, &win, tab, agents[sess.ID], cfg.NoITerm)
				if cfg.NeedsConfirm && !row.NeedsAtt {
					continue
				}
				if cfg.SendableOnly && !row.Probe.Sendable {
					continue
				}
				if err := writeHumanRow(stdout, row); err != nil {
					return err
				}
				flushWriter(stdout)
				n++
				if row.NeedsAtt {
					needs++
				}
				if row.Probe.Sendable {
					sendable++
				}
				if cfg.Limit > 0 && n >= cfg.Limit {
					hitLimit = true
					return nil
				}
			}
		}
		return nil
	})
	printLiveWarnings(stderr, warnings)
	if err != nil {
		fmt.Fprintf(stderr, "warning: live capture failed: %v\n", err)
		// Still emit footer so the table is well-formed.
		fmt.Fprintf(stdout, "%d sessions · %d needs attention · %d sendable\n", n, needs, sendable)
		flushWriter(stdout)
		return nil
	}
	fmt.Fprintf(stdout, "%d sessions · %d needs attention · %d sendable\n", n, needs, sendable)
	flushWriter(stdout)
	return nil
}

func printLiveWarnings(stderr io.Writer, warnings []string) {
	for _, w := range warnings {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(w), "warning:") {
			fmt.Fprintln(stderr, w)
		} else {
			fmt.Fprintf(stderr, "warning: %s\n", w)
		}
	}
}

func filterLiveRows(rows []listRow, cfg liveConfig) []listRow {
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
	return filtered
}

func attachAgentFromResolve(sess *snapshot.SnapshotSession, resolve func(int) (*procresolve.Result, error)) *itermsnapshot.SessionAgent {
	if sess.Idle == nil || *sess.Idle {
		return nil
	}
	// Prefer ShellPID so ResolveFromPID walks the full pane tree (descendants
	// only). Preferring sess.PID when it is the agent binary itself labels the
	// root as "input" and hides parent agent-run + child runners.
	var pid int
	if sess.ShellPID != nil && *sess.ShellPID > 0 {
		pid = *sess.ShellPID
	} else if sess.PID != nil {
		pid = *sess.PID
	}
	if pid <= 0 || resolve == nil {
		return nil
	}
	r, err := resolve(pid)
	if err != nil || r == nil {
		return nil
	}
	tree := make([]itermsnapshot.AgentTreeNode, len(r.Tree))
	for i, n := range r.Tree {
		tree[i] = itermsnapshot.AgentTreeNode{PID: n.PID, PPID: n.PPID, Role: n.Role, Cmd: n.Cmd}
	}
	kind := strings.TrimSpace(r.Kind)
	if kind == "none" {
		kind = ""
	}
	// Soft attach: keep tree even without a hard runner session id so AGENT_RUN
	// can be classified from process roles. Infer kind from Role or Cmd.
	if kind == "" {
		kind = inferRunnerKind(tree)
	}
	sid := strings.TrimSpace(r.SessionID)
	// Keep soft attach when we have kind, sid, agent-run, or a runner in tree.
	if kind == "" && sid == "" && !treeHasAgentRun(tree) && !treeHasRunner(tree) {
		return nil
	}
	return &itermsnapshot.SessionAgent{
		Kind:      kind,
		SessionID: sid,
		Title:     r.GrokTitle,
		Tree:      tree,
	}
}

func cmdBasename(cmd string) string {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return ""
	}
	return strings.ToLower(filepath.Base(fields[0]))
}

func inferRunnerKind(tree []itermsnapshot.AgentTreeNode) string {
	for _, n := range tree {
		if n.Role == "grok" || n.Role == "codex" {
			return n.Role
		}
	}
	// Root role is forced to "input"; still detect runners via Cmd basename.
	for _, n := range tree {
		base := cmdBasename(n.Cmd)
		if base == "grok" || base == "codex" {
			return base
		}
	}
	return ""
}

func treeHasAgentRun(tree []itermsnapshot.AgentTreeNode) bool {
	for _, n := range tree {
		if n.Role == "agent-run" {
			return true
		}
		// Basename fallback (role agent-run-serve is not managed TTY).
		if n.Role != "agent-run-serve" && cmdBasename(n.Cmd) == "agent-run" {
			return true
		}
	}
	return false
}

func treeHasRunner(tree []itermsnapshot.AgentTreeNode) bool {
	for _, n := range tree {
		if n.Role == "grok" || n.Role == "codex" {
			return true
		}
		base := cmdBasename(n.Cmd)
		if base == "grok" || base == "codex" {
			return true
		}
	}
	return false
}

// agentRunFromAgent classifies AGENT_RUN for a live pane.
// nil agent → unknown (""); tree role agent-run → yes; attached otherwise → no.
func agentRunFromAgent(ag *itermsnapshot.SessionAgent) string {
	if ag == nil {
		return ""
	}
	if treeHasAgentRun(ag.Tree) {
		return "yes"
	}
	// Attached (kind/sid/tree) without agent-run role ⇒ not managed by agent-run.
	if strings.TrimSpace(ag.Kind) != "" || strings.TrimSpace(ag.SessionID) != "" || len(ag.Tree) > 0 {
		return "no"
	}
	return ""
}

func flushWriter(w io.Writer) {
	type flusher interface{ Flush() error }
	if f, ok := w.(flusher); ok {
		_ = f.Flush()
	}
}

func liveRowsFromResult(result *itermsnapshot.Result, noITerm bool) []listRow {
	if result == nil || result.Snapshot == nil {
		return nil
	}
	agents := result.Agents
	if agents == nil {
		agents = map[string]*itermsnapshot.SessionAgent{}
	}

	var rows []listRow
	for wi := range result.Snapshot.Windows {
		win := &result.Snapshot.Windows[wi]
		for ti := range win.Tabs {
			tab := &win.Tabs[ti]
			for si := range tab.Sessions {
				sess := &tab.Sessions[si]
				if !includeLiveSession(sess, agents) {
					continue
				}
				rows = append(rows, mapLiveSession(sess, win, tab, agents[sess.ID], noITerm))
			}
		}
	}
	return rows
}

func includeLiveSession(sess *snapshot.SnapshotSession, agents map[string]*itermsnapshot.SessionAgent) bool {
	if ag := agents[sess.ID]; ag != nil {
		if strings.TrimSpace(ag.Kind) != "" || strings.TrimSpace(ag.SessionID) != "" {
			return true
		}
	}
	return matchAgentLikeToken(sess) != ""
}

func matchAgentLikeToken(sess *snapshot.SnapshotSession) string {
	candidates := []string{sess.Name}
	if sess.Command != nil {
		candidates = append(candidates, *sess.Command)
	}
	if sess.CommandLine != nil {
		candidates = append(candidates, *sess.CommandLine)
	}
	for _, p := range sess.Processes {
		candidates = append(candidates, p.Command)
	}
	for _, c := range candidates {
		if tok := findAgentLikeInText(c); tok != "" {
			return tok
		}
	}
	return ""
}

// findAgentLikeInText returns the first agent-like token found as a path
// basename or whitespace-delimited field (case-insensitive).
func findAgentLikeInText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// Whole-string basename (e.g. /usr/local/bin/mark).
	base := strings.ToLower(filepath.Base(s))
	for _, tok := range agentLikeTokens {
		if base == tok {
			return tok
		}
	}
	// Whitespace-delimited tokens; each field may be a path.
	for _, field := range strings.Fields(s) {
		fbase := strings.ToLower(filepath.Base(field))
		for _, tok := range agentLikeTokens {
			if fbase == tok || strings.EqualFold(field, tok) {
				return tok
			}
		}
	}
	return ""
}

func mapLiveSession(
	sess *snapshot.SnapshotSession,
	win *snapshot.SnapshotWindow,
	tab *snapshot.SnapshotTab,
	ag *itermsnapshot.SessionAgent,
	noITerm bool,
) listRow {
	sessionID := sess.ID
	runner := ""
	agentSID := ""
	if ag != nil {
		if sid := strings.TrimSpace(ag.SessionID); sid != "" {
			sessionID = sid
			agentSID = sid
		}
		if kind := strings.TrimSpace(ag.Kind); kind != "" {
			runner = kind
		}
	}
	if runner == "" {
		runner = matchAgentLikeToken(sess)
	}

	sendable := sess.Idle != nil && *sess.Idle
	state := idleState(sess.Idle)
	live := true
	exited := false
	needs := live && !sendable && !exited

	iterm := "-"
	if !noITerm {
		iterm = formatLiveITerm(win, tab)
	}

	updated := "-"
	if sess.Duration != nil && strings.TrimSpace(*sess.Duration) != "" {
		updated = strings.TrimSpace(*sess.Duration)
	}

	workspace := ""
	if sess.Cwd != nil {
		workspace = *sess.Cwd
	}

	return listRow{
		Meta: SessionMetaView{
			SessionID: sessionID,
			Runner:    runner,
			Workspace: workspace,
			UpdatedAt: updated,
		},
		Probe: ProbeResult{
			Live:     live,
			Sendable: sendable,
			State:    state,
			Reason:   "",
			Exited:   exited,
			TTY:      sess.TTY,
		},
		ITerm:          iterm,
		NeedsAtt:       needs,
		AgentRun:       agentRunFromAgent(ag),
		AgentSessionID: agentSID,
	}
}

func idleState(idle *bool) string {
	if idle == nil {
		return "unknown"
	}
	if *idle {
		return "idle"
	}
	return "busy"
}

// formatLiveITerm builds w=<id> t=<tabIndex>.
// id = WindowID if non-zero else window Index; tabIndex = Tab.Index (1-based).
func formatLiveITerm(win *snapshot.SnapshotWindow, tab *snapshot.SnapshotTab) string {
	var id string
	if win.WindowID != 0 {
		id = fmt.Sprintf("%d", win.WindowID)
	} else {
		id = fmt.Sprintf("%d", win.Index)
	}
	return fmt.Sprintf("w=%s t=%d", id, tab.Index)
}
