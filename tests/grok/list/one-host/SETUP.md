# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	addKckListHost(req, 5001, "ttys148", fixtureKckListSID, "3", 1)
	req.Args = []string{"grok", "list"}
	return nil
}
```
