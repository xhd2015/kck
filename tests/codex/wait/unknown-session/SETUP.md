# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"codex", "wait", "00000000-0000-0000-0000-000000000000"}
	return nil
}
```
