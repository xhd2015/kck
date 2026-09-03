package run

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/agent-pro/pkgs/agenttty"
	"github.com/xhd2015/dot-pkgs/go-pkgs/cron/easycron"
	lessflags "github.com/xhd2015/less-flags"
)

const codexHelp = `Usage: kck codex <command> [ARGS]

Commands:
  list      list Codex ids hosted in iTerm tabs
  open      focus hosting tab or resume (--tab / --tab-index / <id>)
  focus     focus hosting iTerm tab only when live (no resume)
  snapshot  capture visible pane text (--tab / --tab-index / <id>)
  send      type text into hosting pane (--session-id / --tab / --open)
  messages  print recent chat messages (--limit / --grep / --offset-from-end)
  prompts   list user prompts (--first / --grep / --this-window / --tab)
  info      show session detail + Active block
  status    PID liveness + rollout path
  resolve   resolve Codex session id (ancestor walk or --tab)
  pickup    new empty session staged from a base session (kck-pickup-a-session)
  new       open a new empty Codex session via agent-run

Run 'kck codex <command> --help' for command-specific options.
`

const codexListHelp = `Usage: kck codex list [OPTIONS]

List Codex session ids that currently have a hosting iTerm2 tab.
Same discovery as: agent-pro codex session list-live.
Sessions with a live PID but no iTerm tab are omitted.

Options:
  --json        machine-readable JSON (no ANSI)
  --limit N     show at most N sessions (0 = unlimited)
  -h,--help     show help
`

const codexFocusHelp = `Usage: kck codex focus <session-id> [--index N]

Focus the iTerm2 tab that already hosts this live Codex session.
Lighter than open: never resumes or creates a window when no live host.

Options:
  --index N     select candidate N when multiple tabs host the same session
  -h,--help     show help
`

const codexOpenHelp = `Usage: kck codex open (<session-id> | --tab SEL | --tab-index N) [OPTIONS]

Focus the iTerm2 tab that already hosts this Codex session when one exists.
Otherwise open a new iTerm2 window and run: codex resume <session-id>
When the Codex id is bound in agent-run, prefers agent-run (live → focus;
exited → agent-run resume) instead of bare codex resume.

Session source (exactly one):
  <session-id>          explicit Codex session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --index N             select candidate N when multiple tabs host the same session
                        (positional <session-id> only; not with --tab/--tab-index)
  --dir DIR             workspace for resume (default: session cwd)
  --no-agent-run        force bare codex resume (skip agent-run prefer)
  --dry-run             resolve only; do not focus or open a window
  -h,--help             show help

A successful --tab/--tab-index resolve focuses that tab (never resumes).
`

const codexSnapshotHelp = `Usage: kck codex snapshot (<session-id> | --tab SEL | --tab-index N) [OPTIONS]

Print currently visible pane text for a live Codex session host.
Does not focus the pane. No resume when no host.

When the Codex id is bound to a live agent-run codex-tty session, prefers that
TTY snapshot (sanitized single frame). Otherwise uses iTerm2 Contents.
Bare codex (not under agent-run) always uses iTerm.

Session source (exactly one):
  <session-id>          explicit Codex session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --index N             select candidate N when multiple tabs host the same session
                        (positional <session-id> only; not with --tab/--tab-index)
  --json                emit {"session_id","iterm_session_id","app","source","contents"}
  -o, --output FILE     write output to FILE instead of stdout
  --dry-run             resolve only; do not capture pane text
  --iterm               force iTerm Contents (skip agent-run prefer path)
  -h,--help             show help
`

