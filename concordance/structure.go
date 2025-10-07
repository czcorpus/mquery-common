// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Department of Linguistics,
//                Faculty of Arts, Charles University
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package concordance

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CloseStruct struct {
	Name  string
	Error error
}

func (s *CloseStruct) String() string {
	return fmt.Sprintf("</%s>", s.Name)
}

type closeStructJson struct {
	Type          string `json:"type"`
	StructureType string `json:"structureType"`
	Name          string `json:"name"`
	Error         error  `json:"error,omitempty"`
}

func (s *CloseStruct) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		closeStructJson{
			Type:          "markup",
			StructureType: "close",
			Name:          s.Name,
			Error:         s.Error,
		},
	)
}

func (s *CloseStruct) UnmarshalJSON(data []byte) error {
	var tmp closeStructJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	s.Name = tmp.Name
	s.Error = tmp.Error
	return nil
}

func (s *CloseStruct) HasError() bool {
	return s.Error != nil
}

// -------

type Struct struct {
	Name  string
	Attrs map[string]string
	// ErrMsg is an error message in case problems occured
	// with parsing related to the structure.
	ErrMsg      string
	IsSelfClose bool
}

func (t *Struct) String() string {
	var ans strings.Builder
	ans.WriteString("<" + t.Name)
	for k, v := range t.Attrs {
		ans.WriteString(" " + k + "=" + v)
	}
	if t.IsSelfClose {
		ans.WriteString(" />")

	} else {
		ans.WriteString(">")
	}
	return ans.String()
}

func (t *Struct) HasError() bool {
	return t.ErrMsg != ""
}

type structJson struct {
	Type          string            `json:"type"`
	StructureType string            `json:"structureType"`
	Name          string            `json:"name"`
	ErrMsg        string            `json:"error,omitempty"`
	Attrs         map[string]string `json:"attrs,omitempty"`
}

func (t *Struct) MarshalJSON() ([]byte, error) {
	sType := "open"
	if t.IsSelfClose {
		sType = "self-close"
	}
	return json.Marshal(
		structJson{
			Type:          "markup",
			StructureType: sType,
			Name:          t.Name,
			ErrMsg:        t.ErrMsg,
			Attrs:         t.Attrs,
		},
	)
}

func (t *Struct) UnmarshalJSON(data []byte) error {
	var tmp structJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	t.Name = tmp.Name
	t.Attrs = tmp.Attrs
	t.ErrMsg = tmp.ErrMsg
	if tmp.Type == "self-close" {
		t.IsSelfClose = true
	}
	return nil
}

func parseStructure(src string) LineElement {
	if isSelfCloseElement(src) {
		values := tagSrchRegexpSC.FindStringSubmatch(src)
		if len(values) > 0 {
			attrSrch := attrValRegexp.FindAllStringSubmatch(values[2], -1)
			attrs := make(map[string]string)
			for _, a := range attrSrch {
				attrs[a[1]] = a[2]
			}
			return &Struct{
				IsSelfClose: true,
				Name:        values[1],
				Attrs:       attrs,
			}
		}

	} else if isOpenElement(src) {
		values := tagSrchRegexp.FindStringSubmatch(src)
		if len(values) > 0 {
			attrSrch := attrValRegexp.FindAllStringSubmatch(values[2], -1)
			attrs := make(map[string]string)
			for _, a := range attrSrch {
				attrs[a[1]] = a[2]
			}
			return &Struct{
				Name:  values[1],
				Attrs: attrs,
			}
		}

	} else if isCloseElement(src) {
		srch := closeTagRegexp.FindStringSubmatch(src)
		if len(srch) > 0 {
			return &CloseStruct{
				Name: srch[1],
			}
		}
	}
	return &Struct{}
}
