# Scenario

**Feature**: `kck skill --list` prints skill name and topic paths

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"skill", "--list"}
	return nil
}
```
