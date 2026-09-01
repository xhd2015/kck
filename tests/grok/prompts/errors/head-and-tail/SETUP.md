# Scenario

**Feature**: `--head` and `--tail` are mutually exclusive

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "prompts", "--head", "2", "--tail", "2"}
	return nil
}
```
