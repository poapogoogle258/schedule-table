package util

import "strings"

func GetExpressionFile(s string) string {
	i := strings.IndexRune(s, '.')
	return s[i+1:]
}
