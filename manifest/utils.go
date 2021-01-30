package manifest

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var estSizeRegex *regexp.Regexp
var units map[string]int
var unitsList []string

func init() {
	estSizeRegex = regexp.MustCompile(`^(?P<amount>\d+(?:.\d+)?)[ ]*(?P<unit>[a-zA-Z]+)$`)
	units = map[string]int{"kb": 1000, "mb": 1000000, "gb": 1000000000, "tb": 1000000000000}
	unitsList = []string{"kb", "mb", "gb", "tb"}
}

func GetEstimateToBytes(est string) (int, error) {
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
		return int(float64(multiplier)*amount), nil
	} else {
		return 0, errors.New(fmt.Sprintf("%s -> Could not recognize unit", fn))
	}
}

func GetBytesToEstimate(size int) string {
	unit := ""
	for i, v := range unitsList {
		if (float64(size)/float64(units[v])) < 1.0 {
			if i > 0 {
				unit = unitsList[i-1]
			} else {
				unit = unitsList[i]
			}
		} else if v == unitsList[len(unitsList)-1] {
			unit = unitsList[len(unitsList)-1]
		}

		if unit != "" {
			break;
		}
	}

	return fmt.Sprintf("%.2f %s", (float64(size)/float64(units[unit])), strings.ToUpper(unit))
}