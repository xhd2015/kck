# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{"grok", "new", "--here", "--new-terminal", "x"}
	return nil
}
```
