// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
//   This file is part of MQUERY.
//
//  MQUERY is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  MQUERY is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with MQUERY.  If not, see <https://www.gnu.org/licenses/>.

package concordance

import (
	"encoding/json"
	"fmt"
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
	return fmt.Sprintf("Token{Value: %s}", t.Word)
}

// ----------------------------------------------

// TokenSlice represents a flow of tokens and markup
// in a concordance line
type TokenSlice []LineElement

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
