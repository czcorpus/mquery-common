// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Martin Zimandl <martin.zimandl@gmail.com>
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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ts1 = `#75308554 ` + RefsEndMark + " " +
		`která {} /který/zavádět/+1 attr  zavádí {} /zavádět/země/-5 attr  ` +
		`celoplošný {} /celoplošný/provoz/+1 attr provoz {col0 coll} /provoz/zavádět/-2 attr ` +
		`těchto {} /tento/služba/+1 attr  služeb {} /služba/provoz/-2 attr  . {} /.//0 attr`

	ts2 = `#108182398 ` + RefsEndMark + " " +
		`. {} /./Z:------------- attr  ?? {} /??/Z:------------- attr  KDYŽ {} /když/J,------------- attr` +
		`   {}VEJCE {col0 coll coll coll1} /vejce/NNNS1-----A---- attr  K {col0 coll} /k/RR--3---------- attr` +
		`   {col0 coll} VEJCI {col0 coll coll coll2} /vejce/NNNS3-----A---- attr` +
		`  SEDÁ {col0 coll} /sedat/VB-S---3P-AAI-- attr Z {} /z/RR--2---------- attr` +
		`  váz {} /váza/NNFP2-----A---- attr  a {} /a/J^------------- attr'`

	ts3_struct = `#61705575 ` + RefsEndMark + " " +
		`pasti {} /past/NNFS2-----A---- attr <g foo=bar /> strc . {} /./Z:------------- attr </hi></s><s id=picko_knihaofyzi:1:1144:4><hi> strc` +
		` 1982 {} /1982/C=------------- attr  / {} ///Z:------------- attr <g/> strc / {} ///Z:------------- attr <g/> strc Kvazikrystaly` +
		` {col0 coll} /kvazikrystal/NNIP1-----A---- attr</s><s id=picko_knihaofyzi:1:1145:1 strong=true> strc Na {} /na/RR--4---------- attr  exotické` +
		` {} /exotický/AAIP4----1A---- attr  kvazikrystaly {} /kvazikrystal/NNIP4-----A---- attr  si {} /se/P7--3---------- attr  často {}` +
		` /často/Dg-------1A---- attr`

	ts4_coll = `#57713857{refs:end}kvalita {} \u001fkvalita\u001fNNFS1-----A---- attr  života {} \u001fživot\u001fNNIS2-----A---- attr  pomáhajících {} \u001fpomáhající\u001fAGFP2-----A---- attr  profesí {} \u001fprofese\u001fNNFP2-----A---- attr  , {} \u001f,\u001fZ:------------- attr kvalita {col0 coll} \u001fkvalita\u001fNNFS1-----A---- attr  {} života {coll coll1} \u001fživot\u001fNNIS2-----A---- attr  rodinných {} \u001frodinný\u001fAAMP2----1A---- attr  příslušníků {} \u001fpříslušník\u001fNNMP2-----A---- attr  pečujících {} \u001fpečující\u001fAGMP2-----A---- attr  o {} \u001fo\u001fRR--4---------- attr`
	ts5_coll = `#37219897{refs:end}smysl {} \x1fsmysl\x1fNNIS4-----A---- attr  . {} \x1f.\x1fZ:------------- attr  Teplý {} \x1fteplý\x1fAAIS1----1A---- attr  tábor {} \x1ftábor\x1fNNIS1-----A---- attr  je {} \x1fbýt\x1fVB-S---3P-AAI-- attr dojná {col0 coll} \x1fdojný\x1fAAFS1----1A---- attr  {} kráva {coll coll1} \x1fkráva\x1fNNFS1-----A---- attr  . {} \x1f.\x1fZ:------------- attr  Všem {} \x1fvšechen\x1fPLMP3---------- attr  , {} \x1f,\x1fZ:------------- attr  kdo {} \x1fkdo\x1fPKM-1---------- attr`

	ts6_refs = `#40281069,doc.title=Pastička,doc.txtype=SCR: drama{refs:end}ANNA {} AnnaNNFS1-----A---- attr  - {} -Z:------------- attr  Tak {} takDb------------- attr ahoj {col0 coll} ahojII------------- attr . {} .Z:------------- attr`

	ts7_refs = `#1246234523,doc.title=Captain Brand of the 'Centipede'<BR>A Pirate of Eminence in the West Indies: His Love and Exploits, Together with Some Account of the Singular Manner` + RefsEndMark + ` pasti {} /past/NNFS2-----A---- attr`

	// note: the last token uses \x1f as attr. delimiter as its tag contains '/'
	ts8_markup_in_attr = `#1` + RefsEndMark + ` foo {} /foo/jj_<sub> attr <s id=1> strc bar {} /bar/a<b attr` +
		` <P> {} /-/null attr` +
		" baz {} \x1fbaz\x1f<i>NN</i> attr</s> strc"

	ts9_attr_const_sens = `#1` + RefsEndMark + `succumbed {} /succumb/vvn attr to {} /to/ii attr the {} /the/at attr gayer {col0 coll} /gay/jjr attr attractions {} /attraction/nn2 attr of {} /of/io attr Parky {} /parky/np1 attr`
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
		assert.Zero(t, a.ErrMsg)
		structList := make([]*Struct, 0, 5)
		closesList := make([]*CloseStruct, 0, 5)
		for _, item := range a.Text {
			switch tItem := item.(type) {
			case *Struct:
				structList = append(structList, tItem)
			case *CloseStruct:
				closesList = append(closesList, tItem)
			}
		}
		assert.Len(t, structList, 6)
		assert.Equal(t, "g", structList[0].Name)
		assert.True(t, structList[0].IsSelfClose)
		assert.Equal(t, "bar", structList[0].Attrs["foo"])
		assert.Equal(t, "s", structList[5].Name)
		assert.Equal(t, "picko_knihaofyzi:1:1145:1", structList[5].Attrs["id"])
		assert.Equal(t, "true", structList[5].Attrs["strong"])
		assert.Len(t, closesList, 3)
		assert.Equal(t, "hi", closesList[0].Name)
		assert.Equal(t, "s", closesList[1].Name)
		assert.Equal(t, "s", closesList[2].Name)
	}
}

