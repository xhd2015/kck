# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	req.Args = []string{"codex", "send", "hi", "--session-id", req.SessionID, "--cron", "not-a-cron"}
	return nil
}
```
