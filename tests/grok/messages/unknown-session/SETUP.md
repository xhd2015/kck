# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "messages", "deadbeef-not-a-session", "--limit", "2"}
	return nil
}
```
