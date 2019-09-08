package utils

import "sort"

func IdsFitter(ids []int) []int {
	sort.Ints(ids)
	var newIds []int
	var lastId int
	for k, v := range ids {
		if k == 0 {
			lastId = v
			newIds = append(newIds, v)
		}
		if k > 0 && v != lastId {
			lastId = v
			newIds = append(newIds, v)
		}
	}
	return newIds
}
