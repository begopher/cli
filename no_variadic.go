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

func NoVariadic() noVariadic {
	return noVariadic{}
}

type noVariadic struct{}

func (v noVariadic) Arg() string {
	return ""
}

func (v noVariadic) Allowed() bool {
	return false
}

func (v noVariadic) Extract(args []string) ([]string, error) {
	if len(args) == 0 {
		return args, nil
	}
	if len(args) == 1 {
		msg := fmt.Sprintf("Error: unexpected value (%s)", args[0])
		return args, fmt.Errorf(msg)

	}
	msg := fmt.Sprintf("Error: unexpected values (%s)", strings.Join(args, ", "))
	return args, fmt.Errorf(msg)
}

func (v noVariadic) String() string {
	return ""
}
