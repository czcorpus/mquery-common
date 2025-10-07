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

package corp

import (
	"strings"

	"github.com/czcorpus/cnc-gokit/collections"
)

// PosAttr specifies a configuration of a positional attribute (= token).
type PosAttr struct {
	Name string `json:"name"`

	// Description contains localized description of the attribute.
	Description map[string]string `json:"description"`
}

func (p PosAttr) IsZero() bool {
	return p.Name == ""
}

// LocaleDescription returns a localized description
// of the attribute. In case the `lang` is not present,
// "en" version is returned.
func (p PosAttr) LocaleDescription(lang string) string {
	d := p.Description[lang]
	if d != "" {
		return d
	}
	return p.Description["en"]
}

// --------------------

// PosAttrList is a slice of PosAttr with some methods
// added for easier use.
type PosAttrList []PosAttr

func (pal PosAttrList) GetIDs() []string {
	ans := make([]string, len(pal))
	for i, v := range pal {
		ans[i] = v.Name
	}
	return ans
}

func (pal PosAttrList) Contains(ident string) bool {
	for _, v := range pal {
		if v.Name == ident {
			return true
		}
	}
	return false
}

// ------

type SyntaxConcordance struct {
	ParentAttr string `json:"parentAttr"`

	// ResultAttrs is a list of positional attributes
	// we need to provide all the required information about
	// syntax in for the "syntax-conc-examples" endpoint
	ResultAttrs []string `json:"resultAttrs"`
}

// ------

// ---------------

type CorpusVariant struct {
	ID          string            `json:"id"`
	FullName    map[string]string `json:"fullName"`
	Description map[string]string `json:"description"`
}

// ---------------

type ContextWindow int

// LeftAndRight converts the window into the left-right intervals.
// In case the value is an odd number, the reminder is added to the right.
func (cw ContextWindow) LeftAndRight() (lft int, rgt int) {
	tmp := int(cw) / 2
	lft = tmp
	rgt = tmp + (int(cw) - 2*tmp)
	return
}

// ---------------

// CorpusSetup is a general configuration of a corpus in MQuery and other apps.
type CorpusSetup struct {
	ID                   string             `json:"id"`
	FullName             map[string]string  `json:"fullName"`
	Description          map[string]string  `json:"description"`
	SyntaxConcordance    SyntaxConcordance  `json:"syntaxConcordance"`
	PosAttrs             PosAttrList        `json:"posAttrs"`
	ConcMarkupStructures []string           `json:"concMarkupStructures"`
	ConcTextPropsAttrs   []string           `json:"concTextPropsAttrs"`
	TextProperties       TextTypeProperties `json:"textProperties"`
	MaximumRecords       int                `json:"maximumRecords"`

	// MaximumTokenContextWindow specifies the total width of token's context
	// with the token in the middle. Odd numbers are applied in a way giving one
	// more token to the right.
	MaximumTokenContextWindow ContextWindow `json:"MaximumTokenContextWindow"`

	// Subcorpora defines named transient subcorpora created as part of the query.
	// MQuery also supports so called saved subcorpora which are files created via Manatee-open
	// (or in a more user-friendly way using KonText or NoSkE).
	Subcorpora map[string]Subcorpus `json:"subcorpora"`

	// ViewContextStruct is a structure used to specify "units"
	// for KWIC left and right context. Typically, this is
	// a structure representing a sentence or a speach.
	ViewContextStruct string `json:"viewContextStruct"`

	// Variants allows for specifying multiple corpora based on a common "core".
	// This typicallay applies for parallel corpora where their structure is the same
	// and they differ just in a language.
	Variants       map[string]CorpusVariant `json:"variants"`
	SrchKeywords   []string                 `json:"srchKeywords"`
	WebURL         string                   `json:"webUrl"`
	HasPublicAudio bool                     `json:"hasPublicAudio"`

	// BibLabelAttr is an attribute specifying a (possibly non-unique) title/label of an
	// unique text work (e.g. a book). It always comes with BibIDAttr.
	//
	// Note that not all corpora have to have this attribute specified.
	BibLabelAttr string `json:"bibLabelAttr"`

	// BibIDAttr is an attribute specifying a unique text work by providing a unique ID.
	// It always comes with BibLabelAttr.
	//
	// Note that not all corpora have to have this attribute specified.
	BibIDAttr string `json:"bibIdAttr"`

	Tagsets []SupportedTagset `json:"tagsets"`

	// Size represents size of corpus in tokens. In MQuery, this does not
	// have to be configured as MQuery can get the value via Manatee.
	Size int64 `json:"size,omitempty"`
}

func (cs *CorpusSetup) LocaleDescription(lang string) string {
	d := cs.Description[lang]
	if d != "" {
		return d
	}
	return cs.Description["en"]
}

func (cs *CorpusSetup) IsDynamic() bool {
	return strings.Contains(cs.ID, "*")
}

func (cs *CorpusSetup) GetPosAttr(name string) PosAttr {
	for _, v := range cs.PosAttrs {
		if v.Name == name {
			return v
		}
	}
	return PosAttr{}
}

func (cs *CorpusSetup) KnownStructures() []string {
	ans := make([]string, 0, len(cs.ConcMarkupStructures)+len(cs.ConcTextPropsAttrs))
	ans = append(ans, cs.ConcMarkupStructures...)
	tmp := collections.SliceMap(cs.ConcTextPropsAttrs, func(s string, i int) string {
		return strings.Split(s, ".")[0]
	})
	ans = append(ans, tmp...)
	return ans
}
