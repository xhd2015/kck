# Scenario

**Feature**: empty home prints zero-session summary

```
MainWith([]) on empty home -> footer "0 sessions · 0 needs attention · 0 sendable"
```

## Steps

1. Args empty (default list).
2. No sessions.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{}
	req.Sessions = nil
	return nil
}
```
