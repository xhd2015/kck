# Scenario

**Feature**: `--json` + `--grep` returns filtered page; no ANSI; total = match count

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writeKckMessagesSession(t, req, fixtureMultiTurnUpdates())
	req.Args = []string{"grok", "messages", req.SessionID, "--grep", "a1", "--json", "--limit", "0"}
	return nil
}
```
