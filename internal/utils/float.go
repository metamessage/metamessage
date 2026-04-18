package utils

import "strconv"

func FormatFloat32(val float32) string {
	f := float64(val)

	s := strconv.FormatFloat(f, 'f', -1, 32)

	if _, err := strconv.Atoi(s); err == nil {
		return s + ".0"
	}

	return s
}

func FormatFloat64(val float64) string {
	s := strconv.FormatFloat(val, 'f', -1, 64)

	if _, err := strconv.Atoi(s); err == nil {
		return s + ".0"
	}

	return s
}
