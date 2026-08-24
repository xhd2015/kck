# Scenario

**Feature**: ListITerm error soft-fails with warning and `-`

```
ListITerm returns error -> warning: on stderr; ITERM -; exit 0
```

## Steps

1. Set ListITermErr.
2. Matching would have been possible if list succeeded (optional).

```go
import (
	"testing"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{}
	req.ListITermErr = "osascript failed: permission denied"
	req.ITermSessions = []run.ITermSession{
		{WindowID: "42", TabIndex: 3, TTY: "/dev/ttys042"},
	}
	return nil
}
```
