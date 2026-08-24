# Scenario

**Feature**: single iTerm match formats `w=42 t=3`

```
TTY /dev/ttys042 matches one ref WindowID=42 TabIndex=3 -> ITERM w=42 t=3
```

## Steps

1. Parent session with TTY `/dev/ttys042`.
2. One matching ITermSession.

```go
import (
	"testing"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{}
	req.ITermSessions = []run.ITermSession{
		{WindowID: "42", TabIndex: 3, TTY: "/dev/ttys042"},
		{WindowID: "99", TabIndex: 1, TTY: "/dev/ttys999"},
	}
	return nil
}
```
