package cli

import (
	"fmt"
	"strings"
)

func Argument(name, description string) argument {
	name = strings.TrimSpace(name)
	if name == "" {
		panic("cli.Argument: name cannot be empty")
	}
	description = strings.TrimSpace(description)
	return argument{
		name:        name,
		description: description,
	}
}

type argument struct {
	name        string
	description string
}

func (a argument) Name() string {
	return a.name
}

func (a argument) Description() string {
	return a.description
}

func (a argument) Extract(namedArgs map[string]string, args []string) []string {
	if len(args) == 0 {
		return args
	}
	namedArgs[a.name] = args[0]
	return args[1:]
}
func (a argument) String(width int) string {
	return fmt.Sprintf("  %-[1]*s  %s\n", width, a.name, a.description)
}
