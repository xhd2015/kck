# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckInfoSessionID
	writeKckInfoSession(t, req)
	req.Procs = []FixtureProc{{PID: 6001, PPID: 1, Cmd: "/usr/local/bin/codex"}}
	req.OpenFiles[6001] = []string{codexOpenPath(req.SessionID)}
	req.Args = []string{"codex", "info", req.SessionID, "--no-pid"}
	return nil
}
```
