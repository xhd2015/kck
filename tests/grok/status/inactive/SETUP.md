# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckStatusSessionID
	writeKckStatusSession(t, req)
	req.Args = []string{"grok", "status", req.SessionID}
	return nil
}
```