const codexSendHelp = `Usage: kck codex send [text] (--session-id <id> | --tab SEL | --tab-index N) [OPTIONS]

Type text and/or key sequences into the live iTerm2 pane that hosts a Codex session.
Same write-text path as: kool iterm2 session <iterm-uuid> send …
By default requires a hosting iTerm tab. With --open, resumes in a new
window when no host is found, waits for the tab to appear, then sends.
When the Codex id is bound in agent-run, --session-id prefers agent-run
auto-send-or-resume directly (live → send queue; exited → resume) with no
iTerm discovery or SendText. --tab / --tab-index still target iTerm panes.

Session source (exactly one):
  --session-id ID       Codex session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --index N             select candidate N when multiple tabs host the same session
                        (--session-id only; not with --tab/--tab-index)
  --no-submit           write without newline (stage; user presses Enter)
  --focus               switch to the session's window/tab before writing
  --no-ctrl-u           do not prefix Ctrl-U (default prefixes Ctrl-U)
  --open                if no hosting tab: resume in a new window, then send
  --no-agent-run        force iTerm path (skip agent-run prefer for --session-id)
  --dir DIR             workspace for --open resume (default: session cwd)
  --dry-run             resolve only; do not open or call SendText
  --enter               append Enter (\n) to the send sequence
  --up,--down,--left,--right   append arrow key (CSI)
  --esc                 append Escape
  --ctrl-c,--ctrl-d     append Ctrl-C / Ctrl-D
  --text STR            append text in sequence order (interleaves with keys)
  --cron EXPR           repeat send on an easy-cron schedule (foreground until done)
                        every-1h | every-1h-at-4m | every-5m-until-19h00m |
                        every-5m-not-within-19h00m-to-06h30m
                        With --dry-run: print next fire times + one would-send (no loop)
  -h,--help             show help

At least one of [text], --text, or a key flag is required.
Sequence flags keep CLI order and may repeat; positional [text] is always last.
--open cannot be combined with --tab/--tab-index.
`

const codexResolveHelp = `Usage: kck codex resolve [OPTIONS]

Resolve a Codex session id either by walking ancestors to the nearest
codex runner (default), or from a sibling iTerm2 tab in this window.

Options:
  --pid PID         start pid for ancestor walk (default: current process)
  --tab SEL         1-based tab index, or next|left|right (right ≡ next)
  --tab-index N     0-based tab index in this iTerm window
  --dry-run         print resolution plan ([dry-run] lines); same discovery path
  -v,--verbose      print detail fields on stderr (ancestor or tab)
  --json            print session id + detail fields as JSON
  -h,--help         show help

Exactly one session source: ancestor walk (default / --pid), or --tab, or
--tab-index. --pid cannot combine with --tab/--tab-index.
Relative next/left/right do not wrap; edges error.
Tab discovery matches: kool iterm2 window status.
Multiple unrelated Codex sessions on the same tab refuse.
`

const codexPromptsHelp = `Usage:
  kck codex prompts (<session-id> | --session-id ID | --tab SEL | --tab-index N | --this-tab)
    [--first] [--grep P]... [--exclude Q] [--head N | --tail N] [--max-body N]
    [--color|--no-color]
  kck codex prompts [--this-window | --this-space]
    [--first] [--recent <window>] [--limit N]
    [--grep P]... [--exclude Q] [--head N | --tail N] [--max-body N]
    [--color|--no-color]
  kck codex prompts [--recent <window>] [--limit N]
    [--first] [--grep P]... [--exclude Q] [--head N | --tail N] [--max-body N]
    [--color|--no-color]

Show user prompts only as compact lines (same shape as kck grok prompts).

Session source (exactly one when scoping):
  <session-id> / --session-id ID
  --tab SEL | --tab-index N | --this-tab
  --this-window / --this-space

Options:
  --first               only the first user prompt per session
  --grep P              repeatable; AND; case-insensitive literal
  --exclude Q           drop prompts matching Q
  --head N | --tail N   mutually exclusive with each other and --first
  --max-body N          soft-cap body runes + …
  --recent WINDOW       Nd|Nh|Nm
  --limit N             session cap (>= 1)
  --color / --no-color  force ANSI on/off
  -h,--help             show help
`

const codexMessagesHelp = `Usage: kck codex messages (<session-id> | --tab SEL | --tab-index N) [OPTIONS]

Print the most recent coalesced Codex chat messages (msgfmt-style),
with per-kind rune caps (user 4096, tool 128, thinking 512, response 8192).
Each line is prefixed with a local [YYYY-MM-DD HH:MM:SS] timestamp, or [—]
when the wire time is unknown. AGENTS/skills preambles are skipped.

Session source (exactly one):
  <session-id>          explicit Codex session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --limit N             page size (default 32; 0 = all remaining after offset)
  --offset-from-end N   skip N newest messages before applying --limit (default 0)
  --grep P              keep messages whose body contains P (repeatable; AND;
                        case-insensitive literal). Applied before offset/limit.
  --color               force ANSI color on (even when stdout is not a TTY)
  --no-color            force ANSI color off
  --json                machine-readable (includes total, offset, limit; no ANSI)
  -h,--help             show help
`

