# Scenario

```go
import "github.com/xhd2015/agent-pro/agent/grok/sessions"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	addLiveGrok(req, 4242, "/dev/ttys148")
	req.ITerm = oneITermTab()
	req.AgentRunByID = map[string]*sessions.AgentRunOpenResult{
		req.SessionID: {
			AgentRunSessionID: "ar-kck-live",
			Mode:              sessions.AgentRunOpenModeSend,
			Delivered:         true,
		},
	}
	req.Args = []string{"grok", "send", "hello", "--session-id", req.SessionID, "--open"}
	return nil
}
```
