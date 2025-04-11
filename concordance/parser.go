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
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	// RefsEndMark is a custom separator which is used by Mquery
	// to separate the "refs" section of a line output to
	RefsEndMark = "{refs:end}"
)

var (
	CollColl1Srch = regexp.MustCompile(`{}|{col.?( col\w+)*}|attr`)
)

// LineParser parses Manatee-encoded concordance lines and converts
// them into (more structured) MQuery format.
type LineParser struct {
	attrs []string
}

func (lp *LineParser) parseTokenQuadruple(s []string) *Token {
	mAttrs := make(map[string]string)
	attrString := s[2]
	delimiter := attrString[:1] // we can use such value access as delim. is never > 1 byte
	rawAttrs := strings.Split(attrString, delimiter)[1:]
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
		if s[1] == "{kwic}" {
			token.MatchType = MatchTypeKWIC

		} else if s[1] == "{coll}" {
			token.MatchType = MatchTypeColl
		}
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

func (lp *LineParser) splitToTokens(line string) ([]string, string) {
	line = strings.ReplaceAll(line, "{col0 coll}", "{kwic}")
	line = strings.ReplaceAll(line, "{coll coll1}", "{coll}")
	line = collIDPatt.ReplaceAllString(line, "{coll}") // transform possible unknown stuff to coll

	refsAndRest := strings.Split(line, RefsEndMark)
	var refsText string
	if len(refsAndRest) > 1 {
		refsText = refsAndRest[0]
		line = strings.Join(refsAndRest[1:], " ")
	}
	rtokens := splitPatt.Split(line, -1)
	ansTokens := make([]string, 0, len(rtokens)+5)
	for _, rtk := range rtokens {
		srch := mrgTokPatt.FindStringSubmatch(rtk)
		if len(srch) > 1 {
			ansTokens = append(ansTokens, srch[2])

		} else {
			ansTokens = append(ansTokens, rtk)
		}
	}
	return ansTokens, refsText
}

// extractStructures is a first stage parsing of Manatee concordance output which
// isolates text and structural (markup) chunks.
func (lp *LineParser) extractStructures(line string) []lineChunk {
	chunks := tagsAndNoTags.FindAllString(line, -1)
	ans := make([]lineChunk, len(chunks))
	for i, ch := range chunks {
		if strings.HasPrefix(ch, "<") && strings.HasSuffix(ch, "strc") {
			ans[i] = lineChunk{value: ch, isStruct: true}

		} else {
			ans[i] = lineChunk{value: ch}
		}
	}
	return ans
}

// parseRefs parses text metadata (aka "refs" in KonText/NoSkE)
func (lp *LineParser) parseRefs(refs string) (ans map[string]string, ref string) {
	srch := refsRegexp.FindAllStringSubmatch(refs, -1)
	for _, item := range srch {
		if strings.HasPrefix(item[0], "#") {
			ref = item[0]

		} else {
			if ans == nil {
				ans = make(map[string]string)
			}
			ans[item[2]] = item[3]
		}
	}
	return
}

// normalizeCurlyMarkup solves the situation when the simplest pattern:
//
// `foo {} SEPfoo_attr2SEPfoo_attr3SEP...foo_attrN attr`
//
// is replaced by:
//
// {col0 coll ... etc. } foo {coll ...etc...} SEPfoo_attr2SEPfoo_attr3SEP...foo_attrN attr
//
// I.e. when a kwic/collocate is enclosed in curly braced markup. In such case
// we want to transform the string to the simplest form so we can easily parse it.
func (lp *LineParser) normalizeCurlyMarkup(s string) string {
	// note - it this method, we use indexing within
	// a string to cut pieces which is normally a bad
	// idea as the indexes are pointing to bytes and
	// not utf8 runes. But we take advantage of the fact,
	// that regexp's FindAllStringIndex returns byte-aware
	// indexing.
	srch := CollColl1Srch.FindAllStringIndex(s, -1)
	pos1 := [2]int{-1, -1} // position of latest '{}'
	lastPos := 0           // last position of the cut once we encounter '{col* col...}'
	numCurlyItems := 0
	var ans strings.Builder
	for _, x := range srch {
		token := s[x[0]:x[1]]
		if strings.Contains(token, "attr") {
			numCurlyItems = 0
			ans.WriteString(s[lastPos:x[1]] + " ")
			lastPos = x[1]

		} else {
			numCurlyItems++
			if numCurlyItems == 1 {
				pos1 = [2]int{x[0], x[1]}

			} else if numCurlyItems > 1 {
				if pos1[0]-1 > lastPos {
					ans.WriteString(s[lastPos : pos1[0]-1])
				}
				ans.WriteString(s[pos1[1]:x[1]] + " ")
				lastPos = x[1] + 1
			}
		}
	}
	if lastPos < len(s)-1 {
		ans.WriteString(s[lastPos:])
	}
	return ans.String()
}

// parseRawLine
func (lp *LineParser) parseRawLine(rawLine string) Line {
	rawLine = lp.normalizeCurlyMarkup(rawLine)
	chunks := lp.extractStructures(rawLine)
	line := Line{}
	for i, chunk := range chunks {
		if chunk.isStruct {
			multiStructSrch := splitTags.FindAllStringSubmatch(chunk.value, -1)
			for _, item := range multiStructSrch {
				line.Text = append(line.Text, parseStructure(item[1]))
			}

		} else {
			rtokens, refs := lp.splitToTokens(chunk.value)
			if i == 0 {
				line.Props, line.Ref = lp.parseRefs(refs)
			}
			items := lp.normalizeTokens(rtokens)
			if len(items)%4 != 0 {
				line.Text = append(line.Text, &Token{Word: "---- ERROR (unparseable) ----"})
				line.ErrMsg = fmt.Sprintf(
					"unparseable Manatee KWIC line: expected 4N elms, found %d: `%s`",
					len(items),
					chunk.value,
				)

			} else {
				for i := 0; i < len(items); i += 4 {
					line.Text = append(line.Text, lp.parseTokenQuadruple(items[i:i+4]))
				}
			}
		}
	}
	return line
}

// Parse parses custom Manatee-open concordance output format into
// a structured one. The `attrDelim` specifies ASCII character used
// as a separator between individual positional attributes. Manatee-open
// uses "/" as a hardcoded value but our (CNC) forks may provide
// more suitable selection (e.g. the "unit separator" character) so it
// does not collide with text itself.
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
