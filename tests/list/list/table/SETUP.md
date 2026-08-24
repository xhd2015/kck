# Scenario

**Feature**: multi-session human table with injected probe

```
seeded metas + Probe map -> human rows newest-first; columns; D1 REASON
```

## Steps

1. Leaves seed 2+ sessions and probe fields.
2. Human format (no `--json`).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	return nil
}
```
