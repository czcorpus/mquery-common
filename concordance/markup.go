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
	"strings"
)

var (
	splitPatt       = regexp.MustCompile(`\s+`)
	mrgTokPatt      = regexp.MustCompile(`(\{[^}]*\})([^\s]+)`)
	collIDPatt      = regexp.MustCompile(`\{col\w+(\s+col\w+)*}`)
	tagSrchRegexpSC = regexp.MustCompile(`^<([\w\d\p{Po}]+)(\s+.*?|)/>$`)
	tagSrchRegexp   = regexp.MustCompile(`^<([\w\d\p{Po}]+)(\s+.*?|)/?>$`)
	tagsAndNoTags   = regexp.MustCompile(`((<[^>]+>)+ strc)|([^<]+)`)
	splitTags       = regexp.MustCompile(`(<[^>]+>)`)
	attrValRegexp   = regexp.MustCompile(`(\w+)=([^"^\s]+)`)
	closeTagRegexp  = regexp.MustCompile(`</([^>]+)\s*>`)
	refsRegexp      = regexp.MustCompile(`((\w+\.\w+)=([^,]+))|(#\d+)`)
)

func isElement(tagSrc string) bool {
	return strings.HasPrefix(tagSrc, "<") && strings.HasSuffix(tagSrc, ">")
}

func isOpenElement(tagSrc string) bool {
	return isElement(tagSrc) && !strings.HasPrefix(tagSrc, "</") &&
		!strings.HasSuffix(tagSrc, "/>")
}

func isCloseElement(tagSrc string) bool {
	return isElement(tagSrc) && strings.HasPrefix(tagSrc, "</")
}

func isSelfCloseElement(tagSrc string) bool {
	return isElement(tagSrc) && strings.HasSuffix(tagSrc, "/>")
}
