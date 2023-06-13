package text

import (
	"strings"

	"github.com/alexsergivan/transliterator"
)

// TransliterateIDX performs transliteration (en) of the input text.
func TransliterateIDX(in string) string {
	if in == "" {
		return ""
	}

	tr := transliterator.NewTransliterator(nil)

	return strings.ToLower(strings.TrimSpace(tr.Transliterate(in, "")))
}
