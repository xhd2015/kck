# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{"skill", "--show", "kck-pickup-a-session"}
	return nil
}
```
