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

	"github.com/begopher/cli/api"
)

func Group(name string, cmds ...api.Cmd) group {
	if name == "" {
		panic("group.New: cannot be created from empty name")
	}
	if len(cmds) < 1 {
		panic("group.New: cmds cannot be empty or nil")
	}
	return group{
		name:    name,
		cmds:    Cmds(cmds),
	}
}

type group struct {
	name    string
	cmds    api.Cmds
}

func (g group) Exec(path []string, options map[string]string, flags map[string]bool, args []string) (bool, error){
	ok, err := g.cmds.Exec(path, options, flags, args)
	if err != nil {
		return ok, err
	}
	if ok {
		return ok, err
	}
	return false, nil
}

func (g group) Name() string{
	return g.name
}

func (g group) Names() []string {
	return g.cmds.Names()
}

func (g group) Namespace() api.Namespace {
	return g.cmds.Namespace()
}

func (g group) String() string {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("%s:\n", g.name))
	text.WriteString(g.cmds.String())
	return text.String()
}
