package compute

import (
	"in-memory-key-value-db/internal/app/parser"
)

type Compute struct{}

func NewCompute() *Compute { return &Compute{} }

func (c *Compute) Parse(input string) (parser.Command, error) {
	return parser.ParseInputString(input)
}
