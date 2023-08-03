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
	"github.com/begopher/cli/api"
	"strings"
)

func Cmd(name string, description string, statement Statement, opts api.Options, flgs api.Flags, args []string, command Command) cmd {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		panic("cli.Cmd: name cannot be empty")
	}
	if strings.HasPrefix(name, "-") {
		panic("cli.Parent: name cannot start with -")
	}
	if description == "" {
		panic("cli.Cmd: description cannot be empty")
	}
	if statement == nil {
		statement = NoStatement()
	}
	if opts == nil {
		opts = Options()
	}
	if flgs == nil {
		flgs = Flags()
	}
	if args == nil {
		args = []string{}
	}
	if command == nil {
		panic("cli.Cmd: command (user implementation) cannot be empty")
	}
	namespace := Namespace()
	namespace.AddAll(opts.Names())
	if err := namespace.AddAll(flgs.Names()); err != nil {
		panic("cli.Cmd: options and flags have identical names")
	}
	if err := namespace.Add(name); err != nil {
		msg := fmt.Sprintf("cli.Cmd: (%s) is identical to flag name or option name", name)
		panic(msg)
	}
	separator := len(args)
	for i, arg := range args {
		if arg == "" {
			panic("cli.Cmd: empty string is not allowed as argument")
		}
		if arg == "..." {
			separator = i
			break
		}
	}
	namedArgs := args[:separator]
	variadicArgs := []string{}
	if separator < len(args) {
		variadicArgs = args[separator+1:]
	}
	if separator == len(args)-1 {
		panic("cli.Cmd: hint must be given when using variadic arguments")
	}
	return cmd{
		name:         name,
		description:  description,
		statement:    statement,
		opts:         opts,
		flags:        flgs,
		command:      command,
		namedArgs:    namedArgs,
		variadicArgs: variadicArgs,
		namespace:    namespace,
		optionWidth:  opts.LNameWidth(),
		flagWidth:    flgs.LNameWidth(),
	}
}

type cmd struct {
	name         string
	description  string
	statement    Statement
	opts         api.Options
	flags        api.Flags
	command      Command
	namedArgs    []string
	variadicArgs []string
	namespace    api.Namespace
	optionWidth  int
	flagWidth    int
}

func (c cmd) Name() string {
	return c.name
}

func (c cmd) Description() string {
	return c.description
}

func (c cmd) Exec(path []string, options map[string]string, flags map[string]bool, args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	if c.name != args[0] {
		return false, nil
	}

	path = append(path, c.name)
	fullPath := strings.Join(path, " ")
	args = args[1:]
	args = c.extract(options, flags, args)
	c.opts.Default(options)
	c.flags.Default(flags)
	if len(args) > 0 {
		if args[0] == "--help" { // done
			return false, fmt.Errorf(c.usage(fullPath))
		}
		if args[0] == "--" {
			// done
			args = args[1:]
		} else if strings.HasPrefix(args[0], "-") {
			if c.opts.Has(args[0]) { // done
				msg := fmt.Sprintf("Error: missing value for %s option (e.g. %[1]s value)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			if c.opts.Count() > 0 && c.flags.Count() > 0 { // done
				msg := fmt.Sprintf("Error: unknown option or flag (%s)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			if c.opts.Count() > 0 { // done
				msg := fmt.Sprintf("Error: unknown option (%s)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			if c.flags.Count() > 0 { // done
				msg := fmt.Sprintf("Error: unknown flag (%s)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			if len(c.namedArgs) > 0 { // done
				msg := fmt.Sprintf("Error: double hyphens (--) is missing before (%s)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			if len(c.variadicArgs) > 0 { // done
				msg := fmt.Sprintf("Error: double hyphens (--) is missing before (%s)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			// done
			if len(args) == 1 {
				msg := fmt.Sprintf("Error: unexpected value (%s)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			msg := fmt.Sprintf("Error: unexpected values (%s)", strings.Join(args, ", "))
			return false, fmt.Errorf(c.usage(fullPath, msg))
		}
	} //end if invalid option of flag
	if len(args) < len(c.namedArgs) {
		missing := len(c.namedArgs) - len(args)
		start := len(c.namedArgs) - missing
		if missing == 1 { //done
			msg := fmt.Sprintf("Error: missing value for (%s) argument", strings.Join(c.namedArgs[start:], ", "))
			return false, fmt.Errorf(c.usage(fullPath, msg))
		}
		// done
		msg := fmt.Sprintf("Error: missing values for (%s) arguments", strings.Join(c.namedArgs[start:], ", "))
		return false, fmt.Errorf(c.usage(fullPath, msg))
	}
	if len(args) > len(c.namedArgs) && len(c.variadicArgs) == 0 {
		extraArgs := args[len(c.namedArgs):]
		if len(extraArgs) == 1 { // done
			msg := fmt.Sprintf("Error: unexpected value (%s)", extraArgs[0])
			return false, fmt.Errorf(c.usage(fullPath, msg))
		}
		// done
		msg := fmt.Sprintf("Error: unexpected values (%s)", strings.Join(extraArgs, ", "))
		return false, fmt.Errorf(c.usage(fullPath, msg))
	}
	mappedArgs := make(map[string]string, len(c.namedArgs))
	for i, key := range c.namedArgs {
		mappedArgs[key] = args[i]
	}
	variadicArgs := args[len(mappedArgs):]
	usage := func(summaries ...string) error {
		return fmt.Errorf(c.usage(fullPath, summaries...))
	}
	ctx := context(path, options, flags, mappedArgs, variadicArgs, args, usage)
	if err := c.command.Exec(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (c cmd) extract(options map[string]string, flags map[string]bool, args []string) []string {
	length := len(args)
	args = c.opts.Extract(options, args)
	args = c.flags.Extract(flags, args)
	if length != len(args) {
		args = c.extract(options, flags, args)
	}
	return args
}

func (c cmd) String(width int) string {
	return fmt.Sprintf("%-[1]*s  %s\n", width, c.name, c.description)
}

func (c cmd) usage(path string, summaries ...string) string {
	var text, args strings.Builder
	if len(c.namedArgs)+len(c.variadicArgs) > 0 {
		args.WriteString("[--] ")
		args.WriteString(strings.Join(c.namedArgs, " "))
		if len(c.namedArgs) > 0 {
			args.WriteString(" ")
		}
		args.WriteString(strings.Join(c.variadicArgs, " "))
	}
	var optFlg string
	if c.opts.Count() > 0 && c.flags.Count() > 0 {
		optFlg = "[OPTIONS|FLAGS] "
	} else if c.opts.Count() > 0 {
		optFlg = "[OPTIONS] "
	} else if c.flags.Count() > 0 {
		optFlg = "[FLAGS] "
	}
	text.WriteString(fmt.Sprintf("\nUsage: %s %s%s\n\n", path, optFlg, args.String()))
	text.WriteString(fmt.Sprintf("%s\n\n", c.description))
	text.WriteString(c.opts.String(c.optionWidth))
	text.WriteString(c.flags.String(c.flagWidth))
	if len(summaries) > 0 {
		for _, msg := range summaries {
			text.WriteString(msg)
			text.WriteString("\n")
		}
		text.WriteString("\n")
	}
	text.WriteString(c.statement.String(path))
	return text.String()
}

func (c cmd) Namespace() api.Namespace {
	return c.namespace
}

func (c cmd) Help() string {
	return c.usage(c.name)
}
