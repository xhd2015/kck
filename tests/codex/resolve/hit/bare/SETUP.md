# Scenario

**Feature**: default bare session id on stdout

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	seedHit(req, fixtureSessionID, pidCodex)
	req.Args = []string{"codex", "resolve"}
	return nil
}
```
