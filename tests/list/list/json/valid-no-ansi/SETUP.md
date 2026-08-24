# Scenario

**Feature**: JSON list is valid and ANSI-free

```
--json + two sessions -> parseable JSON; no \x1b; summary counts
```

## Steps

1. Parent twoSessionFixture.
2. Args = `--json`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--json"}
	return nil
}
```
