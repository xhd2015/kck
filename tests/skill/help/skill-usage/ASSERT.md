## Expected

- Exit 0; stderr empty.
- Usage mentions `--show`, `--install`, `--list`.
- Available topics includes `overview` and `send`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	if !strings.HasPrefix(resp.Stdout, "Usage: kck skill") {
		t.Fatalf("stdout=%q missing usage prefix", resp.Stdout)
	}
	for _, want := range []string{"--show", "--install", "--list", "Available topics:", "overview", "send"} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
}
```
