# Scenario

**Feature**: known session but no live host → not found (no resume)

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckFocusSessionID
	writeKckFocusSession(t, req)
	req.Args = []string{"grok", "focus", req.SessionID}
	return nil
}
```
