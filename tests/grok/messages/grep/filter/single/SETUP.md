# Scenario

**Feature**: single `--grep` keeps matching message bodies only

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, fixtureMultiTurnUpdates())
	req.Args = []string{"grok", "messages", req.SessionID, "--grep", "a1", "--limit", "0"}
	return nil
}
```
