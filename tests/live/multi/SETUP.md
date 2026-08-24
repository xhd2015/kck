# Scenario

**Feature**: multi-agent live list and footer aggregation

```
two Agents panes -> both session ids + footer N/M/K
```

## Steps

1. Live path.
2. Child injects twoAgentsResult.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Home = ""
	return nil
}
```
