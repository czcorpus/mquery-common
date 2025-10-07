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
