# Scenario

**Feature**: `--home` only sees sessions under provided home

```
seed only under req.Home; --home req.Home -> lists seeded id;
other home path not used (inject Home + flag same path; isolation via seed only there)
```

## Steps

1. Seed `only-here` under isolated Home.
2. Args include `--home` pointing at that Home (and Options.Home same).
3. Assert listed id is only-here (proves store path is the isolated one).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Sessions = []SessionSeed{
		{
			SessionID: "only-here",
			Runner:    "grok",
			Workspace: "/ws/only",
			UpdatedAt: "2026-08-05T16:00:00Z",
			Live:      true,
			Sendable:  true,
			State:     "idle",
		},
	}
	// Flag isolation: --home must be honored (Options.Home also set by root).
	req.Args = []string{"--home", req.Home}
	return nil
}
```
