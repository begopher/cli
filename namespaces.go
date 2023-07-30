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
	"github.com/begopher/cli/api"
)

func Namespaces(namespaces []api.Namespace) _namespaces {
	//ns := make([]api.Namespace, len(namespaces) + 1)
	//ns[0] = namespace
	//for i, namespace := range namespaces {
	//	ns[i] = namespace
	//}
	return _namespaces{
		namespaces: namespaces,
	}
}

type _namespaces struct {
	namespaces []api.Namespace
}

/*
func(n _namespaces) Has(key string) bool {
	if _, ok := n.namespace[key]; ok {
		return true
	}
	for _, namespace := range n.namespaces {
		if ok := namespace.Has(key); ok {
			return ok
		}
	}
	return false
}
*/ 
func (n _namespaces) Add(name string) error {
	for _, namespace := range n.namespaces {
		if err := namespace.Add(name); err != nil {
			return err
		}
	}
	return nil
}

func (n _namespaces) AddAll(names []string) error {
	for _, namespace := range n.namespaces {
		if err := namespace.AddAll(names); err != nil {
			return err
		}
	}
	return nil
}
/*
func (n _namespaces) AddNamespace(namespace api.Namespace) {
	n.namespaces = append(n.namespaces, namespace)
        }
*/
