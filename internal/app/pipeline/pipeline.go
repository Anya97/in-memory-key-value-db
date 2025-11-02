package pipeline

import (
	"fmt"
	"in-memory-key-value-db/internal/app/compute"
	"in-memory-key-value-db/internal/app/parser"
)

type KVStorage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Pipeline struct {
	compute *compute.Compute
	storage KVStorage
}

func NewPipeline(c *compute.Compute, s KVStorage) *Pipeline {
	return &Pipeline{
		compute: c,
		storage: s,
	}
}

func (p *Pipeline) Execute(input string) error {
	command, err := p.compute.Parse(input)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	switch command.Name {
	case parser.Set:
		if err := p.storage.Set(command.Args[0], command.Args[1]); err != nil {
			return fmt.Errorf("set error: %w", err)
		}
	case parser.Get:
		value, err := p.storage.Get(command.Args[0])
		if err != nil {
			return fmt.Errorf("get %q: %w", command.Args[0], err)
		}

		fmt.Println(value)
	case parser.Del:
		if err := p.storage.Delete(command.Args[0]); err != nil {
			return fmt.Errorf("delete error: %w", err)
		}
	default:
		return fmt.Errorf("unknown command: %q", command.Name)

	}

	return nil
}
