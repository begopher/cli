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
)

func SimpleApp(name, description string, statement Statement, options api.Options, flags api.Flags, args api.Arguments, vars api.Variadic, implementation Command) simpleApp {
	name = removeAbsolutePath(name)
	cmd := Cmd(
		name,
		description,
		statement,
		options,
		flags,
		args,
		vars,
		implementation)
	return simpleApp{
		cmd: cmd,
	}
}

type simpleApp struct {
	cmd api.Cmd
}

func (s simpleApp) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf(s.cmd.Help())
	}
	args[0] = removeAbsolutePath(args[0])
	options := make(map[string]string, 0)
	flags := make(map[string]bool, 0)
	path := make([]string, 0)
	ok, err := s.cmd.Exec(path, options, flags, args)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(s.cmd.Help())
	}
	return nil
}
