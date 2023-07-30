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

func Statements(statements ...Statement) _statements {
	if len(statements) == 0 {
		panic("cli.Statements: statements cannot be empty")
	}
	for _, stat := range statements {
		if stat == nil {
			panic("cli.Statements: nil value is not accepted")
		}
	}
	return _statements{statements}
}

type _statements struct {
	statements []Statement
}

func (s _statements) String(path string) string {
	var text strings.Builder
	for _, statement := range s.statements {
		text.WriteString(statement.String(path))
	}
	return text.String()
}

