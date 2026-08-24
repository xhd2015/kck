# Scenario

**Feature**: empty agent-run home lists zero sessions

```
empty Home/sessions -> header + footer 0 counts; exit 0
```

## Steps

1. No session seeds.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Sessions = nil
	return nil
}
```
