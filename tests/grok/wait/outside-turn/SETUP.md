# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckWaitSessionID
	writeKckWaitSession(t, req, []string{
		`{"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"hi"}}`,
		`{"sessionUpdate":"turn_completed","prompt_id":"p1","stop_reason":"end_turn"}`,
	})
	markRunning(t, req)
	req.Args = []string{"grok", "wait", req.SessionID}
	return nil
}
```
