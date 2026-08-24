# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckInfoSessionID
	writeKckInfoSession(t, req)
	writeActive(t, req, req.SessionID)
	req.Procs = []FixtureProc{{PID: 6001, PPID: 1, Cmd: "/usr/local/bin/grok"}}
	req.OpenFiles[6001] = []string{grokOpenPath(req.SessionID)}
	req.Args = []string{"grok", "info", req.SessionID, "--no-pid"}
	return nil
}
```
