package cli

import (
	"fmt"
	"strings"
)

func Variadic(arg, description string) variadic {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		panic("cli.Variadic: arg cannot be empty, look at cli.NoVariadic()")
	}
	description = strings.TrimSpace(description)
	return variadic{
		arg:         arg,
		description: description,
	}
}

type variadic struct {
	arg         string
	description string
}

func (v variadic) Arg() string {
	return fmt.Sprintf("[%s]", v.arg)
}

func (v variadic) Allowed() bool {
	return true
}

func (v variadic) Extract(args []string) ([]string, error) {
	return args, nil
}

func (v variadic) String() string {
	if v.description == "" {
		return ""
	}
	var text strings.Builder
	text.WriteString("Variadic:\n")
	msg := fmt.Sprintf("  [%s]  %s\n", v.arg, v.description)
	text.WriteString(msg)
	text.WriteString("\n")
	return text.String()
}
