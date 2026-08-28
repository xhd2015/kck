# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckMessagesSessionID
	writeKckMessagesSession(t, req, fixtureKindsLines()...)
	req.Args = []string{"codex", "messages", req.SessionID, "--limit", "10"}
	return nil
}
```
