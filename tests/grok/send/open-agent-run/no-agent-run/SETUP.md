# Scenario

```go
import "github.com/xhd2015/agent-pro/agent/grok/sessions"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	req.AfterOpenHost = true
	req.NoAgentRun = true
	req.AgentRunByID = map[string]*sessions.AgentRunOpenResult{
		req.SessionID: {
			AgentRunSessionID: "ar-skip",
			Mode:              sessions.AgentRunOpenModeSend,
			Delivered:         true,
		},
	}
	req.Args = []string{"grok", "send", "hello", "--session-id", req.SessionID, "--open", "--no-agent-run"}
	return nil
}
```
