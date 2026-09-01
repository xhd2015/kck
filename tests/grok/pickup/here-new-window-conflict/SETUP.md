# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{
		"grok", "pickup", "hi",
		"--session-id", req.SessionID,
		"--here",
		"--new-window",
	}
	return nil
}
```
