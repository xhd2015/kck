# Scenario

**Feature**: `--color` bold-red highlights grep match spans

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, fixtureMultiTurnUpdates())
	req.Args = []string{"grok", "messages", req.SessionID, "--grep", "a1", "--color", "--limit", "0"}
	return nil
}
```
