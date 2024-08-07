package concordance

import "strings"

func isElement(tagSrc string) bool {
	return strings.HasPrefix(tagSrc, "<") && strings.HasSuffix(tagSrc, ">")
}

func isOpenElement(tagSrc string) bool {
	return isElement(tagSrc) && !strings.HasPrefix(tagSrc, "</") &&
		!strings.HasSuffix(tagSrc, "/>")
}

func isCloseElement(tagSrc string) bool {
	return isElement(tagSrc) && strings.HasPrefix(tagSrc, "</")
}

func isSelfCloseElement(tagSrc string) bool {
	return isElement(tagSrc) && strings.HasSuffix(tagSrc, "/>")
}
