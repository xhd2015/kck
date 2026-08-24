# Scenario

**Feature**: live `--json` machine-readable list without ANSI

```
LiveCapture Result + --json -> {sessions, summary} JSON
```

## Steps

1. Live path.
2. Child sets --json and LiveResult.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Home = ""
	return nil
}
```
