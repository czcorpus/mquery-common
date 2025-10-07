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

// lineChunk is a partially parsed conconcrdance line.
// Typically this comes from initial parsing when we
// detect markup and normal text.
type lineChunk struct {
	value    string
	isStruct bool
}

// LineElement is a generalization of tokens and structures (markup)
// within a line
type LineElement interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
	HasError() bool
	String() string
}

// ------------------------------------------------

type MatchType string

const (
	MatchTypeKWIC MatchType = "kwic"
	MatchTypeColl MatchType = "coll"
)

// ------------- token and methods
// -------------------------------

// Token is a single text position in a corpus text.
type Token struct {
	Word string `json:"word"`

	// Strong is a general flag for emphasizing the token
	Strong bool `json:"strong"`

	// MatchType specifies how is the token related to the search
	// query. In general, we recognize "kwic" and "coll" matches here.
	MatchType MatchType `json:"matchType,omitempty"`

	// Attrs store additional attributes (e.g. PoS, lemma, syntax node parent)
	// of a respective position.
	Attrs map[string]string `json:"attrs"`

	// ErrMsg is an error message in case problems occured
	// with parsing related to the token. The policy here is
	// to always return a token with value replaced by a placeholder
	// in case of an error.
	ErrMsg string `json:"errMsg,omitempty"`
}

func (t *Token) HasError() bool {
	return t.ErrMsg != ""
}

type tokenJson struct {
	Type      string            `json:"type"`
	Word      string            `json:"word"`
	Strong    bool              `json:"strong"`
	MatchType MatchType         `json:"matchType,omitempty"`
	Attrs     map[string]string `json:"attrs"`
	ErrMsg    string            `json:"errMsg,omitempty"`
}

func (t *Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		tokenJson{
			Type:      "token",
			Word:      t.Word,
			Strong:    t.Strong,
			MatchType: t.MatchType,
			Attrs:     t.Attrs,
			ErrMsg:    t.ErrMsg,
		},
	)
}

func (t *Token) UnmarshalJSON(data []byte) error {
	var tmp tokenJson
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	t.Word = tmp.Word
	t.Strong = tmp.Strong
	t.MatchType = tmp.MatchType
	t.Attrs = tmp.Attrs
	t.ErrMsg = tmp.ErrMsg
	return nil
}

func (t *Token) String() string {
	return t.Word
}

// ----------------------------------------------

// TokenSlice represents a flow of tokens and markup
// in a concordance line
type TokenSlice []LineElement

func (ts TokenSlice) String() string {
	var ans strings.Builder
	for i, tok := range ts {
		if i > 0 {
			ans.WriteString(" ")
		}
		ans.WriteString(tok.String())
	}
	return ans.String()
}

// Tokens returns all the line elements which are tokens
// (i.e. it filters out all the structures)
func (ts TokenSlice) Tokens() []*Token {
	ans := make([]*Token, 0, len(ts))
	for _, v := range ts {
		if tv, ok := v.(*Token); ok {
			ans = append(ans, tv)
		}
	}
	return ans
}

func (ts *TokenSlice) UnmarshalJSON(data []byte) error {
	var rawElements []json.RawMessage
	if err := json.Unmarshal(data, &rawElements); err != nil {
		return err
	}

	*ts = make([]LineElement, len(rawElements))

	for i, rawElm := range rawElements {
		var typeInfo struct {
			Type          string `json:"type"`
			StructureType string `json:"structureType"`
		}
		if err := json.Unmarshal(rawElm, &typeInfo); err != nil {
			return err
		}
		var lineElm LineElement

		if typeInfo.Type == "markup" {
			if typeInfo.StructureType == "open" || typeInfo.StructureType == "self-close" {
				lineElm = &Struct{}
			}

		} else if typeInfo.Type == "token" {
			lineElm = &Token{}

		} else {
			return fmt.Errorf("unknown LineElement type %s", typeInfo.Type)
		}
		if err := json.Unmarshal(rawElm, lineElm); err != nil {
			return err
		}
		(*ts)[i] = lineElm
	}
	return nil
}

// Line represents a concordance line and its metadata (properties)
type Line struct {

	// Text contains positional text data (= tokens)
	Text TokenSlice `json:"text"`

	// AlignedText contains possible aligned text chunk in case
	// the queried corpus is a parallel one
	AlignedText TokenSlice `json:"alignedText"`

	// Ref contains numeric ID of the first token of the KWIC
	// It is typically used when referring back to the concordance
	Ref string `json:"ref"`

	// Props contains information about the text this
	// line comes from (typically information like author,
	// publication year etc.)
	Props map[string]string `json:"props,omitempty"`

	// ErrMsg is an error message in case problems occured
	// with parsing related to the line. The policy here is
	// to always return a line with value replaced by a placeholder
	// in case of an error.
	ErrMsg string `json:"errMsg,omitempty"`
}
