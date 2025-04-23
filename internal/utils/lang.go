package utils

import "strings"

func ParseLang(lang string) string {
	if lang == "" {
		return "en"
	}
	var langs = "en,de,fr,it,es,zh,ja,pt,ru,nl"
	if strings.Contains(langs, lang) {
		return lang
	}
	return "en"
}
