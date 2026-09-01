package run

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/xhd2015/skills/skillcmd"

	"kck/docs"
	pickupskill "kck/skills/kck-pickup-a-session"
)

// skillcmd prints via fmt to os.Stdout; serialize capture for L2 injectables.
var skillStdoutMu sync.Mutex

func singleSkill() *skillcmd.SingleSkill {
	return &skillcmd.SingleSkill{
		Name:        docs.Name,
		RootContent: docs.SkillMD,
		TreeFS:      docs.TreeFS,
		Usage:       "kck skill --install",
		Help:        skillHelp(),
	}
}

func runSkill(opts Options) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	args := opts.Args[1:]
	out, err := captureStdout(func() error {
		return handleSkillArgs(args)
	})
	if err != nil {
		return writeError(stderr, fmt.Sprintf("kck skill: %v", err))
	}
	if out != "" {
		if _, werr := io.WriteString(stdout, out); werr != nil {
			return writeError(stderr, fmt.Sprintf("kck skill: write stdout: %v", werr))
		}
	}
	return nil
}

// handleSkillArgs keeps Shape 3 SingleSkill for the kck CLI topic tree, and
// exposes kck-pickup-a-session as a sibling agent skill (show/install/list).
func handleSkillArgs(args []string) error {
	parsed, err := skillcmd.ParseSkillArgs(args)
	if err != nil {
		return err
	}

	switch parsed.Action {
	case skillcmd.ActionHelp:
		if err := singleSkill().Handle([]string{"-h"}); err != nil {
			return err
		}
		fmt.Printf("\nAgent skills:\n  %s\n", pickupskill.Name)
		return nil
	case skillcmd.ActionList:
		if err := singleSkill().Handle([]string{"--list"}); err != nil {
			return err
		}
		fmt.Println(pickupskill.Name)
		return nil
	case skillcmd.ActionShow, skillcmd.ActionVersion, skillcmd.ActionInstall:
		if len(parsed.Rest) > 0 && parsed.Rest[0] == pickupskill.Name {
			return handlePickupAgentSkill(parsed)
		}
		// Fall through: kck root / topics / install tree.
		return singleSkill().Handle(args)
	default:
		return singleSkill().Handle(args)
	}
}

func handlePickupAgentSkill(parsed skillcmd.ParsedArgs) error {
	rest := parsed.Rest[1:] // drop skill name
	switch parsed.Action {
	case skillcmd.ActionShow:
		if len(rest) > 0 {
			return fmt.Errorf("unexpected arguments: %v", rest)
		}
		if parsed.Header {
			out, err := skillcmd.FormatHeaderWithDelimiters(pickupskill.SkillMD)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		}
		fmt.Print(pickupskill.SkillMD)
		return nil
	case skillcmd.ActionVersion:
		if len(rest) > 0 {
			return fmt.Errorf("unexpected arguments: %v", rest)
		}
		// skillcmd.printSkillVersion is unexported; reuse SingleSkill.handleVersion path.
		tmp := &skillcmd.SingleSkill{Name: pickupskill.Name, RootContent: pickupskill.SkillMD}
		return tmp.Handle([]string{"--version"})
	case skillcmd.ActionInstall:
		return skillcmd.HandleInstall(skillcmd.InstallOptions{
			SkillDirName: pickupskill.Name,
			SkillContent: pickupskill.SkillMD,
			Usage:        "kck skill --install " + pickupskill.Name,
		}, rest)
	default:
		return fmt.Errorf("unsupported action for %s", pickupskill.Name)
	}
}

func captureStdout(fn func() error) (string, error) {
	skillStdoutMu.Lock()
	defer skillStdoutMu.Unlock()

	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	old := os.Stdout
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fnErr := fn()
	_ = w.Close()
	os.Stdout = old
	out := <-done
	_ = r.Close()
	return out, fnErr
}

func skillHelp() string {
	return `Usage: kck skill --show [--header] [<topic-path>]
       kck skill <topic-path> --show [--header]
       kck skill --show kck-pickup-a-session
       kck skill kck-pickup-a-session --show
       kck skill --install [OPTIONS] [<dir>]
       kck skill --install kck-pickup-a-session [OPTIONS] [<dir>]
       kck skill --list

Show the root SKILL.md index or a nested topic (path/TOPIC.md).
Also hosts the agent skill kck-pickup-a-session (used by grok/codex pickup).
Install copies SKILL.md and nested TOPIC.md topics into agent skill directories,
or installs kck-pickup-a-session when that name is given.
List prints the skill name, every topic path, then kck-pickup-a-session.
--help also lists available topics and agent skills.

Examples:
  kck skill --show
  kck skill --show send
  kck skill send --show
  kck skill --show kck-pickup-a-session
  kck skill --list
  kck skill --install --dry-run
  kck skill --install --help

Options:
  --show [--header] [path]   Print skill or topic content (header-only with --header)
  --install [OPTIONS] [dir]  Install skill files (see --install --help)
  --list                     Print skill name, topics, and agent skills
  -h, --help                 Show this help and available topics
`
}
