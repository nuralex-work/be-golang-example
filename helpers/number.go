package helpers

import "strconv"

func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}
func UintPtr(v uint) *uint {
	return &v
}

func FloatPtr(v float64) *float64 {
	return &v
}
