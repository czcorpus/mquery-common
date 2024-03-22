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
	"regexp"
)

var (
	splitPatt  = regexp.MustCompile(`\s+`)
	mrgTokPatt = regexp.MustCompile(`(\{[^}]*\})([^\s]+)`)
	collIDPatt = regexp.MustCompile(`\{col\w+(\s+col\w+)*}`)
)

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

type TokenSlice []*Token

type Line struct {

	// Text contains positional text data (= tokens)
	Text TokenSlice `json:"text"`

	// Ref contains structural metadata related to the line
	Ref string `json:"ref"`

	// ErrMsg is an error message in case problems occured
	// with parsing related to the line. The policy here is
	// to always return a line with value replaced by a placeholder
	// in case of an error.
	ErrMsg string `json:"errMsg,omitempty"`
}
