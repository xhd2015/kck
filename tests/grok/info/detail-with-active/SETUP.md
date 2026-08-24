# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckInfoSessionID
	writeKckInfoSession(t, req)
	writeActive(t, req, req.SessionID)
	req.Procs = []FixtureProc{{PID: 5001, PPID: 1, Cmd: "/usr/local/bin/grok"}}
	req.OpenFiles[5001] = []string{grokOpenPath(req.SessionID)}
	req.Args = []string{"grok", "info", req.SessionID}
	return nil
}
```