func TestRegression4(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.Parse([]string{ts4_coll})
	for _, a := range ans {
		assert.Zero(t, a.ErrMsg)
	}
}

func TestRegression5(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.Parse([]string{ts5_coll})
	for _, a := range ans {
		assert.Zero(t, a.ErrMsg)
	}
}

func TestRefsExtraction(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.Parse([]string{ts6_refs})
	assert.Equal(t, "Pastička", ans[0].Props["doc.title"])
	assert.Equal(t, "SCR: drama", ans[0].Props["doc.txtype"])
}

func TestRefsWithMarkupExtraction(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.Parse([]string{ts7_refs})
	assert.Equal(t, "#1246234523", ans[0].Ref)
	assert.Equal(
		t,
		"Captain Brand of the 'Centipede'<BR>A Pirate of Eminence in the West Indies: "+
			"His Love and Exploits, Together with Some Account of the Singular Manner",
		ans[0].Props["doc.title"],
	)
}

func TestMarkupInTokenAttrs(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.ParseLine(ts8_markup_in_attr)
	assert.Zero(t, ans.ErrMsg)
	assert.Equal(t, 6, len(ans.Text))
	assert.Equal(t, "jj_<sub>", asTokenOrPanic(ans.Text[0]).Attrs["tag"])
	assert.Equal(t, "s", ans.Text[1].(*Struct).Name)
	assert.Equal(t, "a<b", asTokenOrPanic(ans.Text[2]).Attrs["tag"])
	pTok := asTokenOrPanic(ans.Text[3])
	assert.Equal(t, "<P>", pTok.Word)
	assert.Equal(t, "-", pTok.Attrs["lemma"])
	assert.Equal(t, "null", pTok.Attrs["tag"])
	assert.Equal(t, "<i>NN</i>", asTokenOrPanic(ans.Text[4]).Attrs["tag"])
	assert.Equal(t, "s", ans.Text[5].(*CloseStruct).Name)
}

func TestAttrConstSensitivity(t *testing.T) {
	p := NewLineParser([]string{"word", "lemma", "tag"})
	ans := p.ParseLine(ts9_attr_const_sens)
	assert.Zero(t, ans.ErrMsg)
	assert.Len(t, ans.Text, 7)
	assert.Equal(t, MatchTypeKWIC, asTokenOrPanic(ans.Text[3]).MatchType)
	assert.Equal(t, "attractions", asTokenOrPanic(ans.Text[4]).Word)
	assert.Equal(t, "attraction", asTokenOrPanic(ans.Text[4]).Attrs["lemma"])
	assert.Equal(t, "nn2", asTokenOrPanic(ans.Text[4]).Attrs["tag"])
}
