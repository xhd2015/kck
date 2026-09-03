package run

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/agent-pro/pkgs/agenttty"
	"github.com/xhd2015/dot-pkgs/go-pkgs/cron/easycron"
	lessflags "github.com/xhd2015/less-flags"
)

const grokHelp = `Usage: kck grok <command> [ARGS]

Commands:
  list      list Grok session ids hosted in iTerm tabs
  open      focus hosting iTerm tab, or resume (--tab / --tab-index / <id>)
  focus     focus hosting iTerm tab only when live (no resume)
  snapshot  capture visible pane text (--tab / --tab-index / <id>)
  send      type text into hosting pane (--session-id / --tab / --open)
  messages  print recent chat messages (--limit / --grep / --offset-from-end)
  prompts   list user prompts (--first / --main / --grep / --this-window / --tab)
  info      show session detail + Active block
  status    dual-signal liveness + session path
  wait      block until current turn finishes (updates.jsonl)
  resolve   resolve Grok session id (ancestor walk or --tab)
  pickup    new empty session staged from a base session (kck-pickup-a-session)
  new       open a new empty Grok session via agent-run

Run 'kck grok <command> --help' for command-specific options.
`

const grokPromptsHelp = `Usage:
  kck grok prompts (<session-id> | --session-id ID | --tab SEL | --tab-index N | --this-tab)
    [--first] [--grep P]... [--exclude Q] [--head N | --tail N] [--max-body N]
    [--color|--no-color]
  kck grok prompts [--this-window | --this-space]
    [--first] [--recent <window>] [--limit N]
    [--grep P]... [--exclude Q] [--head N | --tail N] [--max-body N]
    [--color|--no-color]
  kck grok prompts [--recent <window>] [--limit N]
    [--first] [--grep P]... [--exclude Q] [--head N | --tail N] [--max-body N]
    [--color|--no-color]

Show user prompts only as compact lines:
  [YYYY-MM-DD HH:MM:SS] prompt text…

Single mode: all user prompts for one session (full history), optional text filters.
Multi mode (no session source): newest sessions by last_active:

  (no flags)              last 10 sessions that have prompts
  --limit N               last N sessions that have prompts (N >= 1)
  --recent Nd|Nh|Nm       all sessions with ≥1 in-window user prompt
  --this-window           live hosts in this iTerm window (no default cap)
  --this-space            live hosts on this macOS Mission Control desktop

Session source (exactly one when scoping):
  <session-id> / --session-id ID
  --tab SEL             1-based index, or next|left|right|current
  --tab-index N         0-based tab index
  --this-tab            alias for --tab current
  --this-window / --this-space

Options:
  --first               only the first user prompt per session
  --main                only main-agent class sessions (alias: --main-agent)
  --grep P              repeatable; AND; case-insensitive literal
  --exclude Q           drop prompts matching Q
  --head N | --tail N   mutually exclusive with each other and --first
  --max-body N          soft-cap body runes + …
  --recent WINDOW       Nd|Nh|Nm
  --limit N             session cap (>= 1)
  --color / --no-color  force ANSI on/off
  -h,--help             show help
`

const grokMessagesHelp = `Usage: kck grok messages (<session-id> | --tab SEL | --tab-index N) [OPTIONS]

Print the most recent coalesced Grok chat messages (msgfmt-style),
with per-kind rune caps (user 4096, tool 128, thinking 512, response 8192).
Each line is prefixed with a local [YYYY-MM-DD HH:MM:SS] timestamp, or [—]
when the wire time is unknown.

Session source (exactly one):
  <session-id>          explicit Grok session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --limit N             page size (default 32; 0 = all remaining after offset)
  --offset-from-end N   skip N newest messages before applying --limit (default 0)
                        example: --offset-from-end 32  # skip last 32; start next page
  --grep P              keep messages whose body contains P (repeatable; AND;
                        case-insensitive literal). Applied before offset/limit.
  --color               force ANSI color on (even when stdout is not a TTY)
  --no-color            force ANSI color off
  --json                machine-readable (includes total, offset, limit; no ANSI)
  -h,--help             show help
`

