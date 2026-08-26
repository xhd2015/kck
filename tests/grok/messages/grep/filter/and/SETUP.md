# Scenario

**Feature**: multiple `--grep` require every pattern in the same message (AND)

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, fixtureMultiTurnUpdates())
	req.Args = []string{"grok", "messages", req.SessionID, "--grep", "run_terminal", "--grep", "echo", "--limit", "0"}
	return nil
}
```
