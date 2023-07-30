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

func NestedApp(name, description string, statement Statement, options api.Options, flags api.Flags, group api.Group, groups ...api.Group) application {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		panic("cli.Application: cannot be created from empty name")
	}
	if description == "" {
		panic("cli.Application: cannot be created from empty description")
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
	if group == nil {
		panic("cli.Application: cannot be created from nil group")
	}
	grps := make([]api.Group, len(groups) + 1)
	grps[0] = group
	for i, group := range groups {
		grps[i+1] = group
	}
	_groups := Groups(grps)
	namespace := _groups.Namespace()
	if err := namespace.Add(name); err != nil {
		msg := fmt.Sprintf("cli.Application: application name (%s) is used by cmd, option or flag", name)
		panic(msg)
	}
	if err := namespace.AddAll(options.Names()); err != nil {
		msg := fmt.Sprintf("cli.Application: option name (%s) is used by cmd, option or flag", err)
		panic(msg)
	}
	if err := namespace.AddAll(flags.Names()); err != nil {
		msg := fmt.Sprintf("cli.Application: flag name (%s) is used by cmd, option or flag ", err)
		panic(msg)
	}
	if err := namespace.Add("help"); err != nil {
		panic("cli.NestedApp: help cannot be used as a name of any object (reserved for --help)")
	}
	return application{
		name:        name,
		description: description,
		statement:   statement,
		options:        options,
		flags:        flags,
		groups:  _groups,
		flagWidth:  flags.LNameWidth(),
		optionWidth: options.LNameWidth(),
	}
}

type application struct {
	name        string
	description string
	statement Statement
	options        api.Options
	flags        api.Flags
	groups api.Groups
	flagWidth int
	optionWidth int
}

func (a application) Run(args []string) error{
	if len(args) == 0 {
		return fmt.Errorf(a.usage())
	}
	if a.name != args[0] {
		return fmt.Errorf(a.usage())
		
	}
	path := []string{a.name}
	args = args[1:]
	options := make(map[string]string, 0)
	flags := make(map[string]bool, 0)
	args = a.extract(options, flags, args)
	a.options.Default(options)
	a.flags.Default(flags)
	if len(args) == 0 {
		return fmt.Errorf(a.usage("Error: no command was selected"))
	}
	if args[0] == "--help" {
		return fmt.Errorf(a.usage())
	}
	if strings.HasPrefix(args[0], "-") { // command name never start with -
		if a.options.Has(args[0]) {
			msg := fmt.Sprintf("Error: missing value for %s option", args[0])
			return fmt.Errorf(a.usage(msg))
		}
		var optFlg string
		if a.options.Count() > 0 && a.flags.Count() > 0{
			optFlg = "option or flag"
		}else if a.options.Count() > 0 {
			optFlg = "option"
		}else if a.flags.Count() > 0 {
			optFlg = "flag"
		}else{
			optFlg = "command"
		}
		msg := fmt.Sprintf("Error: unknown %s (%s)", optFlg, args[0])
		return fmt.Errorf(a.usage(msg))
	}
	ok, err := a.groups.Exec(path, options, flags, args)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	msg := fmt.Sprintf("Error: unknown command (%s)", args[0])
	return fmt.Errorf(a.usage(msg))
}

func (a application) extract(options map[string]string, flags map[string]bool, args []string) []string {
	length := len(args)
	args = a.options.Extract(options, args)
	args = a.flags.Extract(flags, args)
	if length != len(args) {
		args = a.extract(options, flags, args)
	}
	return args
}

func (a application) usage(errors ...string) string {
	var optFlg string
	if a.options.Count() > 0 && a.flags.Count() > 0{
		optFlg = "[OPTIONS|FLAGS] "
	}else if a.options.Count() > 0 {
		optFlg = "[OPTIONS] "
	}else if a.flags.Count() > 0 {
		optFlg = "[FLAGS] "
	}
	var text strings.Builder
	text.WriteString(fmt.Sprintf("\nUsage: %s %sCOMMAND\n\n", a.name, optFlg))
	text.WriteString(fmt.Sprintf("%s\n\n", a.description))
	text.WriteString(a.groups.String())
	text.WriteString(a.options.String(a.optionWidth))
	text.WriteString(a.flags.String(a.flagWidth))
	if len(errors) > 0 {
		for _, msg := range errors{
			text.WriteString(msg)
			text.WriteString("\n")
		}
		text.WriteString("\n")
	}
	text.WriteString(a.statement.String(a.name))
	return text.String()
}