const codexInfoHelp = `Usage: kck codex info <session-id> [OPTIONS]

Show detailed info for one Codex session from ~/.codex (or $CODEX_HOME).
Appends an Active block (File always no; live PIDs when not --no-pid).

Options:
  --no-pid      skip live PID scan
  -h,--help     show help
`

const codexStatusHelp = `Usage: kck codex status <session-id> [OPTIONS]

Show PID liveness for one Codex session (open-file hard hits on codex runners).
Codex has no active_sessions.json: File is always no.
Also prints the rollout path (~-shortened in text).

State: running | inactive

Options:
  --no-pid      skip live PID scan
  --json        print SessionStatus as JSON (no ANSI; path is absolute)
  -h,--help     show help
`

func runCodex(opts Options) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	args := opts.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		txt := strings.TrimPrefix(codexHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	switch args[0] {
	case "list":
		return runCodexList(opts, args[1:])
	case "open":
		return runCodexOpen(opts, args[1:])
	case "focus":
		return runCodexFocus(opts, args[1:])
	case "snapshot":
		return runCodexSnapshot(opts, args[1:])
	case "send":
		return runCodexSend(opts, args[1:])
	case "messages":
		return runCodexMessages(opts, args[1:])
	case "prompts":
		return runCodexPrompts(opts, args[1:])
	case "info":
		return runCodexInfo(opts, args[1:])
	case "status":
		return runCodexStatus(opts, args[1:])
	case "resolve":
		return runCodexResolve(opts, args[1:])
	case "pickup":
		return runCodexPickup(opts, args[1:])
	case "new":
		return runCodexNew(opts, args[1:])
	default:
		return writeError(stderr, fmt.Sprintf("unknown codex command: %s", args[0]))
	}
}

