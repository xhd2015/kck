# Scenario

**Feature**: `--no-iterm` skips resolution; ITERM is `-`

```
matching refs exist but --no-iterm -> ITERM - for all
```

## Steps

1. Matching iTerm ref present.
2. Args include `--no-iterm`.

```go
import (
	"testing"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--no-iterm"}
	req.ITermSessions = []run.ITermSession{
		{WindowID: "42", TabIndex: 3, TTY: "/dev/ttys042"},
	}
	return nil
}
```
