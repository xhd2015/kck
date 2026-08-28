# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"codex", "status"}
	return nil
}
```
