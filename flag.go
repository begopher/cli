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

func Flag(sname string, lname, description string) flag {
	sname = strings.TrimSpace(sname)
	lname = strings.TrimSpace(lname)
	description = strings.TrimSpace(description)
	if sname == "" && lname == "" {
		panic("flag.New: sname and lname cannot be both empty")
	}
	if strings.HasPrefix(sname, "-") {
		panic("cli.Flag: - cannot be used as a short name")
	}
	if strings.HasPrefix(lname, "-") {
		panic("cli.Flag: - cannot be used as a prefix for a long name")
	}
	if sname != "" && len([]rune(sname)) > 1 {
		panic("cli.Flag: lname cannot be more than one character")
	}
	if lname != "" && len([]rune(lname)) < 2 {
		panic("cli.Flag: lname must be more than one character")
	}
	if description == "" {
		panic("cli.Flag: description cannot be empty")
	}
	return flag{
		sname:       sname,
		lname:       lname,
		description: description,
	}
}

type flag struct {
	sname       string
	lname       string
	description string
}

func (f flag) Extract(opts map[string]bool, args []string) []string {
	if len(args) < 1 {
		return args
	}
	if args[0] == "-" || args[0] == "--" || strings.HasPrefix(args[0], "---") {
		return args
	}
	if key, ok := strings.CutPrefix(args[0], "--"); ok {
		if key == f.lname {
			opts[f.lname] = true
			if f.sname != "" {
				opts[f.sname] = true
			}
			return args[1:]
		}
		return args
	}
	if key, ok := strings.CutPrefix(args[0], "-"); ok {
		if key == f.sname {
			opts[f.sname] = true
			if f.lname != "" {
				opts[f.lname] = true
			}
			return args[1:]
		}
		newKey := strings.ReplaceAll(key, f.sname, "")
		if key == newKey {
			return args
		}
		opts[f.sname] = true
		if f.lname != "" {
			opts[f.lname] = true
		}
		if newKey == "" {
			return args[1:]
		}
		args[0] = fmt.Sprintf("-%s", newKey)
	}
	return args
}

func (f flag) Default(flags map[string]bool) {
	_, ok := flags[f.lname]
	if f.lname != "" && !ok {
		flags[f.lname] = false
	}
	_, ok = flags[f.sname]
	if f.sname != "" && !ok {
		flags[f.sname] = false
	}
}

func (f flag) SName() string {
	return f.sname
}

func (f flag) LName() string {
	return f.lname
}

func (f flag) String(width int) string {
	prefix := "  "
	sflag := "    " // three spaces
	if f.sname != "" {
		sflag = "-" + f.sname + ", "
	}

	lflag := "  "
	if f.lname != "" {
		lflag = "--"
	}

	lflag = fmt.Sprintf("%s%-[2]*s", lflag, width, f.lname)
	msg := "%s%s%s  %s\n"
	return fmt.Sprintf(msg, prefix, sflag, lflag, f.description)
}
