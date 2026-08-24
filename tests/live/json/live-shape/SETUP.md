# Scenario

**Feature**: live JSON envelope matches store shape (sessions + summary)

```
busyGrokAgentResult + --json -> valid JSON, sessions len 1, no ANSI
```

## Steps

1. Inject busyGrokAgentResult.
2. Args = `--json`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{"--json"}
	req.LiveResult = busyGrokAgentResult()
	return nil
}
```
