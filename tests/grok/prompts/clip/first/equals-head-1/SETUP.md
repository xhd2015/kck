# Scenario

**Feature**: `--first` is sugar for `--head 1`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "prompts", req.SessionID, "--first", "--no-color"}
	return nil
}
```
