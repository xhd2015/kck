# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckMessagesSessionID
	lines := []string{
		`{"timestamp":"2026-08-01T12:00:01.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"keep zebra"}]}}`,
		`{"timestamp":"2026-08-01T12:00:02.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"drop"}]}}`,
	}
	writeKckMessagesSession(t, req, lines...)
	req.Args = []string{"codex", "messages", req.SessionID, "--grep", "zebra", "--limit", "10"}
	return nil
}
```
