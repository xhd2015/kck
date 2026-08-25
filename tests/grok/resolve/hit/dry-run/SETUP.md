# Scenario

**Feature**: ancestor dry-run prints plan lines

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	seedHit(req, fixtureSessionID, pidGrok)
	req.Args = []string{"grok", "resolve", "--dry-run"}
	return nil
}
```
