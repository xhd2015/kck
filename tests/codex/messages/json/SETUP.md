# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckMessagesSessionID
	writeKckMessagesSession(t, req, fixtureKindsLines()[0])
	req.Args = []string{"codex", "messages", req.SessionID, "--json", "--limit", "5"}
	return nil
}
```
