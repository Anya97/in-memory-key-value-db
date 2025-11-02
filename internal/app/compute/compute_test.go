package compute

import (
	"in-memory-key-value-db/internal/app/parser"
	"strings"
	"testing"
)

func TestCompute_Parse(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		want            parser.Command
		wantErrContains string
	}{
		{name: "empty input", input: "", wantErrContains: "empty input string"},
		{name: "unknown command", input: "FOO", wantErrContains: "unknown command"},
		{name: "valid set", input: "SET KEY1 VALUE_2", want: parser.Command{Name: parser.Set, Args: []string{"KEY1", "VALUE_2"}}},
		{name: "trim spaces", input: "   SET   A   B   ", want: parser.Command{Name: parser.Set, Args: []string{"A", "B"}}},
		{name: "case sensitive", input: "set a b", wantErrContains: "unknown command"},
	}

	c := NewCompute()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Parse(tt.input)

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
