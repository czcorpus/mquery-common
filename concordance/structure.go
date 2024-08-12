package concordance

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
)

type CloseStruct struct {
	Name  string `json:"name"`
	Error error  `json:"error"`
}

func (s *CloseStruct) String() string {
	return fmt.Sprintf("</%s>", s.Name)
}

func (s *CloseStruct) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(
		struct {
			Type          string `json:"type"`
			Name          string `json:"name"`
			StructureType string `json:"structureType"`
			Error         error  `json:"error"`
		}{
			Type:          "structure",
			StructureType: "close",
			Name:          s.Name,
			Error:         s.Error,
		},
	)
}

func (s *CloseStruct) HasError() bool {
	return s.Error != nil
}

// -------

type Struct struct {
	Name  string            `json:"name"`
	Attrs map[string]string `json:"attrs"`
	// ErrMsg is an error message in case problems occured
	// with parsing related to the structure.
	ErrMsg      string `json:"errMsg,omitempty"`
	IsSelfClose bool   `json:"isSelfClose,omitempty"`
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

func (t *Struct) MarshalJSON() ([]byte, error) {
	sType := "open"
	if t.IsSelfClose {
		sType = "self-close"
	}
	return sonic.Marshal(
		struct {
			Type          string            `json:"type"`
			StructureType string            `json:"structureType"`
			Name          string            `json:"name"`
			ErrMsg        string            `json:"errMsg,omitempty"`
			Attrs         map[string]string `json:"attrs,omitempty"`
		}{
			Type:          "markup",
			StructureType: sType,
			Name:          t.Name,
			ErrMsg:        t.ErrMsg,
			Attrs:         t.Attrs,
		},
	)

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
