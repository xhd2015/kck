# Scenario

**Feature**: `--json` machine-readable list without ANSI

```
MainWith(["--json"], fixtures) -> JSON object; no ESC sequences
```

## Steps

1. Leaves set sessions + `--json`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	twoSessionFixture(req)
	return nil
}
```
