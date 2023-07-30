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

func Options(opts ...api.Option) options {
	namespace := Namespace()
	for _, option := range opts {
		if option == nil {
			panic("cli.Options: nil value is not allowed")
		}
		if err := namespace.Add(option.SName()); err != nil {
			msg := fmt.Sprintf("cli.Options: %s is duplicated", option.SName())
			panic(msg)
		}
		if err := namespace.Add(option.LName()); err != nil {
			msg := fmt.Sprintf("cli.Options: %s is duplicated", option.LName())
			panic(msg)
		}
	}
	return options{
		opts: opts,
	}
}

type options struct {
	opts []api.Option
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

func (o options) LNameWidth() int {
	var width int
	for _, option := range o.opts {
		name := option.LName()
		if width < len(name) {
			width = len(name)
		}
	}
	return width
}

func (o options) String(width int) string {
	if len(o.opts) < 1 {
		return ""
	}
	var text strings.Builder
	text.WriteString("Options:\n")
	for _, opt := range o.opts {
		text.WriteString(opt.String(width))
	}
	text.WriteString("\n")
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
