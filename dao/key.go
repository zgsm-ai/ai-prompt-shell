package dao

import "strings"

func KeyToPath(key, prefix string) string {
	return strings.ReplaceAll(strings.TrimPrefix(key, prefix), ":", ".")
}

func KeyToID(key, prefix string) string {
	return strings.ReplaceAll(strings.TrimPrefix(key, prefix), ":", ".")
}
