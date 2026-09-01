# Scenario

**Feature**: `--tail 2` keeps last two prompts and marks leading omissions

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"grok", "prompts", req.SessionID, "--tail", "2", "--no-color"}
	return nil
}
```
