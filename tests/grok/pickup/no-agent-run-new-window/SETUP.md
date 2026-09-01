# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{
		"grok", "pickup", "continue the refactor",
		"--session-id", req.SessionID,
		"--no-agent-run",
		"--new-window",
	}
	return nil
}
```
