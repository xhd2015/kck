# Scenario

**Feature**: `--tab 2` (1-based) resolves sibling tab codex session

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	seedTabWindow(req)
	req.Args = []string{"codex", "resolve", "--tab", "2"}
	return nil
}
```
