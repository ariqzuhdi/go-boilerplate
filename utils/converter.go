package utils

import "strconv"

func ParseStringToUint(s string) uint {
	idInt, _ := strconv.ParseUint(s, 10, 64)
	return uint(idInt)
}
