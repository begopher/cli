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

func Variadic(arg, description string) variadic {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		panic("cli.Variadic: arg cannot be empty, look at cli.NoVariadic()")
	}
	description = strings.TrimSpace(description)
	return variadic{
		arg:         arg,
		description: description,
	}
}

type variadic struct {
	arg         string
	description string
}

func (v variadic) Arg() string {
	return fmt.Sprintf("[%s]", v.arg)
}

func (v variadic) Allowed() bool {
	return true
}

func (v variadic) Extract(args []string) ([]string, error) {
	return args, nil
}

func (v variadic) String() string {
	if v.description == "" {
		return ""
	}
	var text strings.Builder
	text.WriteString("Variadic:\n")
	msg := fmt.Sprintf("  [%s]  %s\n", v.arg, v.description)
	text.WriteString(msg)
	text.WriteString("\n")
	return text.String()
}
