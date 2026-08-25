## Expected

- `[dry-run]` plan on stdout; includes would-resolve fixture id; not bare-id-only.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{
		"[dry-run] start pid:",
		"[dry-run] ancestor pid:",
		"[dry-run] runner pid:",
		"[dry-run] would resolve: " + fixtureSessionID,
		"[dry-run] source:",
		"[dry-run] confidence:",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("dry-run stdout missing %q:\n%s", want, out)
		}
	}
	if strings.TrimSpace(out) == fixtureSessionID {
		t.Fatalf("dry-run must not be bare-id-only:\n%s", out)
	}
}
```
