# Scenario

**Feature**: `--tail 0` is rejected

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "prompts", "--tail", "0"}
	return nil
}
```
