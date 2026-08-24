package run

import (
	"reflect"
	"testing"
)

func TestPeelCronFlag(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		args    []string
		want    string
		remain  []string
		wantErr string
	}{
		{
			name:   "absent",
			args:   []string{"hi", "--session-id", "abc"},
			remain: []string{"hi", "--session-id", "abc"},
		},
		{
			name:   "space form",
			args:   []string{"hi", "--cron", "every-1h", "--session-id", "abc"},
			want:   "every-1h",
			remain: []string{"hi", "--session-id", "abc"},
		},
		{
			name:   "equals form",
			args:   []string{"--cron=every-5m", "hi", "--session-id", "abc"},
			want:   "every-5m",
			remain: []string{"hi", "--session-id", "abc"},
		},
		{
			name:    "missing value",
			args:    []string{"hi", "--cron"},
			wantErr: "--cron requires an expression",
		},
		{
			name:    "duplicate",
			args:    []string{"--cron", "every-1h", "--cron", "every-5m"},
			wantErr: "duplicate --cron",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, remain, err := peelCronFlag(tc.args)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("err=%v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("cron=%q, want %q", got, tc.want)
			}
			if !reflect.DeepEqual(remain, tc.remain) {
				t.Fatalf("remain=%v, want %v", remain, tc.remain)
			}
		})
	}
}
