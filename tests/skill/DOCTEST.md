# `kck skill`

Shape 3 multi-topic skill surface via `skillcmd.SingleSkill` + embedded
`docs/` tree. L2 in-process — `run.MainWith` with injectable writers.
skillcmd writes to `os.Stdout`; `run` captures under a mutex into `Stdout`.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith(Options{Args: ["skill",…], Stdout, Stderr})`.
- **`kck/run`** — dispatches `skill`, prefixes errors with `Error: kck skill:`.
- **`docs`** — embedded `SKILL.md` + `path/TOPIC.md` TreeFS.
- **`skillcmd.SingleSkill`** — `--show` / `--install` / `--list` / `-h`.

## Behaviors

- Root help mentions `skill`.
- `kck skill -h` documents `--show` / `--install` / `--list` and lists topics.
- `kck skill --list` prints `kck` then sorted topic paths.
- `kck skill --show` prints root index (`name: kck`); no install plumbing flags.
- `kck skill --show send` prints `name: kck/send`.
- Unknown topic → non-zero; `Error: kck skill:` + unknown/missing wording.

## Version

0.0.1

## Decision Tree

```text
skill/
├── help/
│   ├── root-lists-skill/
│   └── skill-usage/
├── list/
├── show/
│   ├── root/
│   ├── topic/
│   └── unknown/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-skill/` | `kck -h` mentions `skill`. |
| `help/skill-usage/` | `kck skill -h` usage + Available topics. |
| `list/` | skill name + topic paths. |
| `show/root/` | `name: kck`; retrieve example; no `--cursor`/`--global`. |
| `show/topic/` | `name: kck/send`. |
| `show/unknown/` | `Error: kck skill:` + unknown/missing topic. |

## How to Run

```sh
doctest vet ./tests/skill
doctest test ./tests/skill
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/doctest/session"

	"kck/run"
)

type Request struct {
	Args []string
}

type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:   append([]string(nil), req.Args...),
		Stdout: &stdout,
		Stderr: &stderr,
	})
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
