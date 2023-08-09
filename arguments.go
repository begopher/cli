//   Copyright 2023 Abdulrahman Abdulhamid Alsaedi
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package cli

import (
	"fmt"
	"strings"

	"github.com/begopher/cli/internal/api"
)

// Arguments represent a collection of zero or more Argument.
//
// Client of cli library should not invoke any method of Arguments directly,
// instead, Arguments should be passed to:
//   - cli.Cmd(..., Arguments, ...) function.
//
// # Panic when:
//   - two arguments has the same name.
//   - some arguments have non-empty description, while others are empty.
//     (all arguments must have description or all must be empty).
func Arguments(args ...api.Argument) arguments {
	if len(args) == 0 {
		return arguments{
			documented: false,
			width:      0,
			args:       args,
		}
	}
	namespace := namespace()
	var width int
	for _, arg := range args {
		if err := namespace.Add(arg.Name()); err != nil {
			msg := fmt.Sprintf("cli.Arguments: duplicated argument name (%s)", arg.Name())
			panic(msg)
		}
		length := len(arg.Name())
		if length > width {
			width = length
		}
	}
	documented := args[0].Description() != ""
	for _, arg := range args[1:] {
		isDoc := arg.Description() != ""
		if documented != isDoc {
			panic("cli.Arguments: all arguments must either have an empty or non-empty description")
		}
	}
	return arguments{
		documented: documented,
		width:      width,
		args:       args,
	}
}

type arguments struct {
	documented bool
	width      int
	args       []api.Argument
}

func (a arguments) Names() []string {
	names := make([]string, len(a.args))
	for i, arg := range a.args {
		names[i] = arg.Name()
	}
	return names
}

func (a arguments) Count() int {
	return len(a.args)
}

func (a arguments) Extract(namedArgs map[string]string, args []string) ([]string, error) {
	if len(a.args) > len(args) {
		missing := len(a.args) - len(args)
		start := len(a.args) - missing
		if missing == 1 {
			arg := a.args[start]
			err := fmt.Errorf("Error: missing value for (%s) argument", arg.Name())
			return args, err
		}
		remains := a.args[start:]
		names := make([]string, len(remains))
		for i, arg := range remains {
			names[i] = arg.Name()
		}
		err := fmt.Errorf("Error: missing values for (%s) arguments", strings.Join(names, ", "))
		return args, err
	}
	for _, arg := range a.args {
		args = arg.Extract(namedArgs, args)
	}
	return args, nil
}

func (a arguments) String() string {
	if !a.documented {
		return ""
	}
	var text strings.Builder
	text.WriteString("\n")
	text.WriteString("Arguments:\n")
	for _, arg := range a.args {
		text.WriteString(arg.String(a.width))
	}
	return text.String()
}
