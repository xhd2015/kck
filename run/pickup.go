package run

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	codexsessions "github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/agent-pro/pkgs/agenttty"
	"github.com/xhd2015/agent-pro/pkgs/shell"
	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
	lessflags "github.com/xhd2015/less-flags"
)

const grokPickupHelp = `Usage: kck grok pickup "msg..." (--session-id <id> | --tab SEL | --tab-index N) [OPTIONS]

Open a NEW empty Grok session and stage (do not submit) a kck-pickup-a-session
draft: pick up from the base session (read from the bottom) and continue with
the new instruction. Not a native fork.

By default runs in the current terminal. Use --new-window for a new iTerm2 window.
The embedded skill is hydrated to ~/.cache/kck-pickup-a-session/SKILL.md
(skipped when on-disk MD5 already matches the binary). Draft paths use ~.

Session source (exactly one):
  --session-id ID       base Grok session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --dir DIR             workspace for the new session (default: base cwd)
  --here                run in the current terminal (default)
  -n,--new-window,--new-terminal
                        open a new iTerm2 window instead
  --no-agent-run        bare grok + stage draft --no-submit (skip agent-run)
  --dry-run             resolve + print plan; do not launch
  -h,--help             show help

--here and --new-window cannot be combined.
`

const codexPickupHelp = `Usage: kck codex pickup "msg..." (--session-id <id> | --tab SEL | --tab-index N) [OPTIONS]

Open a NEW empty Codex session and stage (do not submit) a kck-pickup-a-session
draft: pick up from the base session (read from the bottom) and continue with
the new instruction. Not a native fork.

By default runs in the current terminal. Use --new-window for a new iTerm2 window.
The embedded skill is hydrated to ~/.cache/kck-pickup-a-session/SKILL.md
(skipped when on-disk MD5 already matches the binary). Draft paths use ~.

Session source (exactly one):
  --session-id ID       base Codex session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --dir DIR             workspace for the new session (default: base cwd)
  --here                run in the current terminal (default)
  -n,--new-window,--new-terminal
                        open a new iTerm2 window instead
  --no-agent-run        bare codex + stage draft --no-submit (skip agent-run)
  --dry-run             resolve + print plan; do not launch
  -h,--help             show help

--here and --new-window cannot be combined.
`

// PickupOpts injects resolve/launch hooks for L2. Nil → production.
type PickupOpts struct {
	ListProcs        func() []sessions.FocusProc
	Lsof             func(int) []string
	ListITerm        func() ([]iterm2.SessionRef, error)
	CurrentSessionID func() string
	ControllingTTY   func() string
	AncestorTTYs     func() []string

	// Codex tab-resolve probes (codex pickup only).
	CodexListProcs func() []codexsessions.FocusProc
	CodexLsof      func(int) []string

	// ResolveBaseCWD returns workspace for the base session when --dir is unset.
	// nil → sessions.Info / codexsessions.Info.
	ResolveBaseCWD func(home, sessionID string) (string, error)

	UserHomeDir func() (string, error)
	LookPath    func(file string) (string, error)
	// TildeHome shortens paths for draft display. nil → pathfmt.TildeHome.
	TildeHome func(path string) string

	// Skill cache hydrate injectables (nil → production).
	SkillContent   string // override embedded SKILL.md (tests)
	CacheSkillPath string // override ~/.cache/.../SKILL.md path (tests)
	ReadFile       func(path string) ([]byte, error)
	WriteFile      func(path string, data []byte, perm os.FileMode) error
	MkdirAll       func(path string, perm os.FileMode) error
	Rename         func(oldpath, newpath string) error

	// OpenInNewWindow opens an iTerm ForceNew window running followUp in dir.
	OpenInNewWindow func(dir, followUp string) error
	// RunForeground runs bin+argv in the current terminal (here + agent-run).
	RunForeground func(bin string, argv []string, dir string) error
	// RunBareHere runs bare grok/codex in the current terminal and stages draft.
	// nil → Start + sleep + StageDraft + Wait.
	RunBareHere func(bin, dir, draft string) error
	// StageDraft stages text into the front iTerm session without submit.
	// Used for --no-agent-run. nil → production AppleScript helper.
	StageDraft func(draft string) error
	// Sleep is used before bare StageDraft. nil → time.Sleep.
	Sleep func(time.Duration) error
}

