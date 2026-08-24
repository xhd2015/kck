# Scenario

**Feature**: `--limit` caps live rows after inclusion/filter

```
two agents + --limit 1 -> exactly one agent session id in data rows
```

## Steps

1. Live path.
2. Child sets limit args + twoAgentsResult.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Home = ""
	return nil
}
```
