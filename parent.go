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

func Parent(name, description string, statement Statement, options api.Options, flags api.Flags, manyCmds ...api.Cmd) api.Cmd {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		panic("cmd.Parent: cannot be created from empty name")
	}
	if description == "" {
		panic("cmd.Parent: cannot be created from empty description")
	}
	if strings.HasPrefix(name, "-") {
		panic("cli.Parent: name cannot start with -")
	}
	if len(manyCmds) < 1 {
		panic("cmd.Parent: cannot be created from empty/nil cmds")
	}
	if statement == nil {
		statement = NoStatement()
	}
	if options == nil {
		options = Options()
	}
	if flags == nil {
		flags = Flags()
	}
	_cmds := cmds(manyCmds)
	namespaces := _cmds.Namespace()
	if err := namespaces.Add(name); err != nil {
		msg := fmt.Sprintf("cli.Parent: name(%s) is duplicated, with a cmd child or one of its flag/option", name)
		panic(msg)
	}
	if err := namespaces.AddAll(options.Names()); err != nil {
		msg := fmt.Sprintf("cli.Parent: option name (%s) is used by a cmd, option or flag", err)
		panic(msg)
	}
	if err := namespaces.AddAll(flags.Names()); err != nil {
		msg := fmt.Sprintf("cli.Parent: flag name (%s) is used by a cmd, option or flag ", err)
		panic(msg)
	}
	return parent{
		name:        name,
		description: description,
		statement:   statement,
		options:     options,
		flags:       flags,
		cmds:        _cmds,
		flagWidth:   flags.LNameWidth(),
		optionWidth: options.LNameWidth(),
		namespace:   namespaces,
	}
}

type parent struct {
	name        string
	description string
	statement   Statement
	options     api.Options
	flags       api.Flags
	cmds        api.Cmds
	flagWidth   int
	optionWidth int
	namespace   api.Namespace
}

func (p parent) Name() string {
	return p.name
}

func (p parent) Description() string {
	return p.description
}

func (p parent) Exec(path []string, options map[string]string, flags map[string]bool, args []string) (bool, error) {
	path = append(path, p.name)
	fullPath := strings.Join(path, " ")
	if len(args) == 0 {
		return false, nil
	}
	if args[0] != p.name {
		return false, nil
	}
	args = args[1:]
	args = p.extract(options, flags, args)
	p.options.Default(options)
	p.flags.Default(flags)
	if len(args) == 0 {
		return false, fmt.Errorf(p.usage(fullPath, "Error: no command was selected"))
	}
	if args[0] == "--help" {
		return false, fmt.Errorf(p.usage(fullPath))
	}
	if strings.HasPrefix(args[0], "-") {
		if p.options.Has(args[0]) {
			msg := fmt.Sprintf("Error: missing value for %s option (e.g. %[1]s value)", args[0])
			return false, fmt.Errorf(p.usage(fullPath, msg))
		}
		if p.options.Count() > 0 && p.flags.Count() > 0 {
			msg := fmt.Sprintf("Error: unknown option or flag (%s)", args[0])
			return false, fmt.Errorf(p.usage(fullPath, msg))
		}
		if p.options.Count() > 0 {
			msg := fmt.Sprintf("Error: unknown option (%s)", args[0])
			return false, fmt.Errorf(p.usage(fullPath, msg))
		}
		if p.flags.Count() > 0 {
			msg := fmt.Sprintf("Error: unknown flag (%s)", args[0])
			return false, fmt.Errorf(p.usage(fullPath, msg))
		}
		msg := fmt.Sprintf("Error: unknown command (%s)", args[0])
		return false, fmt.Errorf(p.usage(fullPath, msg))
	}
	ok, err := p.cmds.Exec(path, options, flags, args)
	if err != nil {
		return ok, err
	}
	if ok {
		return ok, err
	}
	msg := fmt.Sprintf("Error: unknown command (%s)", args[0])
	return false, fmt.Errorf(p.usage(fullPath, msg))
}

func (p parent) extract(options map[string]string, flags map[string]bool, args []string) []string {
	length := len(args)
	args = p.options.Extract(options, args)
	args = p.flags.Extract(flags, args)
	if length != len(args) {
		args = p.extract(options, flags, args)
	}
	return args
}

func (p parent) usage(path string, errors ...string) string {
	var optFlg string
	if p.options.Count() > 0 && p.flags.Count() > 0 {
		optFlg = "[OPTIONS|FLAGS] "
	} else if p.options.Count() > 0 {
		optFlg = "[OPTIONS] "
	} else if p.flags.Count() > 0 {
		optFlg = "[FLAGS] "
	}
	var text strings.Builder
	text.WriteString(fmt.Sprintf("\nUsage: %s %sCOMMAND\n\n", path, optFlg))
	text.WriteString(fmt.Sprintf("%s\n\n", p.description))
	text.WriteString("Commands:\n")
	text.WriteString(p.cmds.String())
	text.WriteString("\n")
	text.WriteString(p.options.String(p.optionWidth))
	text.WriteString(p.flags.String(p.flagWidth))
	if len(errors) > 0 {
		for _, msg := range errors {
			text.WriteString(msg)
			text.WriteString("\n")
		}
		text.WriteString("\n")
	}
	text.WriteString(p.statement.String(path))
	return text.String()
}

func (p parent) Namespace() api.Namespace {
	return p.namespace
}

func (p parent) String(width int) string {
	return fmt.Sprintf("%-[1]*s  %s\n", width, p.name, p.description)
}

func (p parent) Help() string {
	return p.usage(p.name)
}
