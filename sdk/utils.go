package sdk

import (
	"fmt"
	"regexp"
)

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func contains(list []int64, elem int64) bool {
	for _, val := range list {
		if val == elem {
			return true
		}
	}
	
	return false
}

func containsStr(list []string, elem string) bool {
	for _, val := range list {
		if val == elem {
			return true
		}
	}
	
	return false
}

var imageUrlRegex *regexp.Regexp

func init() {
	imageUrlRegex = regexp.MustCompile(`^(?P<body>.*)(?P<suffix>[.][a-z]{3,6})$`)
}

func ensureFileNameSuffix(name string, suffix string) string {
	if !imageUrlRegex.MatchString(name) {
		return fmt.Sprintf("%s%s", name, suffix)
	}

	match := imageUrlRegex.FindStringSubmatch(name)
	return fmt.Sprintf("%s%s", string(match[1]), suffix)
}