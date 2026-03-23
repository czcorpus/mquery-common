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

package corp

const (
	TextPropertyAuthor      TextProperty = "author"
	TextPropertyTitle       TextProperty = "title"
	TextPropertyPubYear     TextProperty = "publication-year"
	TextPropertyPubDate     TextProperty = "publication-date"
	TextPropertyTranslator  TextProperty = "translator"
	TextPropertyOriginaLang TextProperty = "original-language"
	TextPropertyTextType    TextProperty = "text-type"
	TextPropertyMedium      TextProperty = "medium"
)

// TextProperty is a generalized text type property used in APIs across multiple
// corpora. E.g. instead of structural attributes like doc.author text.orig_author,
// we offer "author" and map it automatically to the raw structural attribute.
type TextProperty string

func (tp TextProperty) Validate() bool {
	return tp == TextPropertyAuthor || tp == TextPropertyTitle ||
		tp == TextPropertyPubYear || tp == TextPropertyPubDate ||
		tp == TextPropertyTextType || tp == TextPropertyTranslator ||
		tp == TextPropertyOriginaLang || tp == TextPropertyMedium
}

func (tp TextProperty) String() string {
	return string(tp)
}

func (tp TextProperty) IsZero() bool {
	return tp == ""
}

// StructAttr is a raw corpus structural attribute (doc.id, text.author, doc.pubyear etc.)
type StructAttr struct {
	Name        string            `json:"name"`
	Description map[string]string `json:"description"`
}

func (s StructAttr) LocaleDescription(lang string) string {
	d := s.Description[lang]
	if d != "" {
		return d
	}
	return s.Description["en"]
}

func (s StructAttr) IsZero() bool {
	return s.Name == ""
}

// -------------

type TextTypes map[string][]string

// ------

// TTPropertyConf is a universal/common text property (mapped to
// some real physical one in a corpus).
type TTPropertyConf struct {
	Name         string `json:"name"`
	IsInOverview bool   `json:"isInOverview"`

	// DateFormat can be used in case the attribute represents a date (year, year+month+day,...).
	// The format is the same as Golang time parsing layout string,
	// e.g.: "2006-01-02" for days granularity in ISO8601 format, "2006" for years etc.
	DateFormat string `json:"dateFormat"`
}

func (ttpc TTPropertyConf) IsDateType() bool {
	return ttpc.DateFormat != ""
}

// ------

// TextTypeProperties maps between generalized text properties
// and specific corpus structural attributes.
type TextTypeProperties map[TextProperty]TTPropertyConf

// Prop returns a generalized property based on provided struct. attribute
// If nothing is found, empty TextProperty is returned
func (ttp TextTypeProperties) Prop(attr string) TextProperty {
	for k, v := range ttp {
		if v.Name == attr {
			return k
		}
	}
	return ""
}

func (ttp TextTypeProperties) List() []TextProperty {
	ans := make([]TextProperty, len(ttp))
	var i int
	for k := range ttp {
		ans[i] = k
	}
	return ans
}

func (ttp TextTypeProperties) ListOverviewProps() []TextProperty {
	ans := make([]TextProperty, 0, len(ttp))
	for _, v := range ttp {
		if v.IsInOverview {
			ans = append(ans, TextProperty(v.Name))
		}
	}
	return ans
}

// Attr returns a struct. attribute name based on generalized property.
// If nothing is found, empty string is returned.
func (ttp TextTypeProperties) Attr(prop TextProperty) string {
	return ttp[prop].Name
}

// ---------------

// Subcorpus represents a subcorpus created by selecting specific
// values out of different structural attributes.
type Subcorpus struct {
	ID          string            `json:"id"`
	TextTypes   TextTypes         `json:"textTypes"`
	Description map[string]string `json:"description"`
}
