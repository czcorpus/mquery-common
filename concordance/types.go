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
	"fmt"

	"github.com/bytedance/sonic"
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
	HasError() bool
}

// ------------- token and methods
// -------------------------------

// Token is a single text position in a corpus text.
type Token struct {
	Word string `json:"word"`

	// Strong is a general flag for emphasizing the token
	Strong bool `json:"strong"`

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

func (t *Token) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(
		struct {
			Type   string            `json:"type"`
			Word   string            `json:"word"`
			Strong bool              `json:"strong"`
			Attrs  map[string]string `json:"attrs"`
			ErrMsg string            `json:"errMsg,omitempty"`
		}{
			Type:   "token",
			Word:   t.Word,
			Strong: t.Strong,
			Attrs:  t.Attrs,
			ErrMsg: t.ErrMsg,
		},
	)
}

func (t *Token) String() string {
	return fmt.Sprintf("Token{Value: %s}", t.Word)
}

// ----------------------------------------------

// TokenSlice represents a flow of tokens and markup
// in a concordance line
type TokenSlice []LineElement

// Line represents a concordance line and its metadata (properties)
type Line struct {

	// Text contains positional text data (= tokens)
	Text TokenSlice `json:"text"`

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
