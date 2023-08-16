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

func Command(name string, description string, statement Statement, opts api.Options, flgs api.Flags, arguments api.Arguments, variadic api.Variadic, implementation Implementation) command {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		panic("cli.Command: name cannot be empty")
	}
	if strings.HasPrefix(name, "-") {
		panic("cli.Command: name cannot start with -")
	}
	if description == "" {
		panic("cli.Command: description cannot be empty")
	}
	if statement == nil {
		panic("cli.Command: statement cannot be nil")
	}
	if opts == nil {
		panic("cli.Command: opts cannot be nil")
	}
	if flgs == nil {
		panic("cli.Command: flgs cannot be nil")
	}
	if arguments == nil {
		panic("cli.Command: arguments cannot be nil")
	}
	if variadic == nil {
		panic("cli.Command: variadic cannot be nil")
	}
	if implementation == nil {
		panic("cli.Command: implementation cannot be nil")
	}
	namespace := namespace()
	namespace.AddAll(opts.Names())
	if err := namespace.AddAll(flgs.Names()); err != nil {
		panic("cli.Command: options and flags have identical names")
	}
	if err := namespace.Add(name); err != nil {
		msg := fmt.Sprintf("cli.Command: name (%s) is identical to a flag name or an option name", name)
		panic(msg)
	}
	return command{
		name:           name,
		description:    description,
		statement:      statement,
		opts:           opts,
		flags:          flgs,
		arguments:      arguments,
		variadic:       variadic,
		implementation: implementation,
		namespace:      namespace,
	}
}

type command struct {
	name           string
	description    string
	statement      Statement
	opts           api.Options
	flags          api.Flags
	implementation Implementation
	arguments      api.Arguments
	variadic       api.Variadic
	namespace      api.Namespace
}

func (c command) Name() string {
	return c.name
}

func (c command) Description() string {
	return c.description
}

func (c command) Exec(path []string, options map[string]string, flags map[string]bool, args []string) (bool, error) {
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
			if c.arguments.Count() > 0 { // done
				msg := fmt.Sprintf("Error: double hyphens (--) is missing before (%s)", args[0])
				return false, fmt.Errorf(c.usage(fullPath, msg))
			}
			if c.variadic.Allowed() { // done
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
	namedArgs := make(map[string]string, c.arguments.Count())
	args, err := c.arguments.Extract(namedArgs, args)
	if err != nil {
		return false, fmt.Errorf(c.usage(fullPath, err.Error()))
	}
	variadicArgs, err := c.variadic.Extract(args)
	if err != nil {
		return false, fmt.Errorf(c.usage(fullPath, err.Error()))
	}
	usage := func(summaries ...string) error {
		return fmt.Errorf(c.usage(fullPath, summaries...))
	}
	ctx := context(path, options, flags, namedArgs, variadicArgs, usage)
	if err := c.implementation.Exec(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (c command) extract(options map[string]string, flags map[string]bool, args []string) []string {
	length := len(args)
	args = c.opts.Extract(options, args)
	args = c.flags.Extract(flags, args)
	if length != len(args) {
		args = c.extract(options, flags, args)
	}
	return args
}

func (c command) String(width int) string {
	return fmt.Sprintf("%-[1]*s  %s\n", width, c.name, c.description)
}

func (c command) usage(path string, summaries ...string) string {
	var text, args strings.Builder
	if c.arguments.Count() > 0 || c.variadic.Allowed() {
		args.WriteString("[--] ")
		args.WriteString(strings.Join(c.arguments.Names(), " "))
		if c.arguments.Count() > 0 {
			args.WriteString(" ")
		}
		args.WriteString(c.variadic.Arg())
	}
	var optFlg string
	if c.opts.Count() > 0 && c.flags.Count() > 0 {
		optFlg = "[OPTIONS|FLAGS] "
	} else if c.opts.Count() > 0 {
		optFlg = "[OPTIONS] "
	} else if c.flags.Count() > 0 {
		optFlg = "[FLAGS] "
	}
	text.WriteString(fmt.Sprintf("Usage: %s %s%s\n\n", path, optFlg, args.String()))
	text.WriteString(fmt.Sprintf("%s\n", c.description))
	text.WriteString(c.opts.String())      //done
	text.WriteString(c.flags.String())     //done
	text.WriteString(c.arguments.String()) //done
	text.WriteString(c.variadic.String())  //done
	if len(summaries) > 0 {
		text.WriteString("\n")
		for _, msg := range summaries {
			text.WriteString(msg)
			text.WriteString("\n")
		}
	}
	if !c.statement.Empty() {
		text.WriteString("\n")
	}
	text.WriteString(c.statement.String(path))
	return text.String()
}

func (c command) Namespace() api.Namespace {
	return c.namespace
}

func (c command) Help() string {
	return c.usage(c.name)
}
