# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "focus", "019f283a-dead-7ead-dead-deaddeaddead"}
	return nil
}
```
