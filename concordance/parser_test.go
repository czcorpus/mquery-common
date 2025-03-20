// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Martin Zimandl <martin.zimandl@gmail.com>
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
