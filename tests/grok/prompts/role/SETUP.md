# Scenario

```go
import "time"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	base := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	writePromptSession(t, req, promptSessionOpts{
		ID:           "019f283a-prmt-7ccc-cccc-cccccccccc11",
		Title:        "main agent",
		Kind:         "main",
		LastActiveAt: base,
		Updates:      updatesJSONL(userChunkAt("main-prompt", base)),
	})
	writePromptSession(t, req, promptSessionOpts{
		ID:           "019f283a-prmt-7ccc-cccc-cccccccccc22",
		Title:        "sub agent",
		Kind:         "subagent",
		ParentID:     "019f283a-prmt-7ccc-cccc-cccccccccc11",
		LastActiveAt: base.Add(-time.Minute),
		Updates:      updatesJSONL(userChunkAt("sub-prompt", base.Add(-time.Minute))),
	})
	return nil
}
```
