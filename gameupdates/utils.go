package gameupdates

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