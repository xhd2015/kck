# Scenario

**Feature**: `--main` multi mode keeps main-agent session only

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "prompts", "--main", "--limit", "10", "--no-color"}
	return nil
}
```
