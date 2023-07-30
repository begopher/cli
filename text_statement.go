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
	"strings"
)

func TextStatement(lines ...string) textStatement {
	if len(lines) == 0 {
		panic("cli.TextStatement: lines argument cannot be empty")
	}
	return textStatement{lines}
}

type textStatement struct {
	lines []string
}

func (s textStatement) String(string) string {
	var text strings.Builder
	for _, line := range s.lines {
		text.WriteString(line)
		text.WriteString("\n")
	}
	return text.String()
}
