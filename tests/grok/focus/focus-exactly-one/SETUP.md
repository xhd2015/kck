# Scenario

**Feature**: one live host tab → focused

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckFocusSessionID
	writeKckFocusSession(t, req)
	addLiveGrok(req, 5001, "ttys148")
	req.ITerm = oneITermTab()
	req.Args = []string{"grok", "focus", req.SessionID}
	return nil
}
```
