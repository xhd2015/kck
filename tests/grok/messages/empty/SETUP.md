# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, "")
	req.Args = []string{"grok", "messages", req.SessionID, "--limit", "5"}
	return nil
}
```
