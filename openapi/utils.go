package openapi

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func PathConverter(p string) string {
	return strings.ReplaceAll(strings.ReplaceAll(p, "{", ":"), "}", "")
}

func RefObject(p string) string {
	list := strings.Split(p, "/")
	return list[len(list)-1]
}

func FirstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}

func FirstToUpper(s string) string {
	return strings.Title(s)
}
