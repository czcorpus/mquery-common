// Copyright 2024 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2024 Department of Linguistics,
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRefs(t *testing.T) {
	p := NewLineParser([]string{"word"})
	props, ref := p.parseRefs("#36940724,doc.title=Snídaně v poledne,doc.id=snidane-v-poledne")
	assert.Equal(t, "#36940724", ref)
	assert.Equal(t, 2, len(props))
	assert.Equal(t, "Snídaně v poledne", props["doc.title"])
	assert.Equal(t, "snidane-v-poledne", props["doc.id"])
}

func TestParseRefsCommaInValue(t *testing.T) {
	p := NewLineParser([]string{"word"})
	props, ref := p.parseRefs("#1,doc.title=Foo,bar, baz,doc.id=x,y")
	assert.Equal(t, "#1", ref)
	assert.Equal(t, "Foo,bar, baz", props["doc.title"])
	assert.Equal(t, "x,y", props["doc.id"])
}

func TestParseRefsNoProps(t *testing.T) {
	p := NewLineParser([]string{"word"})
	props, ref := p.parseRefs("#75308554 ")
	assert.Equal(t, "#75308554", ref)
	assert.Nil(t, props)
}
