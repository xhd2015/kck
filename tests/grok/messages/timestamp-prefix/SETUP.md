# Scenario

```go
import "time"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, fixtureMultiTurnUpdates())
	req.Loc = time.FixedZone("CST", 8*3600)
	req.Args = []string{"grok", "messages", req.SessionID, "--limit", "2"}
	return nil
}
```
