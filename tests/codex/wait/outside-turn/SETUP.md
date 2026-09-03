# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckCodexWaitSessionID
	writeKckCodexWaitSession(t, req, []string{
		eventMsgLine("task_started"),
		eventMsgLine("task_complete"),
	})
	markRunning(t, req)
	req.Args = []string{"codex", "wait", req.SessionID}
	return nil
}
```
