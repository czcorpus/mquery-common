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
	"html"
	"strings"
	"unicode/utf8"
)

// LineParser parses Manatee-encoded concordance lines and converts
// them into (more structured) MQuery format.
type LineParser struct {
	attrs []string
}

func (lp *LineParser) parseTokenQuadruple(s []string) *Token {
	mAttrs := make(map[string]string)
	rawAttrs := strings.Split(s[2], "/")[1:]
	var token Token
	if len(rawAttrs) != len(lp.attrs)-1 {
		token.ErrMsg = fmt.Sprintf(
			"cannot parse token quadruple from `%s` (expected num of attrs: %d)",
			s[0], len(lp.attrs)-1)
		token.Word = s[0]
		for _, attr := range lp.attrs[1:] {
			mAttrs[attr] = "N/A"
		}

	} else {
		for i, attr := range lp.attrs[1:] {
			mAttrs[attr] = rawAttrs[i]
		}
		token.Word = s[0]
		token.Strong = len(s[1]) > 2
		token.Attrs = mAttrs
	}
	return &token
}

func (lp *LineParser) normalizeTokens(tokens []string) []string {
	ans := make([]string, 0, len(tokens))
	var parTok strings.Builder
	for _, tok := range tokens {
		tokLen := utf8.RuneCountInString(tok)
		if tok == "" {
			continue

		} else if tokLen == 1 {
			ans = append(ans, tok)

		} else if tok[0] == '{' {
			if tok[tokLen-1] != '}' {
				parTok.WriteString(tok)

			} else {
				ans = append(ans, tok)
			}

		} else if tok[tokLen-1] == '}' {
			parTok.WriteString(tok)
			ans = append(ans, parTok.String())
			parTok.Reset()

		} else {
			ans = append(ans, tok)
		}
	}
	return ans
}

func (lp *LineParser) splitToTokens(line string) []string {
	line = collIDPatt.ReplaceAllString(line, "{coll}")

	rtokens := splitPatt.Split(html.EscapeString(line), -1)
	ansTokens := make([]string, 0, len(rtokens)+5)
	for _, rtk := range rtokens {
		srch := mrgTokPatt.FindStringSubmatch(rtk)
		if len(srch) > 1 {
			ansTokens = append(ansTokens, srch[2])

		} else {
			ansTokens = append(ansTokens, rtk)
		}
	}
	return ansTokens
}

func (lp *LineParser) rmExtraColl(tokens []string) []string {
	if len(tokens)%4 == 0 {
		return tokens
	}
	ans := make([]string, 0, len(tokens))
	var prev string
	for _, tk := range tokens {
		if prev == "attr" && tk == "{coll}" {
			prev = tk
			continue
		}
		ans = append(ans, tk)
		prev = tk
	}
	return ans
}

// parseRawLine
func (lp *LineParser) parseRawLine(line string) Line {
	rtokens := lp.splitToTokens(line)
	items := lp.normalizeTokens(rtokens[1:])
	items = lp.rmExtraColl(items)
	if len(items)%4 != 0 {
		return Line{
			Text:   []*Token{{Word: "---- ERROR (unparseable) ----"}},
			Ref:    rtokens[0],
			ErrMsg: fmt.Sprintf("unparseable Manatee KWIC line: `%s`", line),
		}
	}
	tokens := make(TokenSlice, 0, len(items)/4)
	for i := 0; i < len(items); i += 4 {
		tokens = append(tokens, lp.parseTokenQuadruple(items[i:i+4]))
	}
	return Line{Text: tokens, Ref: rtokens[0]}
}

// It also escapes strings to make them usable in XML documents.
func (lp *LineParser) Parse(lines []string) []Line {
	pLines := make([]Line, len(lines))
	for i, line := range lines {
		pLines[i] = lp.parseRawLine(line)
	}
	return pLines
}

// NewLineParser is a recommended factory function
// to instantiate a `LineParser` value.
func NewLineParser(attrs []string) *LineParser {
	return &LineParser{
		attrs: attrs,
	}
}
