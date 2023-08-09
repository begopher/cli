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
)

// Option represents a command line option which may has a short and/or a long name.
// description is used when Usage is printed to the standard output so clients of
// the your command line application know the meaning of any particular option.
//
// Option may or may not have a default value. When a default value is set, it will
// be printed between parentheses right after description when Usage is printed.
//
// The value of option wheither it is the default one, or the one which has been
// given by the client of your application can be accessed by cli.Context.Option("key").
// Where key can be ethier Option's short name and/or long name.
//
// Value can also be accessed by cli.Context.Options() which returns map[string]string
// of all options with associated values.
//
// Client of cli library should not invoke any method of Option directly, instead, it should
// be passed to the cli.Options(...Option) function.
//
// # Panic when:
//   - Both short and long name are empty string. (option must have name).
//   - Short name is two or more characters.
//   - Long name is one character long.
//   - Short or long name starts with hyphen e.g. "--full-name" instead of "full-name".
//   - Description is empty. (Must tell your client the purpose of this option).
func Option(sname string, lname, description, value string) option {
	sname = strings.TrimSpace(sname)
	lname = strings.TrimSpace(lname)
	//value = strings.TrimSpace(value)
	description = strings.TrimSpace(description)
	if sname == "" && lname == "" {
		panic("cli.Option: sname and lname cannot be both empty")
	}
	if strings.HasPrefix(sname, "-") {
		panic("cli.Option: - cannot be used as a short name")
	}
	if strings.HasPrefix(lname, "-") {
		panic("cli.Option: - cannot be used as a prefix for a long name")
	}
	if sname != "" && len([]rune(sname)) > 1 {
		panic("cli.Option: lname cannot be more than one character")
	}
	if lname != "" && len([]rune(lname)) < 2 {
		panic("cli.Option: lname must be more than one character")
	}
	if description == "" {
		panic("cli.Option: description cannot be empty")
	}
	//if value == "" {
	//	panic("cli.Option: default value cannot be empty")
	//}
	return option{
		sname:       sname,
		lname:       lname,
		description: description,
		value:       value,
	}
}

type option struct {
	sname       string
	lname       string
	description string
	value       string
}

func (o option) Extract(opts map[string]string, args []string) []string {
	if len(args) < 2 {
		return args
	}
	if args[0] == "-" || args[0] == "--" || strings.HasPrefix(args[0], "---") {
		return args
	}
	if key, ok := strings.CutPrefix(args[0], "--"); ok && key == o.lname {
		opts[o.lname] = args[1]
		if o.sname != "" {
			opts[o.sname] = args[1]
		}
		return args[2:]
	}
	if key, ok := strings.CutPrefix(args[0], "-"); ok && key == o.sname {
		opts[o.sname] = args[1]
		if o.lname != "" {
			opts[o.lname] = args[1]
		}
		return args[2:]
	}
	return args
}

func (o option) Default(opts map[string]string) {
	_, ok := opts[o.lname]
	if !ok && o.lname != "" {
		opts[o.lname] = o.value
	}
	//if value := opts[o.lname]; value != "" && !ok {
	//	opts[o.lname] = o.value
	//}
	_, ok = opts[o.sname]
	if !ok && o.sname != "" {
		opts[o.sname] = o.value
	}
	//if value := opts[o.sname]; value != "" && !ok {
	//	opts[o.sname] = o.value
	//}
}

func (o option) SName() string {
	return o.sname
}

func (o option) LName() string {
	return o.lname
}

func (o option) String(width int) string {
	prefix := "  "
	sflag := "    "
	if o.sname != "" {
		sflag = "-" + o.sname + ", "
	}

	lflag := "  "
	if o.lname != "" {
		lflag = "--"
	}
	lflag = fmt.Sprintf("%s%-[2]*s", lflag, width, o.lname)
	var def string = ``
	if o.value != "" {
		def = fmt.Sprintf(`(%s)`, o.value)
	}
	msg := "%s%s%s  %s %s\n"
	return fmt.Sprintf(msg, prefix, sflag, lflag, o.description, def)
}
