package parser

import (
	"strings"
	"testing"
)

func TestParseInputString(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		want            Command
		wantErrContains string
	}{
		{name: "empty input", input: "", wantErrContains: "empty input string"},
		{name: "unknown coomand", input: "FOO", wantErrContains: "unknown command"},
		{name: "wrong args set too few", input: "SET only", wantErrContains: "wrong number of arguments (want 2, got 1)"},
		{name: "wrong args set too many", input: "SET a b c", wantErrContains: "wrong number of arguments (want 2, got 3)"},
		{name: "wrong args get none", input: "GET", wantErrContains: "wrong number of arguments (want 1, got 0)"},
		{name: "wrong args get many", input: "GET a b", wantErrContains: "wrong number of arguments (want 1, got 2)"},
		{name: "wrong args del none", input: "DEL", wantErrContains: "wrong number of arguments (want 1, got 0)"},
		{name: "wrong args del many", input: "DEL a b", wantErrContains: "wrong number of arguments (want 1, got 2)"},
		{name: "invalid char arg1", input: "SET inv@lid val", wantErrContains: "invalid argument 1"},
		{name: "invalid char arg2", input: "SET key val@ue", wantErrContains: "invalid argument 2"},
		{name: "invalid hyphen", input: "SET a-b c", wantErrContains: "invalid argument 1"},
		{name: "valid set", input: "SET KEY1 VALUE_2", want: Command{Name: Set, Args: []string{"KEY1", "VALUE_2"}}},
		{name: "valid set wildcard", input: "SET path path/to/*", want: Command{Name: Set, Args: []string{"path", "path/to/*"}}},
		{name: "valid get", input: "GET KEY1", want: Command{Name: Get, Args: []string{"KEY1"}}},
		{name: "valid del", input: "DEL KEY1", want: Command{Name: Del, Args: []string{"KEY1"}}},
		{name: "trim spaces", input: "   SET   A   B   ", want: Command{Name: Set, Args: []string{"A", "B"}}},
		{name: "case sensitive", input: "set a b", wantErrContains: "unknown command"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInputString(tt.input)

			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErrContains, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tt.want.Name {
				t.Fatalf("Name mismatch: want %q, got %q", tt.want.Name, got.Name)
			}
			if len(got.Args) != len(tt.want.Args) {
				t.Fatalf("Args length mismatch: want %d, got %d (%v)", len(tt.want.Args), len(got.Args), got.Args)
			}
			for i := range got.Args {
				if got.Args[i] != tt.want.Args[i] {
					t.Fatalf("Arg[%d] mismatch: want %q, got %q", i, tt.want.Args[i], got.Args[i])
				}
			}
		})
	}
}
