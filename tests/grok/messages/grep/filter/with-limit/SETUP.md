# Scenario

**Feature**: `--limit` applies after `--grep` (newest matches)

```go
import "strings"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	updates := strings.Join([]string{
		`{"timestamp":1785493071,"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"hit-0"}}`,
		`{"timestamp":1785493072,"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"miss"}}`,
		`{"timestamp":1785493072,"sessionUpdate":"turn_completed"}`,
		`{"timestamp":1785493073,"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"hit-1"}}`,
		`{"timestamp":1785493074,"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hit-2"}}`,
		`{"timestamp":1785493074,"sessionUpdate":"turn_completed"}`,
	}, "\n")
	writeKckMessagesSession(t, req, updates)
	req.Args = []string{"grok", "messages", req.SessionID, "--grep", "hit", "--limit", "2"}
	return nil
}
```
