# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	addLiveGrok(req, 5001, "ttys148")
	req.ITerm = oneITermTab()
	req.Args = []string{"grok", "send", "partial", "--session-id", req.SessionID,
		"--no-submit", "--focus", "--no-ctrl-u"}
	return nil
}
```
