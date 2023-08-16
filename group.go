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

func Group(name string, cmds ...api.Command) group {
	if name == "" {
		panic("cli.Group: name cannot be empty")
	}
	if len(cmds) == 0 {
		panic("cli.Group: cmds cannot be empty")
	}
	for _, cmd := range cmds {
		if cmd == nil {
			panic("cli.Group: nil value is not allowed as a Command")
		}
	}
	return group{
		name: name,
		commands: commands(cmds),
	}
}

type group struct {
	name string
	commands api.Commands
}

func (g group) Exec(path []string, options map[string]string, flags map[string]bool, args []string) (bool, error) {
	ok, err := g.commands.Exec(path, options, flags, args)
	if err != nil {
		return ok, err
	}
	if ok {
		return ok, err
	}
	return false, nil
}

func (g group) Name() string {
	return g.name
}

func (g group) Names() []string {
	return g.commands.Names()
}

func (g group) Namespace() api.Namespace {
	return g.commands.Namespace()
}

func (g group) String() string {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("%s:\n", g.name))
	text.WriteString(g.commands.String())
	return text.String()
}