const grokListHelp = `Usage: kck grok list [OPTIONS]

List Grok session ids that currently have a hosting iTerm2 tab.
Same discovery as: agent-pro grok session list-live.
Sessions with a live PID but no iTerm tab are omitted.

Options:
  --json        machine-readable JSON (no ANSI)
  --limit N     show at most N sessions (0 = unlimited)
  -h,--help     show help
`

const grokFocusHelp = `Usage: kck grok focus <session-id> [--index N]

Focus the iTerm2 tab that already hosts this live Grok session.
Lighter than open: never resumes or creates a window when no live host.

Options:
  --index N     select candidate N when multiple tabs host the same session
  -h,--help     show help
`

const grokOpenHelp = `Usage: kck grok open (<session-id> | --tab SEL | --tab-index N) [OPTIONS]

Focus the iTerm2 tab that already hosts this Grok session when one exists.
Otherwise open a new iTerm2 window and run: grok --resume <session-id>
When the Grok id is bound in agent-run, prefers agent-run (live → focus;
exited → agent-run resume) instead of bare grok --resume.

Session source (exactly one):
  <session-id>          explicit Grok session id
  --tab SEL             1-based tab index, or next|left|right (right ≡ next)
  --tab-index N         0-based tab index in this iTerm window

Options:
  --index N             select candidate N when multiple tabs host the same session
                        (positional <session-id> only; not with --tab/--tab-index)
  --dir DIR             workspace for resume (default: session cwd)
  --no-agent-run        force bare grok --resume (skip agent-run prefer)
  --dry-run             resolve only; do not focus or open a window
  -h,--help             show help

A successful --tab/--tab-index resolve focuses that tab (never resumes).
`

const grokSnapshotHelp = `Usage: kck grok snapshot (<session-id> | --tab SEL | --tab-index N) [OPTIONS]

Print currently visible pane text for a live Grok session host.
Does not focus the pane. No resume when no host.

When the Grok id is bound to a live agent-run grok-tty session, prefers that
TTY snapshot (sanitized single frame). Otherwise uses iTerm2 Contents.
Bare grok (not under agent-run) always uses iTerm.

Session source (exactly one):
  <session-id>          explicit Grok session id
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

const grokSendHelp = `Usage: kck grok send [text] (--session-id <id> | --tab SEL | --tab-index N) [OPTIONS]

Type text and/or key sequences into the live iTerm2 pane that hosts a Grok session.
Same write-text path as: kool iterm2 session <iterm-uuid> send …
By default requires a hosting iTerm tab. With --open, resumes in a new
window when no host is found, waits for the tab to appear, then sends.
When the Grok id is bound in agent-run, --session-id prefers agent-run
auto-send-or-resume directly (live → send queue; exited → resume) with no
iTerm discovery or SendText. --tab / --tab-index still target iTerm panes.

Session source (exactly one):
  --session-id ID       Grok session id
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

const grokInfoHelp = `Usage: kck grok info <session-id> [OPTIONS]

Show detailed info for one Grok session from ~/.grok (or $GROK_HOME).
Appends a dual-signal Active block (file-active + live PIDs).

Options:
  --no-pid      skip live PID scan; Active state from file-active only
  -h,--help     show help
`

const grokStatusHelp = `Usage: kck grok status <session-id> [OPTIONS]

Show dual-signal liveness for one Grok session:
  file-active (active_sessions.json) + live PIDs (open-file hard hits).
Also prints the session summary.json path (~-shortened in text).

State: running | marked-active | inactive

Options:
  --no-pid      skip live PID scan; state from file-active only
  --json        print SessionStatus as JSON (no ANSI; path is absolute)
  -h,--help     show help
`

const grokWaitHelp = `Usage: kck grok wait <session-id> [OPTIONS]

Block until the current turn finishes, or error if the session is not running.

Turn state is read from updates.jsonl (user_message_chunk vs turn_completed),
not from screen/TTY idle. Mid-turn waits for turn_completed; already outside
a turn returns immediately while the session stays running.

Options:
  --timeout DUR   max wait (default 30m; Go duration, e.g. 30s, 5m, 1h)
  -h,--help       show help
`

