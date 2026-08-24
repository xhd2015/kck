# Scenario

**Feature**: no iTerm match shows `-`

```
TTY /dev/ttys042 with only unrelated refs -> ITERM -
```

## Steps

1. Parent session TTY ttys042.
2. ITerm list has no matching TTY.

```go
import (
	"testing"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{}
	req.ITermSessions = []run.ITermSession{
		{WindowID: "9", TabIndex: 1, TTY: "/dev/ttys999"},
	}
	return nil
}
```