type pickupArgs struct {
	Msg        string
	SessionID  *string
	Tab        *string
	TabIndex   *int
	Dir        string
	Here       bool
	NewWindow  bool
	NoAgentRun bool
	DryRun     bool
}

type pickupKind string

const (
	pickupGrok  pickupKind = "grok"
	pickupCodex pickupKind = "codex"
)

func runGrokPickup(opts Options, args []string) error {
	return runPickup(opts, args, pickupGrok)
}

func runCodexPickup(opts Options, args []string) error {
	return runPickup(opts, args, pickupCodex)
}

func runPickup(opts Options, args []string, kind pickupKind) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		help := grokPickupHelp
		if kind == pickupCodex {
			help = codexPickupHelp
		}
		txt := strings.TrimPrefix(help, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	parsed, err := parsePickupArgs(args)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	if strings.TrimSpace(parsed.Msg) == "" {
		return writeError(stderr, "message is required")
	}
	if parsed.Here && parsed.NewWindow {
		return writeError(stderr, "--here and --new-window cannot be combined")
	}
	newWindow := parsed.NewWindow // default: here (current terminal)

	home := pickupRunnerHome(opts, kind)
	popts := pickupOptsFor(opts, kind)
	applyPickupDefaults(popts, kind, home)

	sessionID, err := resolvePickupBaseID(parsed, home, kind, popts)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	skillPath, err := ensurePickupSkillCached(popts)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	tilde := popts.TildeHome
	if tilde == nil {
		tilde = pathfmt.TildeHome
	}
	skillDisplay := tilde(skillPath)

	cwd, err := resolvePickupCWD(home, sessionID, parsed.Dir, popts)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	draft := buildPickupDraft(skillDisplay, sessionID, parsed.Msg)

	if parsed.DryRun {
		fmt.Fprintf(stdout, "Would pickup from %s\n", sessionID)
		fmt.Fprintf(stdout, "dir: %s\n", cwd)
		if newWindow {
			fmt.Fprintln(stdout, "terminal: new iTerm2 window")
		} else {
			fmt.Fprintln(stdout, "terminal: current")
		}
		fmt.Fprintf(stdout, "draft: %s\n", draft)
		return nil
	}

	if err := launchPickup(kind, cwd, draft, parsed.NoAgentRun, newWindow, popts); err != nil {
		return writeError(stderr, err.Error())
	}

	if newWindow {
		fmt.Fprintf(stdout, "opened: new window; pickup from session %s\n", sessionID)
	} else {
		fmt.Fprintf(stdout, "opened: here; pickup from session %s\n", sessionID)
	}
	return nil
}

func pickupRunnerHome(opts Options, kind pickupKind) string {
	if kind == pickupCodex {
		home := strings.TrimSpace(opts.CodexHome)
		if home == "" {
			home = agenttty.CodexHome()
		}
		return home
	}
	home := strings.TrimSpace(opts.GrokHome)
	if home == "" {
		home = agenttty.GrokHome()
	}
	return home
}

func pickupOptsFor(opts Options, kind pickupKind) *PickupOpts {
	var popts *PickupOpts
	if kind == pickupCodex {
		popts = opts.CodexPickupOpts
	} else {
		popts = opts.GrokPickupOpts
	}
	if popts == nil {
		popts = &PickupOpts{}
	}
	return popts
}

func applyPickupDefaults(popts *PickupOpts, kind pickupKind, home string) {
	if popts.ResolveBaseCWD != nil {
		return
	}
	if kind == pickupCodex {
		popts.ResolveBaseCWD = func(h, sessionID string) (string, error) {
			if strings.TrimSpace(h) == "" {
				h = home
			}
			info, err := codexsessions.Info(h, sessionID, 1)
			if err != nil {
				return "", err
			}
			return strings.TrimSpace(info.CWD), nil
		}
		return
	}
	popts.ResolveBaseCWD = func(h, sessionID string) (string, error) {
		if strings.TrimSpace(h) == "" {
			h = home
		}
		info, err := sessions.Info(h, sessionID)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(info.CWD), nil
	}
}

