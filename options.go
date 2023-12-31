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

// Options represent a collection of zero or more Option.
//
// Client of cli library should not invoke any method of Options directly,
// instead, options should be passed to:
//   - cli.Cmd(..., Options, ...) function.
//   - cli.Parent(..., Options, ...) function.
//   - cli.NestedApp(..., Options, ...) function.
//
// # Panic when:
//   - one of the given option is nil value.
//   - two options has the same short or long name
func Options(opts ...api.Option) options {
	namespace := namespace()
	var width int
	for _, option := range opts {
		if option == nil {
			panic("cli.Options: nil value is not allowed")
		}
		if err := namespace.Add(option.SName()); err != nil {
			msg := fmt.Sprintf("cli.Options: %s is duplicated", option.SName())
			panic(msg)
		}
		name := option.LName()
		if err := namespace.Add(name); err != nil {
			msg := fmt.Sprintf("cli.Options: %s is duplicated", option.LName())
			panic(msg)
		}
		if width < len(name) {
			width = len(name)
		}
	}
	return options{
		opts:  opts,
		width: width,
	}
}

type options struct {
	opts  []api.Option
	width int
}

func (o options) Extract(to map[string]string, args []string) []string {
	length := len(args)
	for _, opt := range o.opts {
		args = opt.Extract(to, args)
		if length != len(args) {
			break
		}
	}
	//o.setDefault(to)
	return args

}

func (o options) Default(to map[string]string) {
	for _, opt := range o.opts {
		opt.Default(to)
	}
}

func (o options) Names() []string {
	names := make([]string, 0, len(o.opts)+len(o.opts))
	for _, option := range o.opts {
		sname := option.SName()
		if sname != "" {
			names = append(names, sname)
		}
		lname := option.LName()
		if lname != "" {
			names = append(names, lname)
		}
	}
	return names
}

func (o options) Count() int {
	return len(o.opts)
}

func (o options) String() string {
	if len(o.opts) < 1 {
		return ""
	}
	var text strings.Builder
	text.WriteString("\n")
	text.WriteString("Options:\n")
	for _, opt := range o.opts {
		text.WriteString(opt.String(o.width))
	}
	return text.String()
}

func (o options) Has(option string) bool {
	if option == "" {
		return false
	}
	if strings.HasPrefix(option, "-") {
		if name, ok := strings.CutPrefix(option, "--"); ok && len(name) > 1 {
			option = name
			goto search
		}
		if name, ok := strings.CutPrefix(option, "-"); ok && len(name) == 1 {
			option = name
			goto search
		}
		return false

	}
search:
	for _, name := range o.Names() {
		if name == option {
			return true
		}
	}
	return false
}
