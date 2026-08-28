# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"codex", "status", "019f283a-ffff-7fff-ffff-ffffffffff99"}
	return nil
}
```
