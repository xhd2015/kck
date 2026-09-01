# Scenario

**Feature**: `--first` cannot combine with `--head`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "prompts", "--first", "--head", "2"}
	return nil
}
```
