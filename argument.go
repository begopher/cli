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