func parsePickupArgs(args []string) (pickupArgs, error) {
	var out pickupArgs
	remain, err := lessflags.
		String("--session-id", &out.SessionID).
		String("--tab", &out.Tab).
		Int("--tab-index", &out.TabIndex).
		String("--dir", &out.Dir).
		Bool("--here", &out.Here).
		Bool("-n,--new-window,--new-terminal", &out.NewWindow).
		Bool("--no-agent-run", &out.NoAgentRun).
		Bool("--dry-run", &out.DryRun).
		HelpFunc("-h,--help", func() {}).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if err == lessflags.ErrHelp {
			return out, nil
		}
		return out, err
	}
	switch len(remain) {
	case 0:
		// missing message checked by caller
	case 1:
		out.Msg = remain[0]
	default:
		return out, fmt.Errorf("pickup: unexpected arguments: %s", strings.Join(remain[1:], " "))
	}
	return out, nil
}

func resolvePickupBaseID(parsed pickupArgs, home string, kind pickupKind, popts *PickupOpts) (string, error) {
	if parsed.SessionID != nil {
		if parsed.Tab != nil || parsed.TabIndex != nil {
			return "", fmt.Errorf("--session-id cannot be combined with --tab/--tab-index")
		}
		id := strings.TrimSpace(*parsed.SessionID)
		if id == "" {
			return "", fmt.Errorf("session id is required")
		}
		return id, nil
	}
	if parsed.Tab == nil && parsed.TabIndex == nil {
		return "", fmt.Errorf("expected --session-id, or --tab / --tab-index")
	}
	if kind == pickupCodex {
		return resolveCodexPickupTab(parsed.Tab, parsed.TabIndex, home, popts)
	}
	return resolveGrokPickupTab(parsed.Tab, parsed.TabIndex, home, popts)
}

func resolveGrokPickupTab(tabFlag *string, tabIndexFlag *int, grokHome string, popts *PickupOpts) (string, error) {
	id, _, err := sessions.ResolveSessionSource(nil, tabFlag, tabIndexFlag, &sessions.SessionSourceOpts{
		ListProcs:        popts.ListProcs,
		Lsof:             popts.Lsof,
		ListITerm:        popts.ListITerm,
		CurrentSessionID: popts.CurrentSessionID,
		ControllingTTY:   popts.ControllingTTY,
		AncestorTTYs:     popts.AncestorTTYs,
		GrokHome:         grokHome,
	})
	return id, err
}

func resolveCodexPickupTab(tabFlag *string, tabIndexFlag *int, _ string, popts *PickupOpts) (string, error) {
	id, _, err := codexsessions.ResolveSessionSource(nil, tabFlag, tabIndexFlag, &codexsessions.SessionSourceOpts{
		ListProcs:        popts.CodexListProcs,
		Lsof:             popts.CodexLsof,
		ListITerm:        popts.ListITerm,
		CurrentSessionID: popts.CurrentSessionID,
		ControllingTTY:   popts.ControllingTTY,
		AncestorTTYs:     popts.AncestorTTYs,
	})
	return id, err
}

func resolvePickupCWD(home, sessionID, dirOverride string, popts *PickupOpts) (string, error) {
	if d := strings.TrimSpace(dirOverride); d != "" {
		return absExistingDir(d)
	}
	cwd, err := popts.ResolveBaseCWD(home, sessionID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(cwd) == "" {
		return "", fmt.Errorf("session %s has empty cwd; pass --dir", sessionID)
	}
	return absExistingDir(cwd)
}

func absExistingDir(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("workspace dir: %w", err)
	}
	if real, e := filepath.EvalSymlinks(abs); e == nil {
		abs = real
	}
	st, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("workspace dir: %w", err)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("workspace dir: not a directory: %s", abs)
	}
	return abs, nil
}

func buildPickupDraft(skillPath, sessionID, msg string) string {
	return fmt.Sprintf("read %s, session-id: %s, %s", skillPath, sessionID, msg)
}

