package run

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/xhd2015/agent-pro/pkgs/agentrunapi"
	lessflags "github.com/xhd2015/less-flags"
)

const (
	grokNewHelp = `Usage: kck grok new "prompt..." [OPTIONS]

Open a NEW empty Grok session via agent-run (wrk create agent-launch shape).
Prepends /brainstorm to the prompt. Runner: grok-tty.

By default opens a new iTerm2 window (--new-terminal), waits until the
Grok provider session id is resolved (registry PID → open-files), prints it,
then exits. Use --here / --no-new-terminal to run in the current terminal
(silent; no kck output so the TUI is not corrupted).

Options:
  --dir DIR             workspace (default: cwd)
  --here                run in the current terminal (silent)
  --no-new-terminal     same as --here
  --new-terminal        open a new iTerm2 window (default)
  --submit              auto-submit prompt (default: stage with --no-submit)
  --dry-run             print plan; do not launch
  -h,--help             show help

--here/--no-new-terminal and --new-terminal cannot be combined.
`

	codexNewHelp = `Usage: kck codex new "prompt..." [OPTIONS]

Open a NEW empty Codex session via agent-run (wrk create agent-launch shape).
Prepends $brainstorm to the prompt. Runner: codex-tty.

By default opens a new iTerm2 window (--new-terminal), waits until the
Codex provider session id is resolved (registry PID → open-files), prints it,
then exits. Use --here / --no-new-terminal to run in the current terminal
(silent; no kck output so the TUI is not corrupted).

Options:
  --dir DIR             workspace (default: cwd)
  --here                run in the current terminal (silent)
  --no-new-terminal     same as --here
  --new-terminal        open a new iTerm2 window (default)
  --submit              auto-submit prompt (default: stage with --no-submit)
  --dry-run             print plan; do not launch
  -h,--help             show help

--here/--no-new-terminal and --new-terminal cannot be combined.
`

	newSessionIDBaseMaxRunes = 128
	newSessionIDFallbackBase = "sess"
)

// NewOpts injects launch/wait hooks for L2. Nil → production.
type NewOpts struct {
	LookPath func(file string) (string, error)
	// RunForeground runs bin+argv in the current terminal (here path).
	RunForeground func(bin string, argv []string, dir string) error
	// RunNewTerminal runs bin+argv that includes agent-run --new-terminal.
	// nil → same as RunForeground (agent-run owns ForceNew).
	RunNewTerminal func(bin string, argv []string, dir string) error
	// WaitProviderSession waits for the grok/codex provider session id after
	// new-terminal launch. nil → agentrunapi.WaitProviderSessionID.
	WaitProviderSession func(home, runner, agentRunSessionID string) (providerSessionID string, err error)
	// AgentRunHome overrides AGENT_RUN_HOME / ~/.agent-run (tests + wait path).
	AgentRunHome string
	Getwd        func() (string, error)
	Abs          func(string) (string, error)
	Stat         func(string) (os.FileInfo, error)
	UserHomeDir  func() (string, error)
	Getenv       func(string) string
	Now          func() time.Time
	Sleep        func(time.Duration) error
	// SessionExists reports whether an agent-run session id is taken.
	// nil → sessions/<id>/meta.json or sessions/<id>/ exists.
	SessionExists func(home, sessionID string) bool
}

type newArgs struct {
	Msg            string
	Dir            string
	Here           bool
	NoNewTerminal  bool
	NewTerminal    bool
	Submit         bool
	DryRun         bool
}

type newKind string

const (
	newGrok  newKind = "grok"
	newCodex newKind = "codex"
)

func runGrokNew(opts Options, args []string) error {
	return runNew(opts, args, newGrok)
}

func runCodexNew(opts Options, args []string) error {
	return runNew(opts, args, newCodex)
}

