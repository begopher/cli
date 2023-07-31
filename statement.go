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
//	"fmt"
//
// "strings"
)

type Statement interface {
	String(path string) string
}

/*
func Summary(lines ...string) summary {
	return summary{lines}
}

type summary struct {
	lines []string
}

func (s summary) String(path string) string {
	var text strings.Builder
	if len(s.lines) == 0 {
		text.WriteString(fmt.Sprintf("Run '%s COMMAND --help' for more information on a command.\n", path))
	}else {
		for _, line := range s.lines {
			text.WriteString(line)
			text.WriteString("\n")
		}
	}
	return text.String()
}
*/
