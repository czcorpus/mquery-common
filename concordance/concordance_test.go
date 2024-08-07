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

	ts2 = `#108182398 ` +
		`. {} /./Z:------------- attr  ?? {} /??/Z:------------- attr  KDYŽ {} /když/J,------------- attr` +
		`   {}VEJCE {col0 coll coll coll1} /vejce/NNNS1-----A---- attr  K {col0 coll} /k/RR--3---------- attr` +
		`   {col0 coll} VEJCI {col0 coll coll coll2} /vejce/NNNS3-----A---- attr` +
		`  SEDÁ {col0 coll} /sedat/VB-S---3P-AAI-- attr Z {} /z/RR--2---------- attr` +
		`  váz {} /váza/NNFP2-----A---- attr  a {} /a/J^------------- attr'`

	ts3_struct = `#61705575 ` +
		`pasti {} /past/NNFS2-----A---- attr <g foo=bar /> strc . {} /./Z:------------- attr </hi></s><s id=picko_knihaofyzi:1:1144:4><hi> strc` +
		` 1982 {} /1982/C=------------- attr  / {} ///Z:------------- attr <g/> strc / {} ///Z:------------- attr <g/> strc Kvazikrystaly` +
		` {col0 coll} /kvazikrystal/NNIP1-----A---- attr</s><s id=picko_knihaofyzi:1:1145:1 strong=true> strc Na {} /na/RR--4---------- attr  exotické` +
		` {} /exotický/AAIP4----1A---- attr  kvazikrystaly {} /kvazikrystal/NNIP4-----A---- attr  si {} /se/P7--3---------- attr  často {}` +
		` /často/Dg-------1A---- attr`
)

func asTokenOrPanic(v LineElement) *Token {
	tmp, ok := v.(*Token)
	if !ok {
		panic("not a token")
	}
	return tmp
}

func TestExampleLines(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "p_lemma", "parent"})
	ans := p.Parse([]string{ts1})
	assert.Equal(t, "", ans[0].ErrMsg)
	assert.Equal(t, "#75308554", ans[0].Ref)
	tok := asTokenOrPanic(ans[0].Text[0])
	assert.Equal(t, "která", tok.Word)
	assert.Equal(t, "který", tok.Attrs["lemma"])
	assert.Equal(t, "zavádět", tok.Attrs["p_lemma"])
	assert.Equal(t, "+1", tok.Attrs["parent"])

	assert.Equal(t, 7, len(ans[0].Text))

	assert.Equal(t, "#75308554", ans[0].Ref)
	tok = asTokenOrPanic(ans[0].Text[6])
	assert.Equal(t, ".", tok.Word)
	assert.Equal(t, ".", tok.Attrs["lemma"])
	assert.Equal(t, "", tok.Attrs["p_lemma"])
	assert.Equal(t, "0", tok.Attrs["parent"])
}

func TestRegression1(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.Parse([]string{ts2})
	for _, a := range ans {
		assert.Zero(t, a.ErrMsg)
	}
}

func TestParsingLineWithStructs(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.Parse([]string{ts3_struct})
	for _, a := range ans {
		assert.NotZero(t, a.ErrMsg)
	}
}
