// Copyright 2017 The kubecfg authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package registry

import (
	"encoding/json"

	"github.com/ksonnet/ksonnet/metadata/app"
)

const (
	DefaultApiVersion = "0.1"
	DefaultKind       = "ksonnet.io/registry"
)

type Spec struct {
	APIVersion string              `json:"apiVersion"`
	Kind       string              `json:"kind"`
	GitVersion *app.GitVersionSpec `json:"gitVersion"`
	Libraries  LibraryRefSpecs     `json:"libraries"`
}

func (s *Spec) Marshal() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

type Specs []*Spec

type LibraryRef struct {
	Version string `json:"version"`
	Path    string `json:"path"`
}

type LibraryRefSpecs map[string]*LibraryRef
