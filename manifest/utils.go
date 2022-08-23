package manifest

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var estSizeRegex *regexp.Regexp
var units map[string]int64
var unitsList []string

func init() {
	estSizeRegex = regexp.MustCompile(`^(?P<amount>\d+(?:.\d+)?)[ ]*(?P<unit>[a-zA-Z]+)$`)
	units = map[string]int64{"kb": 1000, "mb": 1000000, "gb": 1000000000, "tb": 1000000000000}
	unitsList = []string{"kb", "mb", "gb", "tb"}
}

func GetEstimateToBytes(est string) (int64, error) {
	fn := fmt.Sprintf("getEstimateInBytes(est=%s)", est)
	if !estSizeRegex.MatchString(est) {
		return 0, errors.New(fmt.Sprintf("%s -> Could not parse input", fn))
	}
	match := estSizeRegex.FindStringSubmatch(est)

	amount, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("%s -> Could not convert amount to float", fn))
	}

	unit := strings.ToLower(match[2])
	if multiplier, ok := units[unit]; ok {
		return int64(float64(multiplier) * amount), nil
	} else {
		return 0, errors.New(fmt.Sprintf("%s -> Could not recognize unit", fn))
	}
}

func GetBytesToEstimate(size int64) string {
	unit := ""
	for i, v := range unitsList {
		if (float64(size) / float64(units[v])) < 1.0 {
			if i > 0 {
				unit = unitsList[i-1]
			} else {
				unit = unitsList[i]
			}
		} else if v == unitsList[len(unitsList)-1] {
			unit = unitsList[len(unitsList)-1]
		}

		if unit != "" {
			break
		}
	}

	return fmt.Sprintf("%.2f %s", (float64(size) / float64(units[unit])), strings.ToUpper(unit))
}

func ConcatStringSlicesUnique(slice1 []string, slice2 []string) []string {
	processed := map[string]bool{}
	for _, str := range slice1 {
		processed[str] = true
	}

	for _, str := range slice2 {
		if _, ok := processed[str]; !ok {
			slice1 = append(slice1, str)
		}
	}

	return slice1
}

func RemoveIdFromList(ids []int64, id int64) []int64 {
	for idx, idInList := range ids {
		if idInList == id {
			return append(ids[:idx], ids[idx+1:]...)
		}
	}
	return ids
}

func containsStr(list []string, elem string) bool {
	for _, val := range list {
		if val == elem {
			return true
		}
	}
	
	return false
}

func RemoveStrFromList(list []string, elem string) []string {
	for idx, val := range list {
		if elem == val {
			return append(list[:idx], list[idx+1:]...)
		}
	}
	return list
}