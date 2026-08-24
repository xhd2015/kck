package run

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/xhd2015/skills/skillcmd"

	"kck/docs"
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
		return singleSkill().Handle(args)
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
       kck skill --install [OPTIONS] [<dir>]
       kck skill --list

Show the root SKILL.md index or a nested topic (path/TOPIC.md).
Install copies SKILL.md and nested TOPIC.md topics into agent skill directories.
List prints the skill name and every available topic path.
--help also lists available topics (see below).

Examples:
  kck skill --show
  kck skill --show send
  kck skill send --show
  kck skill --list
  kck skill --install --dry-run
  kck skill --install --help

Options:
  --show [--header] [path]   Print skill or topic content (header-only with --header)
  --install [OPTIONS] [dir]  Install skill files (see --install --help)
  --list                     Print skill name and all topic paths
  -h, --help                 Show this help and available topics
`
}
