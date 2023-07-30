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

import(
	"fmt"
	"github.com/begopher/cli/api"
)

func Namespace() api.Namespace {
	return _namespace{make(map[string]bool)}
}

type _namespace struct {
	names map[string]bool
}

func (n _namespace) Has(key string) bool {
	if _, ok := n.names[key]; ok {
		return true
	}
	return false
}


func (n _namespace) Add(name string) error {
	if name == "" {
		return nil
	}
	if _, ok := n.names[name]; ok {
		return fmt.Errorf("%s", name)
	}
	n.names[name] = true
	return nil
}

func (n _namespace) AddAll(names []string) error {
	for _, name := range names {
		if err := n.Add(name); err != nil {
			return err
		}
	}
	return nil
}

/*
func (n _namespace) Clone(namespace api.Namespace) (duplicated string, ok bool) {
	for name, _ := range n.names {
		if namespace.Has(name) {
			return name, false
		}
		namespace.Add(name)
	}
	return "", true
}
*/
/*
func (n _namespace) Add(key string) error {
	if _, ok := n.names[key]; ok {
		return fmt.Errorf("name duplication (%s)", key)
	}
	n.names[key] = true
	return nil
}

func (n _namespace) Names() map[string]bool {
	return n.names
}

func (n _namespace) Merge(nspace api.Namespace) (api.Namespace, error) {
	names := make(map[string]bool)
	for name, _ := range n.names {
		names[name] = true
	}
	for name, _ := range nspace.Names() {
		if _, ok := names[name]; ok {
			return namespace(), fmt.Errorf("name duplication (%s)", name)
		}
		names[name] = true
	}
	return _namespace{names}, nil
}
*/
