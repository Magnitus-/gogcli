package sdk

import "strings"

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func replaceUnicodeEncodedAscii(val string) string {
	return strings.Replace(
		strings.Replace(
			strings.Replace(val, "\u0026", "&", -1),
			"\u003e", 
			">",
			-1,
		),
		"\u003c",
		"<",
		-1,
	)
}
