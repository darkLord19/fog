package main

import "testing"

func TestValidateSlackConfig(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		bot     string
		app     string
		wantErr bool
	}{
		{name: "http mode", mode: "http", wantErr: false},
		{name: "socket mode with tokens", mode: "socket", bot: "xoxb-123", app: "xapp-123", wantErr: false},
		{name: "socket mode missing bot", mode: "socket", app: "xapp-123", wantErr: true},
		{name: "socket mode missing app", mode: "socket", bot: "xoxb-123", wantErr: true},
		{name: "invalid mode", mode: "bad", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSlackConfig(tc.mode, tc.bot, tc.app)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for %s: %v", tc.name, err)
			}
		})
	}
}
