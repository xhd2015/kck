# Scenario

**Feature**: zero grep hits prints `(no matching messages)`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, fixtureMultiTurnUpdates())
	req.Args = []string{"grok", "messages", req.SessionID, "--grep", "zzznomatch"}
	return nil
}
```
