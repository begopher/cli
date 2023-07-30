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

type Context interface {
	Flag(string) bool
	Flags() map[string]bool
	Option(string) string
	Options() map[string]string
	Arg(string) string
	Args() []string
	Variadic() []string
	Path() []string
	Usage(...string) error
}

func context(path []string, options map[string]string, flags map[string]bool, namedArgs map[string]string, variadicArgs, args []string, usage func(...string) error) _context {
	return _context{
		path:         path,
		flags:        flags,
		options:      options,
		namedArgs:    namedArgs,
		variadicArgs: variadicArgs,
		args:         args,
		usage:        usage,
	}
}

type _context struct {
	path         []string
	flags        map[string]bool
	options      map[string]string
	namedArgs    map[string]string
	variadicArgs []string
	args         []string
	usage        func(summaries ...string) error
}

func (c _context) Flag(name string) bool {
	return c.flags[name]
}

func (c _context) Flags() map[string]bool {
	return c.flags
}

func (c _context) Option(key string) string {
	return c.options[key]
}

func (c _context) Options() map[string]string {
	return c.options
}

func (c _context) Arg(key string) string {
	return c.namedArgs[key]
}

func (c _context) Args() []string {
	return c.args
}

func (c _context) Variadic() []string {
	return c.variadicArgs
}

func (c _context) Path() []string {
	return c.path
}

func (c _context) Usage(summaries ...string) error {
	return c.usage(summaries...)
}
