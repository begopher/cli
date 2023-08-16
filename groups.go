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

func groups(grps []api.Group) _groups {
	if len(grps) == 0 {
		panic("cli.Groups: grps cannot be empty slice")
	}
	xnamespaces := make([]api.Namespace, len(grps))
	groupNamespace := namespace()
	cmdNamespace := namespace()
	for i, group := range grps {
		if group == nil {
			panic("cli.Groups: nil value is not allowed as a valid Group")
		}
		xnamespaces[i] = group.Namespace()
		if err := groupNamespace.Add(group.Name()); err != nil {
			msg := fmt.Sprintf("cli.Groups: name (%s) is taken by two group", group.Name())
			panic(msg)
		}
		if err := cmdNamespace.AddAll(group.Names()); err != nil {
			msg := fmt.Sprintf("cli.Groups: cmd name (%s) in group (%s) is taken by other cmd in other group", err, group.Name())
			panic(msg)
		}
	}
	return _groups{
		grps:      grps,
		namespace: namespaces(xnamespaces),
	}
}

type _groups struct {
	grps      []api.Group
	namespace api.Namespace
}

func (g _groups) Exec(path []string, options map[string]string, flags map[string]bool, args []string) (bool, error) {
	for _, group := range g.grps {
		ok, err := group.Exec(path, options, flags, args)
		if err != nil {
			return ok, err
		}
		if ok {
			return ok, err
		}
	}
	return false, nil
}

func (g _groups) String() string {
	var text strings.Builder
	for _, group := range g.grps {
		text.WriteString("\n")
		text.WriteString(group.String())
	}
	return text.String()
}

func (g _groups) Namespace() api.Namespace {
	return g.namespace
}
