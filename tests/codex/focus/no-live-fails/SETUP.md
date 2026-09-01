# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckFocusSessionID
	writeKckFocusSession(t, req)
	req.Args = []string{"codex", "focus", req.SessionID}
	return nil
}
```
