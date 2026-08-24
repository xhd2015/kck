# Scenario

**Feature**: plain idle shell without agent is omitted (agents-only)

```
idle zsh only, no Agents -> 0 sessions footer; no iterm-uuid-plain-zsh
```

## Steps

1. Inject idleZshOnlyResult().

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = idleZshOnlyResult()
	return nil
}
```
