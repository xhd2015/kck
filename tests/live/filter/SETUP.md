# Scenario

**Feature**: live list flags --needs-confirm and --sendable filter included rows

```
twoAgentsResult + --needs-confirm -> only busy agent-sess-a
twoAgentsResult + --sendable -> only idle agent-sess-b
```

## Steps

1. Shared fixture twoAgentsResult for filter children.
2. Child sets filter Args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Home = ""
	req.LiveResult = twoAgentsResult()
	return nil
}
```
