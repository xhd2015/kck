# Scenario

# clip — per-session head/tail/first

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	writePromptSession(t, req, promptSessionOpts{
		ID:      fixtureKckPromptsSessionID,
		Updates: fiveUserPrompts(),
	})
	return nil
}
```
