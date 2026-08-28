# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckMessagesSessionID
	lines := []string{
		`{"timestamp":"2026-08-01T12:00:01.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"m0"}]}}`,
		`{"timestamp":"2026-08-01T12:00:02.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"m1"}]}}`,
		`{"timestamp":"2026-08-01T12:00:03.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"m2"}]}}`,
	}
	writeKckMessagesSession(t, req, lines...)
	req.Args = []string{"codex", "messages", req.SessionID, "--limit", "2"}
	return nil
}
```
