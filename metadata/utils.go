package metadata

func RemoveIdFromList(ids []int64, id int64) []int64 {
	for idx, idInList := range ids {
		if idInList == id {
			return append(ids[:idx], ids[idx+1:]...)
		}
	}
	return ids
}
