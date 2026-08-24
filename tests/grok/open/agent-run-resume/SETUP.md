# Scenario

```go
import "github.com/xhd2015/agent-pro/agent/grok/sessions"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckOpenSessionID
	writeKckOpenSession(t, req)
	req.Procs = nil
	req.OpenFiles = map[int][]string{}
	req.ITerm = nil
	req.AgentRunByID = map[string]*sessions.AgentRunOpenResult{
		req.SessionID: {
			AgentRunSessionID: "ar-kck-open",
			Mode:              sessions.AgentRunOpenModeResume,
			Opened:            true,
			CWD:               req.ProjectDir,
			Command:           "/usr/local/bin/agent-run run --session-id ar-kck-open --auto-send-or-resume --open",
		},
	}
	req.Args = []string{"grok", "open", req.SessionID}
	return nil
}
```