func runNew(opts Options, args []string, kind newKind) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		help := grokNewHelp
		if kind == newCodex {
			help = codexNewHelp
		}
		txt := strings.TrimPrefix(help, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	parsed, err := parseNewArgs(args)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	if strings.TrimSpace(parsed.Msg) == "" {
		return writeError(stderr, "message is required")
	}

	here := parsed.Here || parsed.NoNewTerminal
	if here && parsed.NewTerminal {
		return writeError(stderr, "--here/--no-new-terminal and --new-terminal cannot be combined")
	}
	newTerminal := !here // default: new terminal

	nopts := newOptsFor(opts, kind)
	cwd, err := resolveNewCWD(parsed.Dir, nopts)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	prompt := buildNewPrompt(kind, parsed.Msg)
	runner := "grok-tty"
	if kind == newCodex {
		runner = "codex-tty"
	}

	home, err := resolveNewAgentRunHome(opts, nopts)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	sessionID, err := allocateNewSessionID(prompt, home, nopts)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	argv := buildAgentRunNewArgv(runner, cwd, sessionID, prompt, parsed.Submit, newTerminal)

	if parsed.DryRun {
		label := "grok"
		if kind == newCodex {
			label = "codex"
		}
		fmt.Fprintf(stdout, "Would open new %s session\n", label)
		if newTerminal {
			fmt.Fprintln(stdout, "terminal: new")
		} else {
			fmt.Fprintln(stdout, "terminal: current")
		}
		fmt.Fprintf(stdout, "dir: %s\n", cwd)
		fmt.Fprintf(stdout, "runner: %s\n", runner)
		fmt.Fprintf(stdout, "agent-run-session-id: %s\n", sessionID)
		fmt.Fprintf(stdout, "prompt: %s\n", prompt)
		fmt.Fprintf(stdout, "submit: %v\n", parsed.Submit)
		fmt.Fprintf(stdout, "cmd: agent-run %s\n", strings.Join(argv, " "))
		return nil
	}

	lookPath := nopts.LookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	bin, err := lookPath("agent-run")
	if err != nil {
		return writeError(stderr, fmt.Sprintf("agent-run not found on PATH: %v", err))
	}

	if here {
		runFg := nopts.RunForeground
		if runFg == nil {
			runFg = defaultNewRunForeground
		}
		// Silent: any kck stdout/stderr would corrupt the TUI.
		if err := runFg(bin, argv, cwd); err != nil {
			return err
		}
		return nil
	}

	runNT := nopts.RunNewTerminal
	if runNT == nil {
		runNT = nopts.RunForeground
	}
	if runNT == nil {
		runNT = defaultNewRunForeground
	}
	if err := runNT(bin, argv, cwd); err != nil {
		return writeError(stderr, fmt.Sprintf("run new-terminal: %v", err))
	}

	wait := nopts.WaitProviderSession
	if wait == nil {
		wait = defaultWaitProviderSession(nopts)
	}
	providerID, err := wait(home, runner, sessionID)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	label := "grok"
	if kind == newCodex {
		label = "codex"
	}
	fmt.Fprintf(stdout, "opened: new terminal; new %s session\n", label)
	fmt.Fprintf(stdout, "session-id: %s\n", providerID)
	return nil
}

func newOptsFor(opts Options, kind newKind) *NewOpts {
	var nopts *NewOpts
	if kind == newCodex {
		nopts = opts.CodexNewOpts
	} else {
		nopts = opts.GrokNewOpts
	}
	if nopts == nil {
		nopts = &NewOpts{}
	}
	return nopts
}

func parseNewArgs(args []string) (newArgs, error) {
	var out newArgs
	remain, err := lessflags.
		String("--dir", &out.Dir).
		Bool("--here", &out.Here).
		Bool("--no-new-terminal", &out.NoNewTerminal).
		Bool("--new-terminal", &out.NewTerminal).
		Bool("--submit", &out.Submit).
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
		return out, fmt.Errorf("unexpected arguments: %s", strings.Join(remain[1:], " "))
	}
	return out, nil
}

func buildNewPrompt(kind newKind, msg string) string {
	msg = strings.TrimSpace(msg)
	if kind == newCodex {
		return "$brainstorm " + msg
	}
	return "/brainstorm " + msg
}

func buildAgentRunNewArgv(runner, dir, sessionID, prompt string, submit, newTerminal bool) []string {
	argv := []string{
		"run",
		"--open",
		"--color",
		"--agent-runner", runner,
		"--session-id", sessionID,
	}
	if !submit {
		argv = append(argv, "--no-submit")
	}
	if newTerminal {
		argv = append(argv, "--new-terminal")
	}
	if strings.TrimSpace(dir) != "" {
		argv = append(argv, "--dir", dir)
	}
	if p := strings.TrimSpace(prompt); p != "" {
		argv = append(argv, "--", p)
	}
	return argv
}

