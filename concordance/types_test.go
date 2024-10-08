// Copyright 2024 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2024 Institute of the Czech National Corpus,
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

func TestRefsRegexp(t *testing.T) {
	s := "#36940724,doc.title=Snídaně v poledne,doc.id=snidane-v-poledne"
	srch := refsRegexp.FindAllString(s, -1)
	assert.Equal(t, 3, len(srch))
	assert.Equal(t, "#36940724", srch[0])
	assert.Equal(t, "doc.title=Snídaně v poledne", srch[1])
	assert.Equal(t, "doc.id=snidane-v-poledne", srch[2])
}
