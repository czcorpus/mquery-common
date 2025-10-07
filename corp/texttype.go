package corp

const (
	TextPropertyAuthor      = "author"
	TextPropertyTitle       = "title"
	TextPropertyPubYear     = "publication-year"
	TextPropertyTranslator  = "translator"
	TextPropertyOriginaLang = "original-language"
	TextPropertyTextType    = "text-type"
)

// TextProperty is a generalized text type property used in APIs across multiple
// corpora. E.g. instead of structural attributes like doc.author text.orig_author,
// we offer "author" and map it automatically to the raw structural attribute.
type TextProperty string

func (tp TextProperty) Validate() bool {
	return tp == TextPropertyAuthor || tp == TextPropertyTitle ||
		tp == TextPropertyPubYear || tp == TextPropertyTextType ||
		tp == TextPropertyTranslator || tp == TextPropertyOriginaLang
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

type TTPropertyConf struct {
	Name         string `json:"name"`
	IsInOverview bool   `json:"isInOverview"`
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
