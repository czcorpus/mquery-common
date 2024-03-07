// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Martin Zimandl <martin.zimandl@gmail.com>
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

const (
	ts1 = `#75308554 ` +
		`která {} /který/zavádět/+1 attr  zavádí {} /zavádět/země/-5 attr  ` +
		`celoplošný {} /celoplošný/provoz/+1 attr provoz {col0 coll} /provoz/zavádět/-2 attr ` +
		`těchto {} /tento/služba/+1 attr  služeb {} /služba/provoz/-2 attr  . {} /.//0 attr`
)

func TestExampleLines(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "p_lemma", "parent"})
	ans := p.Parse([]string{ts1})
	assert.Equal(t, "", ans[0].ErrMsg)
	assert.Equal(t, "#75308554", ans[0].Ref)
	assert.Equal(t, "která", ans[0].Text[0].Word)
	assert.Equal(t, "který", ans[0].Text[0].Attrs["lemma"])
	assert.Equal(t, "zavádět", ans[0].Text[0].Attrs["p_lemma"])
	assert.Equal(t, "+1", ans[0].Text[0].Attrs["parent"])

	assert.Equal(t, 7, len(ans[0].Text))

	assert.Equal(t, "#75308554", ans[0].Ref)
	assert.Equal(t, ".", ans[0].Text[6].Word)
	assert.Equal(t, ".", ans[0].Text[6].Attrs["lemma"])
	assert.Equal(t, "", ans[0].Text[6].Attrs["p_lemma"])
	assert.Equal(t, "0", ans[0].Text[6].Attrs["parent"])
}
