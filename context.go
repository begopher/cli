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

// Context gives client of cli library (developer) the ability to access all options,
// flags, arguments and variadic arguments' values, which has been passed by end user
// of the cli application.
// Context also provides access to the entire path of the executed command, which consists of:
//  - NestedApp name, followed by
//  - All Parents names if exist, followed by.
//  - Command name.
// 
// Option, Argument and Variadic methods return their values as strings. Therefore, parsing these values
// in some cases to different data types like int or date might be necessity. In such cases,
// if error occured due to invalid entry, "Context.Usage(...string) error" method can be used to
// return the usage message of the executed command, with additional information given by the developer
// to the end user, describing how the given entry can be fixed.
// "Error: " prefix should be used to follow the library convention, which make it easier and faster
// for the end user to determine where is the mistake.
// 
// # See cli.Error(context, error) function 
type Context interface {
	// Flag accepts either the short or the long name of any cli.Flag in the current executed command
	// and returns true only if end user raised that flag either by the short name or by the long name.
	Flag(string) bool
	Flags() map[string]bool
	// Option accepts either the short or the long name of any cli.Option in the current executed command
	// and returns the assotiated value which has been given by end user. Otherwise the default will be returned.
	Option(string) string
	Options() map[string]string
	// Argument accepts a name of any the cli.Argument in the current executed command and returns the correct value
	// assosiated with that name, which has been given by the end user.
	Argument(string) string
	// Variadic returns slice of string of all additional values (that came after cli.Arguments) which has been given
	// by the end user, if cli.NoVariadic is used instead of cli.Variadic then  empty slice will be returned.
	Variadic() []string
	Path() []string
	Usage(...string) error
}

func context(path []string, options map[string]string, flags map[string]bool, namedArgs map[string]string, variadicArgs []string, usage func(...string) error) _context {
	return _context{
		path:         path,
		flags:        flags,
		options:      options,
		namedArgs:    namedArgs,
		variadicArgs: variadicArgs,
		//args:         args,
		usage: usage,
	}
}

type _context struct {
	path         []string
	flags        map[string]bool
	options      map[string]string
	namedArgs    map[string]string
	variadicArgs []string
	//args         []string
	usage func(summaries ...string) error
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

func (c _context) Argument(key string) string {
	return c.namedArgs[key]
}

/*
func (c _context) Arguments() []string {
	return c.args
}
*/

func (c _context) Variadic() []string {
	return c.variadicArgs
}

func (c _context) Path() []string {
	return c.path
}

func (c _context) Usage(summaries ...string) error {
	return c.usage(summaries...)
}