const grokResolveHelp = `Usage: kck grok resolve [OPTIONS]

Resolve a Grok session id either by walking ancestors to the nearest
grok runner (default), or from a sibling iTerm2 tab in this window.

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
When a parent and its child subagent share a tab, the parent id is returned;
unrelated multiple grok sessions on the same tab still refuse.
`

func runGrok(opts Options) error {
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
		txt := strings.TrimPrefix(grokHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	switch args[0] {
	case "list":
		return runGrokList(opts, args[1:])
	case "open":
		return runGrokOpen(opts, args[1:])
	case "focus":
		return runGrokFocus(opts, args[1:])
	case "snapshot":
		return runGrokSnapshot(opts, args[1:])
	case "send":
		return runGrokSend(opts, args[1:])
	case "messages":
		return runGrokMessages(opts, args[1:])
	case "prompts":
		return runGrokPrompts(opts, args[1:])
	case "info":
		return runGrokInfo(opts, args[1:])
	case "status":
		return runGrokStatus(opts, args[1:])
	case "wait":
		return runGrokWait(opts, args[1:])
	case "resolve":
		return runGrokResolve(opts, args[1:])
	case "pickup":
		return runGrokPickup(opts, args[1:])
	case "new":
		return runGrokNew(opts, args[1:])
	default:
		return writeError(stderr, fmt.Sprintf("unknown grok command: %s", args[0]))
	}
}

func runGrokMessages(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(grokMessagesHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	grokHome := resolveGrokHome(opts)
	msgOpts := opts.GrokMessagesOpts
	err := sessions.RunMessages(args, stdout, stderr, grokHome, msgOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro grok session messages", "kck grok messages")
		return writeError(stderr, msg)
	}
	return nil
}

func runGrokPrompts(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(grokPromptsHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	grokHome := resolveGrokHome(opts)
	promptOpts := opts.GrokPromptsOpts
	err := sessions.RunPrompts(args, stdout, stderr, grokHome, promptOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro grok session prompts", "kck grok prompts")
		return writeError(stderr, msg)
	}
	return nil
}

func runGrokList(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	// Local help text when -h/--help alone (RunListLive also handles help).
	for _, a := range args {
		if a == "-h" || a == "--help" {
			txt := strings.TrimPrefix(grokListHelp, "\n")
			if !strings.HasSuffix(txt, "\n") {
				txt += "\n"
			}
			fmt.Fprint(stdout, txt)
			return nil
		}
	}

	grokHome := strings.TrimSpace(opts.GrokHome)
	if grokHome == "" {
		grokHome = agenttty.GrokHome()
	}
	listOpts := opts.GrokListLiveOpts
	if listOpts == nil {
		listOpts = &sessions.ListLiveOpts{}
	}
	err := sessions.RunListLive(args, stdout, stderr, grokHome, listOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro grok session list-live", "kck grok list")
		return writeError(stderr, msg)
	}
	return nil
}

func runGrokOpen(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(grokOpenHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	grokHome := strings.TrimSpace(opts.GrokHome)
	if grokHome == "" {
		grokHome = agenttty.GrokHome()
	}

	openOpts := opts.GrokOpenOpts
	err := sessions.RunOpen(args, stdout, stderr, grokHome, openOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro grok session open", "kck grok open")
		return writeError(stderr, msg)
	}
	return nil
}

func runGrokFocus(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(grokFocusHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	grokHome := strings.TrimSpace(opts.GrokHome)
	if grokHome == "" {
		grokHome = agenttty.GrokHome()
	}

	focusOpts := opts.GrokFocusOpts
	err := sessions.RunFocus(args, stdout, grokHome, focusOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro grok session focus", "kck grok focus")
		return writeError(stderr, msg)
	}
	return nil
}

func runGrokSnapshot(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(grokSnapshotHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	grokHome := strings.TrimSpace(opts.GrokHome)
	if grokHome == "" {
		grokHome = agenttty.GrokHome()
	}

	snapOpts := opts.GrokSnapshotOpts
	err := sessions.RunSnapshot(args, stdout, stderr, grokHome, snapOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro grok session snapshot", "kck grok snapshot")
		return writeError(stderr, msg)
	}
	return nil
}

func runGrokSend(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(grokSendHelp, "\n")
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
			return runCronDryPreview(opts, expr, cronRaw, sendArgs)
		}
		return runCronSendLoop(opts, expr, cronRaw, sendArgs)
	}

	grokHome := strings.TrimSpace(opts.GrokHome)
	if grokHome == "" {
		grokHome = agenttty.GrokHome()
	}

	sendOpts := opts.GrokSendOpts
	err = sessions.RunSend(sendArgs, stdout, stderr, grokHome, sendOpts)
	if err != nil {
		return writeError(stderr, rewriteSendErr(err))
	}
	return nil
}



func runGrokInfo(opts Options, args []string) error {
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
			txt := strings.TrimPrefix(grokInfoHelp, "\n")
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

	grokHome := resolveGrokHome(opts)
	info, err := sessions.Info(grokHome, sessionID)
	if err != nil {
		return writeError(stderr, err.Error())
	}

	now := opts.GrokNow
	if now.IsZero() {
		now = time.Now()
	}
	home, _ := os.UserHomeDir()
	fmt.Fprintln(stdout, sessions.FormatInfoText(info, home, now))

	st, err := sessions.Status(grokHome, sessionID, !noPID, opts.GrokLiveOpts)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, sessions.FormatActiveBlock(st))
	return nil
}

func runGrokStatus(opts Options, args []string) error {
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
			txt := strings.TrimPrefix(grokStatusHelp, "\n")
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

	grokHome := resolveGrokHome(opts)
	st, err := sessions.Status(grokHome, sessionID, !noPID, opts.GrokLiveOpts)
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

func runGrokWait(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	var timeoutStr string
	remain, err := lessflags.String("--timeout", &timeoutStr).
		HelpFunc("-h,--help", func() {}).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if err == lessflags.ErrHelp {
			txt := strings.TrimPrefix(grokWaitHelp, "\n")
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

	var timeout time.Duration
	if strings.TrimSpace(timeoutStr) != "" {
		timeout, err = time.ParseDuration(timeoutStr)
		if err != nil {
			return writeError(stderr, fmt.Sprintf("invalid --timeout: %v", err))
		}
		if timeout <= 0 {
			return writeError(stderr, "--timeout must be > 0")
		}
	}

	waitOpts := sessions.WaitOpts{Timeout: timeout, Live: opts.GrokLiveOpts}
	if opts.GrokWaitOpts != nil {
		waitOpts = *opts.GrokWaitOpts
		if timeout > 0 {
			waitOpts.Timeout = timeout
		}
		if waitOpts.Live == nil {
			waitOpts.Live = opts.GrokLiveOpts
		}
	}

	res, err := sessions.Wait(resolveGrokHome(opts), sessionID, waitOpts)
	if err != nil {
		return writeError(stderr, err.Error())
	}
	fmt.Fprintf(stdout, "reason: %s\nsession-id: %s\n", res.Reason, res.SessionID)
	return nil
}

func runGrokResolve(opts Options, args []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	if argsHaveHelp(args) {
		txt := strings.TrimPrefix(grokResolveHelp, "\n")
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		fmt.Fprint(stdout, txt)
		return nil
	}

	var resolveOpts *sessions.ResolveOpts
	if opts.GrokResolveOpts != nil {
		cp := *opts.GrokResolveOpts
		resolveOpts = &cp
	} else {
		resolveOpts = &sessions.ResolveOpts{}
	}
	resolveOpts.Stdout = stdout
	resolveOpts.Stderr = stderr
	if strings.TrimSpace(resolveOpts.GrokHome) == "" {
		resolveOpts.GrokHome = resolveGrokHome(opts)
	}

	err := sessions.RunResolve(args, resolveOpts)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "agent-pro grok session resolve", "kck grok resolve")
		return writeError(stderr, msg)
	}
	return nil
}

func resolveGrokHome(opts Options) string {
	grokHome := strings.TrimSpace(opts.GrokHome)
	if grokHome == "" {
		return agenttty.GrokHome()
	}
	return grokHome
}

func argsHaveHelp(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}
