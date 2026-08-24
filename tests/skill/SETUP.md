# Scenario

**Feature**: kck multi-topic skill surface via `kck skill`

```
# skillcmd SingleSkill + embedded docs TreeFS
Caller -> run.MainWith(["skill", …]) -> stdout/stderr/exit
# nested path/TOPIC.md topics under docs/; root docs/SKILL.md
```

## Preconditions

- Skill content is embedded (`docs/` + `skillcmd.SingleSkill`); no home/iTerm
  fixture required.

## Steps

1. Descendant setups set `req.Args` for the skill action under test.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	_ = req
	return nil
}

func assertNoHarnessErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
}

func assertSuccess(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode != 0 || resp.ErrText != "" {
		t.Fatalf("want success exit 0, got exit=%d err=%q stderr=%q stdout=%q",
			resp.ExitCode, resp.ErrText, resp.Stderr, resp.Stdout)
	}
}

func assertFailure(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode == 0 && resp.ErrText == "" {
		t.Fatalf("want non-zero failure, got exit=0 empty err; stdout=%q stderr=%q",
			resp.Stdout, resp.Stderr)
	}
}

func assertContains(t *testing.T, got, want, label string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s must contain %q; got:\n%s", label, want, got)
	}
}
```
