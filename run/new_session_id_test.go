package run

import (
	"strings"
	"testing"
)

func TestSlugifyNewPrompt(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"/brainstorm fix flaky auth", "brainstorm-fix-flaky-auth"},
		{"$brainstorm extract TODOs", "brainstorm-extract-todos"},
		{"!!!", "sess"},
		{"", "sess"},
	}
	for _, tc := range cases {
		got := slugifyNewPrompt(tc.in)
		if got != tc.want {
			t.Fatalf("slugifyNewPrompt(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestBuildNewPrompt(t *testing.T) {
	if got := buildNewPrompt(newGrok, "fix it"); got != "/brainstorm fix it" {
		t.Fatalf("grok prompt=%q", got)
	}
	if got := buildNewPrompt(newCodex, "fix it"); got != "$brainstorm fix it" {
		t.Fatalf("codex prompt=%q", got)
	}
}

func TestBuildAgentRunNewArgv(t *testing.T) {
	argv := buildAgentRunNewArgv("grok-tty", "/tmp/proj", "sid-1", "/brainstorm x", false, true)
	join := strings.Join(argv, " ")
	for _, want := range []string{
		"run", "--open", "--color", "grok-tty", "--session-id", "sid-1",
		"--no-submit", "--new-terminal", "--dir", "/tmp/proj", "/brainstorm x",
	} {
		if !strings.Contains(join, want) {
			t.Fatalf("argv missing %q: %v", want, argv)
		}
	}
	argvSubmit := buildAgentRunNewArgv("codex-tty", "/tmp", "s2", "$brainstorm y", true, false)
	join = strings.Join(argvSubmit, " ")
	if strings.Contains(join, "--no-submit") {
		t.Fatalf("submit must omit --no-submit: %v", argvSubmit)
	}
	if strings.Contains(join, "--new-terminal") {
		t.Fatalf("here must omit --new-terminal: %v", argvSubmit)
	}
}
