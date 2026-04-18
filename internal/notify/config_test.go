package notify

import "testing"

func TestParseChannel_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  Channel
	}{
		{"stdout", ChannelStdout},
		{"STDOUT", ChannelStdout},
		{"", ChannelStdout},
		{"exec", ChannelExec},
		{"Exec", ChannelExec},
	}
	for _, tc := range cases {
		got, err := ParseChannel(tc.input)
		if err != nil {
			t.Errorf("ParseChannel(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseChannel(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseChannel_Invalid(t *testing.T) {
	_, err := ParseChannel("slack")
	if err == nil {
		t.Error("expected error for unknown channel")
	}
}

func TestApplyFlags_Channel(t *testing.T) {
	cfg := DefaultConfig()
	if err := ApplyFlags(&cfg, "exec", "/usr/local/bin/alert.sh"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Channel != ChannelExec {
		t.Errorf("expected ChannelExec, got %q", cfg.Channel)
	}
	if cfg.ExecCmd != "/usr/local/bin/alert.sh" {
		t.Errorf("unexpected ExecCmd: %s", cfg.ExecCmd)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	cfg := DefaultConfig()
	if err := ApplyFlags(&cfg, "", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Channel != ChannelStdout {
		t.Errorf("expected default channel to remain stdout")
	}
}
