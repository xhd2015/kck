# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckMessagesSessionID
	writeKckMessagesSession(t, req)
	req.Args = []string{"codex", "messages", req.SessionID}
	return nil
}
```
