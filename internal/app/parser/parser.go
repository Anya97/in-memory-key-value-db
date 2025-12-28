package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	Set = "SET"
	Get = "GET"
	Del = "DEL"
)

var argRe = regexp.MustCompile(`^[A-Za-z0-9_/*]+$`)

var allowedCommandsWithArgs = map[string]int{
	Set: 2, Get: 1, Del: 1,
}

var usage = map[string]string{
	Set: "SET key value",
	Get: "GET key",
	Del: "DEL key",
}

type Command struct {
	Name string
	Args []string
}

func ParseInputString(input string) (Command, error) {
	inputSlice := strings.Fields(input)

	if len(inputSlice) == 0 {
		return Command{}, errors.New("empty input string")
	}

	name := inputSlice[0]
	wantArgs, ok := allowedCommandsWithArgs[name]
	if !ok {
		return Command{}, fmt.Errorf(
			"unknown command: %q. Supported: %q, %q, %q. Usage: %s | %s | %s",
			name, Set, Get, Del, usage[Set], usage[Get], usage[Del],
		)
	}

	if len(inputSlice)-1 != wantArgs {
		return Command{}, fmt.Errorf(
			"%s: wrong number of arguments (want %d, got %d). Usage: %s",
			name, wantArgs, len(inputSlice)-1, usage[name],
		)
	}

	for i, a := range inputSlice[1:] {
		if !argRe.MatchString(a) {
			return Command{}, fmt.Errorf("invalid argument %d: %q. Allowed characters: letters, digits, '_', '/', '*', '-', '.'", i+1, a)
		}
	}

	return Command{Name: name, Args: inputSlice[1:]}, nil
}
