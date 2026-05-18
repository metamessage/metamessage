package utils

import (
	"strconv"
	"strings"
)

func FormatFloat32(val float32) string {
	f := float64(val)

	s := strconv.FormatFloat(f, 'f', -1, 32)

	if !strings.Contains(s, ".") {
		return s + ".0"
	}

	return s
}

func FormatFloat64(val float64) string {
	s := strconv.FormatFloat(val, 'f', -1, 64)

	if !strings.Contains(s, ".") {
		return s + ".0"
	}

	return s
}
