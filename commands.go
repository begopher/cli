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
	"github.com/begopher/cli/internal/api"
	"strings"
)

func commands(cmds []api.Command) _commands {
	if len(cmds) < 1 {
		panic("cli.Cmds: cannot be created from empty or nil slice")
	}
	xNamespaces := make([]api.Namespace, len(cmds))
	sibling := namespace()
	var nameWidth int
	for i, cmd := range cmds {
		if cmd == nil {
			panic("cli.Cmds: nil is not accepted as a valid cmd value")
		}
		if err := sibling.Add(cmd.Name()); err != nil {
			msg := fmt.Sprintf("cli.Cmds: name (%s) is taken by other Cmd", cmd.Name())
			panic(msg)
		}
		xNamespaces[i] = cmd.Namespace()
		if width := len([]rune(cmd.Name())); width > nameWidth {
			nameWidth = width
		}
	}
	return _commands{
		cmds:      cmds,
		namespace: namespaces(xNamespaces),
		nameWidth: nameWidth,
	}
}

type _commands struct {
	cmds      []api.Command
	namespace api.Namespace
	nameWidth int
}

func (c _commands) Exec(path []string, options map[string]string, flags map[string]bool, args []string) (bool, error) {
	for _, cmd := range c.cmds {
		ok, err := cmd.Exec(path, options, flags, args)
		if err != nil {
			return ok, err
		}
		if ok {
			return ok, err
		}
	}
	return false, nil
}

func (c _commands) Names() []string {
	names := make([]string, len(c.cmds))
	for i, cmd := range c.cmds {
		names[i] = cmd.Name()
	}
	return names
}

func (c _commands) String() string {
	var text strings.Builder
	for _, cmd := range c.cmds {
		text.WriteString("  ")
		text.WriteString(cmd.String(c.nameWidth))
	}
	return text.String()
}

func (c _commands) Namespace() api.Namespace {
	return c.namespace
}
