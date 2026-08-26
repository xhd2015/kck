# Scenario

**Feature**: `--no-color` suppresses ANSI even with `--grep`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, fixtureMultiTurnUpdates())
	req.Args = []string{"grok", "messages", req.SessionID, "--grep", "a1", "--no-color", "--limit", "0"}
	return nil
}
```