func resolveNewCWD(dirOverride string, nopts *NewOpts) (string, error) {
	getwd := nopts.Getwd
	if getwd == nil {
		getwd = os.Getwd
	}
	absFn := nopts.Abs
	if absFn == nil {
		absFn = filepath.Abs
	}
	statFn := nopts.Stat
	if statFn == nil {
		statFn = os.Stat
	}

	raw := strings.TrimSpace(dirOverride)
	if raw == "" {
		wd, err := getwd()
		if err != nil {
			return "", fmt.Errorf("workspace dir: %w", err)
		}
		raw = wd
	}
	abs, err := absFn(raw)
	if err != nil {
		return "", fmt.Errorf("workspace dir: %w", err)
	}
	if real, e := filepath.EvalSymlinks(abs); e == nil {
		abs = real
	}
	st, err := statFn(abs)
	if err != nil {
		return "", fmt.Errorf("workspace dir: %w", err)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("workspace dir: not a directory: %s", abs)
	}
	return abs, nil
}

func resolveNewAgentRunHome(opts Options, nopts *NewOpts) (string, error) {
	if h := strings.TrimSpace(nopts.AgentRunHome); h != "" {
		return filepath.Clean(h), nil
	}
	if h := strings.TrimSpace(opts.Home); h != "" {
		return filepath.Clean(h), nil
	}
	getenv := nopts.Getenv
	if getenv == nil {
		getenv = os.Getenv
	}
	if v := strings.TrimSpace(getenv("AGENT_RUN_HOME")); v != "" {
		return filepath.Clean(v), nil
	}
	userHome := nopts.UserHomeDir
	if userHome == nil {
		userHome = os.UserHomeDir
	}
	dir, err := userHome()
	if err != nil {
		return "", fmt.Errorf("agent-run home: %w", err)
	}
	return filepath.Clean(filepath.Join(dir, ".agent-run")), nil
}

func allocateNewSessionID(prompt, home string, nopts *NewOpts) (string, error) {
	nowFn := nopts.Now
	if nowFn == nil {
		nowFn = time.Now
	}
	exists := nopts.SessionExists
	if exists == nil {
		exists = defaultNewSessionExists
	}
	base := slugifyNewPrompt(prompt)
	ts := nowFn().Local().Format("20060102-150405")
	candidate := base + "-" + ts
	if !exists(home, candidate) {
		return candidate, nil
	}
	for n := 1; n < 10000; n++ {
		id := fmt.Sprintf("%s-%d", candidate, n)
		if !exists(home, id) {
			return id, nil
		}
	}
	return "", fmt.Errorf("could not allocate free session id for base %q", base)
}

func slugifyNewPrompt(prompt string) string {
	s := strings.ToLower(prompt)
	var b strings.Builder
	b.Grow(len(s))
	prevDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash {
			b.WriteByte('-')
			prevDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return newSessionIDFallbackBase
	}
	if utf8.RuneCountInString(out) > newSessionIDBaseMaxRunes {
		runes := []rune(out)
		out = strings.TrimRight(string(runes[:newSessionIDBaseMaxRunes]), "-")
		if out == "" {
			return newSessionIDFallbackBase
		}
	}
	return out
}

func defaultNewSessionExists(home, sessionID string) bool {
	metaPath := filepath.Join(home, "sessions", sessionID, "meta.json")
	if _, err := os.Stat(metaPath); err == nil {
		return true
	}
	dir := filepath.Join(home, "sessions", sessionID)
	if st, err := os.Stat(dir); err == nil && st.IsDir() {
		return true
	}
	return false
}

func defaultWaitProviderSession(nopts *NewOpts) func(home, runner, agentRunSessionID string) (string, error) {
	return func(home, runner, agentRunSessionID string) (string, error) {
		res, err := agentrunapi.WaitProviderSessionID(agentrunapi.WaitProviderSessionOpts{
			Home:      home,
			Runner:    runner,
			SessionID: agentRunSessionID,
			Sleep:     nopts.Sleep,
		})
		if err != nil {
			return "", err
		}
		return res.ProviderSessionID, nil
	}
}

func defaultNewRunForeground(bin string, argv []string, dir string) error {
	cmd := exec.Command(bin, argv...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
