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

func TestRefsRegexp(t *testing.T) {
	s := "#36940724,doc.title=Snídaně v poledne,doc.id=snidane-v-poledne"
	srch := refsRegexp.FindAllString(s, -1)
	assert.Equal(t, 3, len(srch))
	assert.Equal(t, "#36940724", srch[0])
	assert.Equal(t, "doc.title=Snídaně v poledne", srch[1])
	assert.Equal(t, "doc.id=snidane-v-poledne", srch[2])
}
