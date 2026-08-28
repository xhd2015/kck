package run

import (
	"fmt"
	"io"
	"strings"
	"time"

	codexsessions "github.com/xhd2015/agent-pro/agent/codex/sessions"
	groksessions "github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/agent-pro/pkgs/agenttty"
	"github.com/xhd2015/dot-pkgs/go-pkgs/cron/easycron"
)

const cronDryPreviewCount = 3

// peelCronFlag extracts --cron / --cron=EXPR from args. Remain keeps all other
// tokens for sessions.RunSend. Duplicate --cron is an error.
func peelCronFlag(args []string) (cron string, remain []string, err error) {
	remain = make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--cron":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--cron requires an expression")
			}
			if cron != "" {
				return "", nil, fmt.Errorf("duplicate --cron")
			}
			cron = args[i+1]
			i++
		case strings.HasPrefix(a, "--cron="):
			if cron != "" {
				return "", nil, fmt.Errorf("duplicate --cron")
			}
			cron = strings.TrimPrefix(a, "--cron=")
		default:
			remain = append(remain, a)
		}
	}
	return cron, remain, nil
}

func argsHaveDryRun(args []string) bool {
	for _, a := range args {
		if a == "--dry-run" {
			return true
		}
	}
	return false
}

func cronNow(opts Options) func() time.Time {
	if opts.GrokCronNow != nil {
		return opts.GrokCronNow
	}
	return time.Now
}

func cronSleep(opts Options) func(time.Duration) error {
	if opts.GrokCronSleep != nil {
		return opts.GrokCronSleep
	}
	return func(d time.Duration) error {
		time.Sleep(d)
		return nil
	}
}

func cronLoc(opts Options) *time.Location {
	if opts.GrokCronLoc != nil {
		return opts.GrokCronLoc
	}
	return time.Local
}

func formatCronParseErr(err error) string {
	msg := strings.TrimPrefix(err.Error(), "easycron: ")
	return "invalid --cron: " + msg
}

func rewriteSendErr(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	msg = strings.ReplaceAll(msg, "agent-pro grok session send", "kck grok send")
	msg = strings.ReplaceAll(msg, "agent-pro codex session send", "kck codex send")
	msg = strings.ReplaceAll(msg, "agent-pro grok session focus", "kck grok send")
	msg = strings.ReplaceAll(msg, "agent-pro codex session focus", "kck codex send")
	return msg
}

type cronRunner string

const (
	cronRunnerGrok  cronRunner = "grok"
	cronRunnerCodex cronRunner = "codex"
)

func cronRunSend(opts Options, runner cronRunner, sendArgs []string, stdout, stderr io.Writer) error {
	switch runner {
	case cronRunnerCodex:
		return codexsessions.RunSend(sendArgs, stdout, stderr, resolveCodexHome(opts), opts.CodexSendOpts)
	default:
		home := strings.TrimSpace(opts.GrokHome)
		if home == "" {
			home = agenttty.GrokHome()
		}
		return groksessions.RunSend(sendArgs, stdout, stderr, home, opts.GrokSendOpts)
	}
}

func runCronDryPreview(opts Options, expr easycron.Expr, raw string, sendArgs []string) error {
	return runCronDryPreviewFor(opts, cronRunnerGrok, expr, raw, sendArgs)
}

func runCronSendLoop(opts Options, expr easycron.Expr, raw string, sendArgs []string) error {
	return runCronSendLoopFor(opts, cronRunnerGrok, expr, raw, sendArgs)
}

func runCronDryPreviewFor(opts Options, runner cronRunner, expr easycron.Expr, raw string, sendArgs []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	nowFn := cronNow(opts)
	loc := cronLoc(opts)
	anchor := nowFn()

	fmt.Fprintf(stdout, "cron %s\n", raw)
	from := anchor
	for i := 0; i < cronDryPreviewCount; i++ {
		next, ok := expr.Next(anchor, from, loc)
		if !ok {
			break
		}
		fmt.Fprintf(stdout, "  next[%d]  %s\n", i, next.In(loc).Format(time.RFC3339))
		from = next.Add(time.Nanosecond)
	}

	err := cronRunSend(opts, runner, sendArgs, stdout, stderr)
	if err != nil {
		return writeError(stderr, rewriteSendErr(err))
	}
	return nil
}

func runCronSendLoopFor(opts Options, runner cronRunner, expr easycron.Expr, raw string, sendArgs []string) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	nowFn := cronNow(opts)
	sleepFn := cronSleep(opts)
	loc := cronLoc(opts)
	anchor := nowFn()
	from := anchor
	tick := 0

	for {
		next, ok := expr.Next(anchor, from, loc)
		if !ok {
			fmt.Fprintln(stdout, "cron done: until reached")
			return nil
		}

		if wait := next.Sub(nowFn()); wait > 0 {
			if err := sleepFn(wait); err != nil {
				return writeError(stderr, err.Error())
			}
		}

		err := cronRunSend(opts, runner, sendArgs, stdout, stderr)
		tick++
		if err != nil {
			msg := rewriteSendErr(err)
			if tick == 1 {
				return writeError(stderr, msg)
			}
			fmt.Fprintf(stderr, "warning: send failed: %s; will retry\n", msg)
		}

		upcoming, ok := expr.Next(anchor, next.Add(time.Nanosecond), loc)
		if !ok {
			fmt.Fprintln(stdout, "cron done: until reached")
			return nil
		}
		fmt.Fprintf(stdout, "next %s (%s)\n", upcoming.In(loc).Format("2006-01-02 15:04:05"), raw)
		from = next.Add(time.Nanosecond)
	}
}
