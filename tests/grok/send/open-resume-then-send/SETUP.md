# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	req.AfterOpenHost = true
	req.Args = []string{"grok", "send", "hello", "--session-id", req.SessionID, "--open"}
	return nil
}
```
