package compute

import (
	"fmt"
	"in-memory-key-value-db/internal/app/parser"
)

type KVStorage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Computer struct {
	storage KVStorage
}

func NewComputer(storage KVStorage) *Computer {
	return &Computer{storage: storage}
}

func (c *Computer) Execute(input string) error {
	command, commandErr := parser.ParseInputString(input)
	if commandErr != nil {
		return commandErr
	}

	switch command.Name {
	case parser.Set:
		err := c.storage.Set(command.Args[0], command.Args[1])
		if err != nil {
			return fmt.Errorf("set %q: %w", command.Args[0], err)
		}
	case parser.Get:
		value, err := c.storage.Get(command.Args[0])
		if err != nil {
			return fmt.Errorf("get %q: %w", command.Args[0], err)
		}

		fmt.Println(value)
	case parser.Del:
		err := c.storage.Delete(command.Args[0])
		if err != nil {
			return fmt.Errorf("del %q: %w", command.Args[0], err)
		}
	}
	return nil
}
