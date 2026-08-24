# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{"grok", "send", "hi", "--session-id", "019f283a-dead-7ead-dead-deaddeaddead"}
	return nil
}
```
