package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ParseStringToUint64 parses a numeric string into sign + exponent (signed) + mantissa.
// Supported formats: integers (123), decimals (123.456), scientific notation (123.456e7 / -123e-8).
// Return values:
//
//	isNegative: Whether it is a negative number.
//	exponent: Exponent (signed, e.g., for 1.23e5 → exponent = 5 - 2 = 3; for 1.23e-5 → exponent = -5 - 2 = -7).
//	mantissa: Mantissa (unsigned, e.g., for 1.23 → mantissa = 123).
//	err: Parsing error.
func ParseStringToUint64(s string) (isNegative bool, exponent int8, mantissa uint64, err error) {
	if s == "" {
		return false, 0, 0, errors.New("empty numeric string")
	}

	isNegative = strings.HasPrefix(s, "-")
	if isNegative {
		s = s[1:]
		if s == "" {
			return false, 0, 0, errors.New("invalid numeric string: only minus sign")
		}
	}

	var expPart string
	if eIdx := strings.IndexAny(s, "eE"); eIdx != -1 {
		expPart = s[eIdx+1:]
		s = s[:eIdx]
		if expPart == "" {
			return false, 0, 0, errors.New("missing exponent part in scientific notation")
		}
	}

	parts := strings.Split(s, ".")
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = parts[1]
	}

	if intPart == "" {
		intPart = "0"
	}

	// Calculate the base exponent (the number of decimal places, a negative number indicates a left - shift is required).
	baseExp := int64(-len(fracPart)) // For example, "1.23" → baseExp = - 2

	// Handle the exponent in scientific notation
	if expPart != "" {
		var exp int64
		exp, err = strconv.ParseInt(expPart, 10, 64)
		if err != nil {
			return false, 0, 0, fmt.Errorf("invalid exponent: %w", err)
		}

		// Merge exponents: exponent in scientific notation + exponent of the number of decimal places
		baseExp += exp
	}

	if baseExp < math.MinInt8 || baseExp > math.MaxInt8 {
		return false, 0, 0, fmt.Errorf("final exponent out of range (%d ~ %d): %d", math.MinInt8, math.MaxInt8, baseExp)
	}

	exponent = int8(baseExp)

	// Concatenate the integer and fractional parts into the mantissa (remove leading zeros, so "00123" becomes "123").
	mantissaStr := strings.TrimLeft(intPart+fracPart, "0")

	// Handle scenarios where all values are zero (such as "0", "0.000", "0e10").
	if mantissaStr == "" {
		mantissaStr = "0"
	}

	// Parse the mantissa and handle overflow.
	if mantissa, err = strconv.ParseUint(mantissaStr, 10, 64); err != nil {
		// Explicitly indicate overflow to facilitate handling at a higher level (such as storing as a string).
		var numErr *strconv.NumError
		if errors.As(err, &numErr) && errors.Is(numErr.Err, strconv.ErrRange) {
			return false, 0, 0, fmt.Errorf("mantissa overflow (exceeds uint64 max): %s", mantissaStr)
		}
		return false, 0, 0, fmt.Errorf("invalid mantissa: %w", err)
	}

	return isNegative, exponent, mantissa, nil
}

// CamelToSnake
// Example:
// Datetime     → datetime
// ArrArrPrtTime → arr_arr_prt_time
// ID           → id
// UserID       → user_id
// HTTPRequest  → http_request
func CamelToSnake(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	result.Grow(len(s) + 2)

	for i, char := range s {
		if unicode.IsUpper(char) {
			if i > 0 &&
				(!unicode.IsUpper(rune(s[i-1])) || (i+1 < len(s) && !unicode.IsUpper(rune(s[i+1])))) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

var DefaultTime = time.Unix(0, 0).UTC()

func GetLocationOffsetHour(loc *time.Location) int {
	if loc == nil {
		return 0
	}

	_, offsetSec := DefaultTime.In(loc).Zone()
	return offsetSec / 3600
}

func IntToLocation(offsetHours int) *time.Location {
	offsetSeconds := offsetHours * 60 * 60
	return time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetSeconds)
}
