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

func Flags(flgs ...api.Flag) flags {
	namespace := namespace()
	var width int
	for _, flag := range flgs {
		if flag == nil {
			panic("cli.Flags: nil value is not allowed")
		}
		if err := namespace.Add(flag.SName()); err != nil {
			msg := fmt.Sprintf("cli.Flags: %s is duplicated", flag.SName())
			panic(msg)
		}
		name := flag.LName()
		if err := namespace.Add(name); err != nil {
			msg := fmt.Sprintf("cli.Flags: %s is duplicated", flag.LName())
			panic(msg)
		}
		if width < len(name) {
			width = len(name)
		}
	}
	return flags{
		flgs: flgs,
		width: width,
	}
}

type flags struct {
	flgs []api.Flag
	width int
}

func (f flags) Extract(to map[string]bool, args []string) []string {
	args = f.recursive(to, args)
	return args

}

func (f flags) Default(to map[string]bool) {
	for _, flag := range f.flgs {
		flag.Default(to)
	}
}

func (f flags) recursive(to map[string]bool, args []string) []string {
	length := len(args)
	for _, flag := range f.flgs {
		args = flag.Extract(to, args)
		if length != len(args) {
			break
		}
	}
	return args
}

func (f flags) Names() []string {
	names := make([]string, 0, len(f.flgs)+len(f.flgs))
	for _, flag := range f.flgs {
		sname := flag.SName()
		if sname != "" {
			names = append(names, sname)
		}
		lname := flag.LName()
		if lname != "" {
			names = append(names, lname)
		}
	}
	return names
}

func (f flags) Count() int {
	return len(f.flgs)
}

func (f flags) String() string {
	var text strings.Builder
	if len(f.flgs) < 1 {
		return ""
	}
	text.WriteString("\n")
	text.WriteString("Flags:\n")
	for _, flag := range f.flgs {
		text.WriteString(flag.String(f.width))
	}
	return text.String()
}