func launchPickup(kind pickupKind, cwd, draft string, noAgentRun, newWindow bool, popts *PickupOpts) error {
	lookPath := popts.LookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	runFg := popts.RunForeground
	if runFg == nil {
		runFg = defaultPickupRunForeground
	}
	openFn := popts.OpenInNewWindow
	if openFn == nil {
		openFn = defaultPickupOpenInNewWindow
	}

	if !noAgentRun {
		bin, err := lookPath("agent-run")
		if err != nil {
			return fmt.Errorf("agent-run not found on PATH: %w", err)
		}
		runner := "grok-tty"
		if kind == pickupCodex {
			runner = "codex-tty"
		}
		argv := buildAgentRunPickupArgv(runner, cwd, draft)
		if newWindow {
			cmdLine := quotedPickupCommand(bin, argv)
			if err := openFn(cwd, cmdLine); err != nil {
				return fmt.Errorf("open new window: %w", err)
			}
			return nil
		}
		if err := runFg(bin, argv, cwd); err != nil {
			return fmt.Errorf("run here: %w", err)
		}
		return nil
	}

	runnerBin := "grok"
	if kind == pickupCodex {
		runnerBin = "codex"
	}
	bin, err := lookPath(runnerBin)
	if err != nil {
		return fmt.Errorf("%s not found on PATH: %w", runnerBin, err)
	}

	stage := popts.StageDraft
	if stage == nil {
		stage = defaultStageDraftInFrontSession
	}
	sleepFn := popts.Sleep
	if sleepFn == nil {
		sleepFn = func(d time.Duration) error {
			time.Sleep(d)
			return nil
		}
	}

	if newWindow {
		cmdLine := shell.ShellQuote(bin)
		if err := openFn(cwd, cmdLine); err != nil {
			return fmt.Errorf("open new window: %w", err)
		}
		if err := sleepFn(2 * time.Second); err != nil {
			return err
		}
		if err := stage(draft); err != nil {
			return fmt.Errorf("stage draft: %w", err)
		}
		return nil
	}

	if popts.RunBareHere != nil {
		return popts.RunBareHere(bin, cwd, draft)
	}
	return defaultPickupRunBareHere(bin, cwd, draft, stage, sleepFn)
}

func defaultPickupRunBareHere(bin, dir, draft string, stage func(string) error, sleepFn func(time.Duration) error) error {
	cmd := exec.Command(bin)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start here: %w", err)
	}
	if err := sleepFn(2 * time.Second); err != nil {
		_ = cmd.Process.Kill()
		return err
	}
	if err := stage(draft); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("stage draft: %w", err)
	}
	return cmd.Wait()
}

func buildAgentRunPickupArgv(runner, dir, draft string) []string {
	argv := []string{
		"run",
		"--open",
		"--no-submit",
		"--agent-runner", runner,
	}
	if strings.TrimSpace(dir) != "" {
		argv = append(argv, "--dir", dir)
	}
	if p := strings.TrimSpace(draft); p != "" {
		argv = append(argv, "--", p)
	}
	return argv
}

func quotedPickupCommand(bin string, argv []string) string {
	parts := append([]string{bin}, argv...)
	quoted := make([]string, 0, len(parts))
	for _, p := range parts {
		quoted = append(quoted, shell.ShellQuote(p))
	}
	return strings.Join(quoted, " ")
}

func defaultPickupOpenInNewWindow(dir, followUp string) error {
	return iterm2.OpenConfig(dir, &iterm2.Config{
		Mode:             iterm2.ModeForceNew,
		FollowUpCommands: []string{followUp},
		SafeInputIgnore:  true,
	})
}

func defaultPickupRunForeground(bin string, argv []string, dir string) error {
	cmd := exec.Command(bin, argv...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func defaultStageDraftInFrontSession(draft string) error {
	// Best-effort: type into the frontmost iTerm session without newline.
	escaped := iterm2.EscapeCommandForAppleScript(draft)
	script := strings.Join([]string{
		`tell application "iTerm"`,
		`  tell current session of current window`,
		`    write text ((ASCII character 21) & "` + escaped + `") without newline`,
		`  end tell`,
		`end tell`,
	}, "\n")
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
