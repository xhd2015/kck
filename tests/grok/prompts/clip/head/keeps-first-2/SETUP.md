# Scenario

**Feature**: `--head 2` keeps first two prompts and marks trailing omissions

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "prompts", req.SessionID, "--head", "2", "--no-color"}
	return nil
}
```
