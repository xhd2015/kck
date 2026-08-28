# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSnapshotSessionID
	writeKckSnapshotSession(t, req)
	req.Procs = nil
	req.OpenFiles = map[int][]string{}
	req.ITerm = oneITermTab()
	req.Args = []string{"codex", "snapshot", req.SessionID}
	return nil
}
```
