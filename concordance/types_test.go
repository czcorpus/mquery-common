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
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestTokenSliceMarkupJSONRoundTrip(t *testing.T) {
	src := TokenSlice{
		&Struct{Name: "s", Attrs: map[string]string{"id": "s1"}},
		&Token{Word: "pes", Attrs: map[string]string{"lemma": "pes"}},
		&Struct{Name: "lb", IsSelfClose: true},
		&Token{Word: "štěká", Attrs: map[string]string{"lemma": "štěkat"}},
		&CloseStruct{Name: "s"},
	}
	data, err := json.Marshal(src)
	require.NoError(t, err)
	var ans TokenSlice
	require.NoError(t, json.Unmarshal(data, &ans))
	require.Len(t, ans, 5)
	assert.Equal(t, src[0], ans[0])
	assert.Equal(t, "pes", ans[1].(*Token).Word)
	assert.Equal(t, &Struct{Name: "lb", IsSelfClose: true}, ans[2])
	assert.Equal(t, &CloseStruct{Name: "s"}, ans[4])
	assert.Equal(t, "<s id=s1> pes <lb /> štěká </s>", ans.String())
}

func TestCloseStructErrorJSON(t *testing.T) {
	data, err := json.Marshal(&CloseStruct{Name: "s", Error: errors.New("unmatched")})
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"markup","structureType":"close","name":"s","error":"unmatched"}`, string(data))
	var ans CloseStruct
	require.NoError(t, json.Unmarshal(data, &ans))
	assert.Equal(t, "s", ans.Name)
	require.True(t, ans.HasError())
	assert.Equal(t, "unmatched", ans.Error.Error())
}

func TestTokenSliceUnknownMarkup(t *testing.T) {
	var ans TokenSlice
	err := json.Unmarshal([]byte(`[{"type":"markup","structureType":"foo","name":"s"}]`), &ans)
	assert.Error(t, err)
}
