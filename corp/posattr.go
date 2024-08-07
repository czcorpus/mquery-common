package corp

// Single corpus configuration types
// ----------------------------------------

type PosAttr struct {
	Name        string            `json:"name"`
	Description map[string]string `json:"description"`
}

func (p PosAttr) IsZero() bool {
	return p.Name == ""
}

func (p PosAttr) LocaleDescription(lang string) string {
	d := p.Description[lang]
	if d != "" {
		return d
	}
	return p.Description["en"]
}

type PosAttrList []PosAttr

func (pal PosAttrList) GetIDs() []string {
	ans := make([]string, len(pal))
	for i, v := range pal {
		ans[i] = v.Name
	}
	return ans
}
