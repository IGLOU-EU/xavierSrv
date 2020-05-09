package tools

import "strings"

func JsonEscapeString(d string) string {
	r := strings.NewReplacer(
		"\\", "\\\\",
		"/", "\\/",
		"\"", "\\\"",
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\x08", "\\f",
		"\x0c", "\\b",
	)
	return r.Replace(d)
}
