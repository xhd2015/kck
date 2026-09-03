# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckWaitSessionID
	writeKckWaitSession(t, req, []string{
		`{"sessionUpdate":"turn_completed","prompt_id":"p0","stop_reason":"end_turn"}`,
	})
	// no live PIDs → inactive
	req.Args = []string{"grok", "wait", req.SessionID}
	return nil
}
```
