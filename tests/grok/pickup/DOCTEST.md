# `kck grok pickup`

Open a new empty Grok session and stage a kck-pickup-a-session draft from a base
session. Default: current terminal; `--new-window` for ForceNew. Draft skill path
uses tilde form.

L2 only — injectable `GrokHome` + `GrokPickupOpts`. No live iTerm.

## Version

0.0.3

## Decision Tree

```text
grok/pickup/
├── help/
│   ├── root-lists-pickup/
│   ├── grok-usage/
│   └── pickup-usage/
├── missing-message/
├── missing-session-source/
├── here-new-window-conflict/
├── dry-run/
├── agent-run-here/
├── agent-run-new-window/
└── no-agent-run-new-window/
```

## How to Run

```sh
doctest vet ./tests/grok/pickup
doctest test ./tests/grok/pickup
```

```go
import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"

	"kck/run"
)

type Request struct {
	Args       []string
	GrokHome   string
	TempDir    string
	ProjectDir string
	SessionID  string
	CachePath  string
	SkillBody  string
	UserHome   string

	LookPathMap map[string]string
}

type Response struct {
	Stdout     string
	Stderr     string
	ErrText    string
	ExitCode   int
	Opened     []string
	Foreground []string
	BareHere   []string
	Staged     []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	var opened, foreground, bareHere, staged []string
	look := req.LookPathMap
	if look == nil {
		look = map[string]string{
			"agent-run": "/usr/local/bin/agent-run",
			"grok":      "/usr/local/bin/grok",
		}
	}
	body := req.SkillBody
	if body == "" {
		body = "# kck-pickup-a-session fixture\n"
	}
	home := req.UserHome
	popts := &run.PickupOpts{
		SkillContent:   body,
		CacheSkillPath: req.CachePath,
		UserHomeDir:    func() (string, error) { return home, nil },
		TildeHome: func(path string) string {
			if home != "" && strings.HasPrefix(path, home) {
				return "~" + strings.TrimPrefix(path, home)
			}
			return path
		},
		ReadFile:  os.ReadFile,
		WriteFile: os.WriteFile,
		MkdirAll:  os.MkdirAll,
		Rename:    os.Rename,
		LookPath: func(file string) (string, error) {
			if p, ok := look[file]; ok {
				return p, nil
			}
			return "", os.ErrNotExist
		},
		OpenInNewWindow: func(dir, followUp string) error {
			opened = append(opened, dir+"|"+followUp)
			return nil
		},
		RunForeground: func(bin string, argv []string, dir string) error {
			foreground = append(foreground, dir+"|"+bin+" "+strings.Join(argv, " "))
			return nil
		},
		RunBareHere: func(bin, dir, draft string) error {
			bareHere = append(bareHere, dir+"|"+bin+"|"+draft)
			staged = append(staged, draft)
			return nil
		},
		StageDraft: func(draft string) error {
			staged = append(staged, draft)
			return nil
		},
		Sleep: func(time.Duration) error { return nil },
		ResolveBaseCWD: func(h, sessionID string) (string, error) {
			_ = h
			if sessionID == req.SessionID {
				return req.ProjectDir, nil
			}
			return "", os.ErrNotExist
		},
	}

	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:           append([]string(nil), req.Args...),
		Stdout:         &stdout,
		Stderr:         &stderr,
		GrokHome:       req.GrokHome,
		GrokPickupOpts: popts,
	})
	resp := &Response{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		Opened:     append([]string(nil), opened...),
		Foreground: append([]string(nil), foreground...),
		BareHere:   append([]string(nil), bareHere...),
		Staged:     append([]string(nil), staged...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
