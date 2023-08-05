package cli

import (
	"fmt"
	"strings"
)

func NoVariadic() noVariadic {
	return noVariadic{}
}

type noVariadic struct {}

func (v noVariadic) Arg() string {
	return ""
}

func (v noVariadic) Allowed() bool {
	return false
}

func (v noVariadic) Extract(args []string) ([]string, error) {
	if len(args) == 0 {
		return args, nil
	}
	if len(args) == 1 {
		msg := fmt.Sprintf("Error: unexpected value (%s)", args[0])
		return args, fmt.Errorf(msg)
		
	}
	msg := fmt.Sprintf("Error: unexpected values (%s)", strings.Join(args, ", "))	
	return args, fmt.Errorf(msg)
}

func (v noVariadic) String() string {
	return ""
}
