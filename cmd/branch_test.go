package cmd

import (
	"testing"

	"github.com/owenrumney/branch/internal/config"
)

func TestParseArgs(t *testing.T) {
	cfg := config.Default()

	tests := []struct {
		name       string
		args       []string
		cfg        *config.Config
		wantTicket string
		wantDesc   []string
	}{
		{
			name:       "with ticket at start",
			args:       []string{"PIP-1234", "implement", "feature"},
			cfg:        cfg,
			wantTicket: "PIP-1234",
			wantDesc:   []string{"implement", "feature"},
		},
		{
			name:       "with GitHub issue ticket",
			args:       []string{"#123", "fix", "bug"},
			cfg:        cfg,
			wantTicket: "#123",
			wantDesc:   []string{"fix", "bug"},
		},
		{
			name:       "without ticket",
			args:       []string{"implement", "new", "feature"},
			cfg:        cfg,
			wantTicket: "",
			wantDesc:   []string{"implement", "new", "feature"},
		},
		{
			name:       "single word not a ticket",
			args:       []string{"feature"},
			cfg:        cfg,
			wantTicket: "",
			wantDesc:   []string{"feature"},
		},
		{
			name:       "empty args",
			args:       []string{},
			cfg:        cfg,
			wantTicket: "",
			wantDesc:   nil,
		},
		{
			name:       "ticket only",
			args:       []string{"PIP-5678"},
			cfg:        cfg,
			wantTicket: "PIP-5678",
			wantDesc:   []string{},
		},
		{
			name:       "underscore ticket format",
			args:       []string{"PIP_1234", "add", "tests"},
			cfg:        cfg,
			wantTicket: "PIP_1234",
			wantDesc:   []string{"add", "tests"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTicket, gotDesc := parseArgs(tt.args, tt.cfg)

			if gotTicket != tt.wantTicket {
				t.Errorf("parseArgs() ticket = %q, want %q", gotTicket, tt.wantTicket)
			}

			if len(gotDesc) != len(tt.wantDesc) {
				t.Errorf("parseArgs() description length = %d, want %d", len(gotDesc), len(tt.wantDesc))
			}

			for i, want := range tt.wantDesc {
				if i >= len(gotDesc) {
					t.Errorf("parseArgs() description[%d] missing, want %q", i, want)
					continue
				}
				if gotDesc[i] != want {
					t.Errorf("parseArgs() description[%d] = %q, want %q", i, gotDesc[i], want)
				}
			}
		})
	}
}
