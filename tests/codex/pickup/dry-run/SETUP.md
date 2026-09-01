# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{
		"codex", "pickup", "extract TODOs",
		"--session-id", req.SessionID,
		"--dry-run",
	}
	return nil
}
```