func runCodexPrompts(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(codexPromptsHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	codexHome := resolveCodexHome(opts)
	promptOpts := opts.CodexPromptsOpts
	err := sessions.RunPrompts(args, stdout, stderr, codexHome, promptOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro codex session prompts", "kck codex prompts")
		return writeError(stderr, msg)
	}
	return nil
}

func runCodexMessages(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(codexMessagesHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	codexHome := resolveCodexHome(opts)
	msgOpts := opts.CodexMessagesOpts
	err := sessions.RunMessages(args, stdout, stderr, codexHome, msgOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro codex session messages", "kck codex messages")
		return writeError(stderr, msg)
	}
	return nil
}

func runCodexList(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	for _, a := range args {
		if a == "-h" || a == "--help" {
			txt := strings.TrimPrefix(codexListHelp, "\n")
			if !strings.HasSuffix(txt, "\n") {
				txt += "\n"
			}
			fmt.Fprint(stdout, txt)
			return nil
		}
	}

	codexHome := resolveCodexHome(opts)
	listOpts := opts.CodexListLiveOpts
	if listOpts == nil {
		listOpts = &sessions.ListLiveOpts{}
	}
	err := sessions.RunListLive(args, stdout, stderr, codexHome, listOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro codex session list-live", "kck codex list")
		return writeError(stderr, msg)
	}
	return nil
}

func runCodexOpen(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(codexOpenHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	codexHome := resolveCodexHome(opts)
	openOpts := opts.CodexOpenOpts
	err := sessions.RunOpen(args, stdout, stderr, codexHome, openOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro codex session open", "kck codex open")
		msg = strings.ReplaceAll(msg, "agent-pro codex session focus", "kck codex open")
		return writeError(stderr, msg)
	}
	return nil
}

func runCodexFocus(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(codexFocusHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	codexHome := resolveCodexHome(opts)
	focusOpts := opts.CodexFocusOpts
	err := sessions.RunFocus(args, stdout, codexHome, focusOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro codex session focus", "kck codex focus")
		return writeError(stderr, msg)
	}
	return nil
}

func runCodexSnapshot(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(codexSnapshotHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	codexHome := resolveCodexHome(opts)
	snapOpts := opts.CodexSnapshotOpts
	err := sessions.RunSnapshot(args, stdout, stderr, codexHome, snapOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro codex session snapshot", "kck codex snapshot")
		return writeError(stderr, msg)
	}
	return nil
}

func runCodexSend(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(codexSendHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	cronRaw, sendArgs, err := peelCronFlag(args)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	if cronRaw != "" {
		expr, perr := easycron.Parse(cronRaw)
		if perr != nil {
			return writeError(stderr, formatCronParseErr(perr))
		}
		if argsHaveDryRun(sendArgs) {
			return runCronDryPreviewFor(opts, cronRunnerCodex, expr, cronRaw, sendArgs)
		}
		return runCronSendLoopFor(opts, cronRunnerCodex, expr, cronRaw, sendArgs)
	}

	codexHome := resolveCodexHome(opts)
	sendOpts := opts.CodexSendOpts
	err = sessions.RunSend(sendArgs, stdout, stderr, codexHome, sendOpts)
	if err != nil {
		return writeError(stderr, rewriteSendErr(err))
	}
	return nil
}

func runCodexResolve(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(codexResolveHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	var resolveOpts *sessions.ResolveOpts
	if opts.CodexResolveOpts != nil {
		cp := *opts.CodexResolveOpts
		resolveOpts = &cp
	} else {
		resolveOpts = &sessions.ResolveOpts{}
	}
	resolveOpts.Stdout = stdout
	resolveOpts.Stderr = stderr
	if strings.TrimSpace(resolveOpts.CodexHome) == "" {
		resolveOpts.CodexHome = resolveCodexHome(opts)
	}

	err := sessions.RunResolve(args, resolveOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro codex session resolve", "kck codex resolve")
		return writeError(stderr, msg)
	}
	return nil
}

func runCodexInfo(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	var noPID bool
	remain, err := lessflags.Bool("--no-pid", &noPID).
		HelpFunc("-h,--help", func() {}).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if err == lessflags.ErrHelp {
			txt := strings.TrimPrefix(codexInfoHelp, "\n")
			if !strings.HasSuffix(txt, "\n") {
				txt += "\n"
			}
			fmt.Fprint(stdout, txt)
			return nil
		}
		return writeError(stderr, err.Error())
	}
	if len(remain) != 1 {
		return writeError(stderr, fmt.Sprintf("expected exactly one session id, got %d arguments", len(remain)))
	}
	sessionID := strings.TrimSpace(remain[0])
	if sessionID == "" {
		return writeError(stderr, "session id is required")
	}

	codexHome := resolveCodexHome(opts)
	info, err := sessions.Info(codexHome, sessionID, 3)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	now := opts.CodexNow
	if now.IsZero() {
		now = time.Now()
	}
	home, _ := os.UserHomeDir()
	fmt.Fprintln(stdout, sessions.FormatInfoText(info, home, now))

	st, err := sessions.Status(codexHome, sessionID, !noPID, opts.CodexLiveOpts)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, sessions.FormatActiveBlock(st))
	return nil
}

func runCodexStatus(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	var noPID bool
	var asJSON bool
	remain, err := lessflags.Bool("--no-pid", &noPID).
		Bool("--json", &asJSON).
		HelpFunc("-h,--help", func() {}).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if err == lessflags.ErrHelp {
			txt := strings.TrimPrefix(codexStatusHelp, "\n")
			if !strings.HasSuffix(txt, "\n") {
				txt += "\n"
			}
			fmt.Fprint(stdout, txt)
			return nil
		}
		return writeError(stderr, err.Error())
	}
	if len(remain) != 1 {
		return writeError(stderr, fmt.Sprintf("expected exactly one session id, got %d arguments", len(remain)))
	}
	sessionID := strings.TrimSpace(remain[0])
	if sessionID == "" {
		return writeError(stderr, "session id is required")
	}

	codexHome := resolveCodexHome(opts)
	st, err := sessions.Status(codexHome, sessionID, !noPID, opts.CodexLiveOpts)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	if asJSON {
		out, err := sessions.FormatStatusJSON(st)
		if err != nil {
			return writeError(stderr, err.Error())
		}
		fmt.Fprintln(stdout, out)
		return nil
	}
	fmt.Fprintln(stdout, sessions.FormatStatusText(st))
	return nil
}

func resolveCodexHome(opts Options) string {
	codexHome := strings.TrimSpace(opts.CodexHome)
	if codexHome == "" {
		return agenttty.CodexHome()
	}
	return codexHome
}
