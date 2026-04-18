package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"testing/quick"
)

// go test internal/utils/** -v
//
// go test internal/utils/** -v -run TestParseStringToUint64/positive_integer
//
// go test internal/utils/** -bench=. -benchmem

func TestParseStringToUint64(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		wantNeg      bool
		wantExp      int16
		wantMantissa uint64
		wantErr      error
	}

	testCases := []testCase{
		{
			name:         "positive_integer",
			input:        "123",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 123,
			wantErr:      nil,
		},
		{
			name:         "negative integer",
			input:        "-456",
			wantNeg:      true,
			wantExp:      0,
			wantMantissa: 456,
			wantErr:      nil,
		},
		{
			name:         "zero integer",
			input:        "0",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      nil,
		},
		{
			name:         "integer with leading zero",
			input:        "00789",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 789,
			wantErr:      nil,
		},

		{
			name:         "positive decimal",
			input:        "123.456",
			wantNeg:      false,
			wantExp:      -3, // The fractional part has 3 digits. → baseExp=-3
			wantMantissa: 123456,
			wantErr:      nil,
		},
		{
			name:         "negative decimal",
			input:        "-78.9",
			wantNeg:      true,
			wantExp:      -1, // The fractional part has one digit. → baseExp=-1
			wantMantissa: 789,
			wantErr:      nil,
		},
		{
			name:         "decimal with leading zero (int part)",
			input:        "012.34",
			wantNeg:      false,
			wantExp:      -2,
			wantMantissa: 1234,
			wantErr:      nil,
		},
		{
			name:         "decimal with trailing zero (frac part)",
			input:        "123.4500",
			wantNeg:      false,
			wantExp:      -4, // The fractional part has four digits. → baseExp=-4
			wantMantissa: 1234500,
			wantErr:      nil,
		},
		{
			name:         "decimal with only frac part",
			input:        ".678",
			wantNeg:      false,
			wantExp:      -3,
			wantMantissa: 678,
			wantErr:      nil,
		},
		{
			name:         "decimal with only int part (trailing dot)",
			input:        "999.",
			wantNeg:      false,
			wantExp:      0, // The fractional part is empty. → baseExp=0
			wantMantissa: 999,
			wantErr:      nil,
		},
		{
			name:         "zero_decimal",
			input:        "0.000",
			wantNeg:      false,
			wantExp:      -3,
			wantMantissa: 0,
			wantErr:      nil,
		},

		{
			name:         "scientific positive exp",
			input:        "1.23e5",
			wantNeg:      false,
			wantExp:      3, // baseExp=-2 + 5 = 3
			wantMantissa: 123,
			wantErr:      nil,
		},
		{
			name:         "scientific negative exp",
			input:        "4.56e-8",
			wantNeg:      false,
			wantExp:      -10, // baseExp=-2 + (-8) = -10
			wantMantissa: 456,
			wantErr:      nil,
		},
		{
			name:         "scientific uppercase E",
			input:        "7.89E10",
			wantNeg:      false,
			wantExp:      8, // baseExp=-2 + 10 = 9
			wantMantissa: 789,
			wantErr:      nil,
		},
		{
			name:         "scientific with negative number",
			input:        "-1.234e-3",
			wantNeg:      true,
			wantExp:      -6, // baseExp=-3 + (-3) = -6
			wantMantissa: 1234,
			wantErr:      nil,
		},
		{
			name:         "scientific zero",
			input:        "0e100",
			wantNeg:      false,
			wantExp:      100,
			wantMantissa: 0,
			wantErr:      nil,
		},

		{
			name:         "exponent min boundary",
			input:        "1e-127",
			wantNeg:      false,
			wantExp:      -127, // baseExp=0 + (-127) = -127
			wantMantissa: 1,
			wantErr:      nil,
		},
		{
			name:         "exponent max boundary",
			input:        "1e126",
			wantNeg:      false,
			wantExp:      126, // baseExp=0 + 128 = 128
			wantMantissa: 1,
			wantErr:      nil,
		},
		{
			name:         "exponent exceed min",
			input:        "1e-129",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      fmt.Errorf("final exponent out of range (%d ~ %d)", -128, 127),
		},
		{
			name:         "exponent exceed max",
			input:        "1e128",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      fmt.Errorf("final exponent out of range (%d ~ %d)", -128, 127),
		},

		{
			name:         "mantissa overflow (uint64 max +1)",
			input:        "18446744073709551616", // uint64 max=18446744073709551615
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      errors.New("mantissa overflow (exceeds uint64 max): 18446744073709551616"),
		},

		{
			name:         "empty string",
			input:        "",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      errors.New("empty numeric string"),
		},
		{
			name:         "only minus sign",
			input:        "-",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      errors.New("invalid numeric string: only minus sign"),
		},
		{
			name:         "invalid scientific (missing exp part)",
			input:        "123e",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      errors.New("missing exponent part in scientific notation"),
		},
		{
			name:         "invalid scientific (non-numeric exp)",
			input:        "123eabc",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      fmt.Errorf("invalid exponent: %w", &strconv.NumError{Func: "ParseInt", Num: "abc", Err: strconv.ErrSyntax}),
		},
		{
			name:         "non-numeric input",
			input:        "abc123",
			wantNeg:      false,
			wantExp:      0,
			wantMantissa: 0,
			wantErr:      fmt.Errorf("invalid mantissa: %w", &strconv.NumError{Func: "ParseUint", Num: "abc123", Err: strconv.ErrSyntax}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotNeg, gotExp, gotMantissa, gotErr := ParseStringToUint64(tc.input)

			// Assertion error (prioritize error checking to avoid subsequent assertion panics)
			if tc.wantErr == nil {
				if gotErr != nil {
					t.Fatalf("expected no error, but got: %v", gotErr)
				}
			} else {
				if gotErr == nil {
					t.Fatalf("expected error: %v, but got no error", tc.wantErr)
				}
				// Match the error message (ignore the specific type of the underlying wrapper and only match the error text)
				if !strings.Contains(gotErr.Error(), tc.wantErr.Error()) {
					t.Errorf("error mismatch: expected %q, got %q", tc.wantErr.Error(), gotErr.Error())
				}
			}

			//Assert non - error fields (only when there is no error)
			if tc.wantErr == nil {
				if gotNeg != tc.wantNeg {
					t.Errorf("isNegative mismatch: expected %t, got %t", tc.wantNeg, gotNeg)
				}
				if int16(gotExp) != tc.wantExp {
					t.Errorf("exponent mismatch: expected %d, got %d", tc.wantExp, gotExp)
				}
				if gotMantissa != tc.wantMantissa {
					t.Errorf("mantissa mismatch: expected %d, got %d", tc.wantMantissa, gotMantissa)
				}
			}
		})
	}
}

func BenchmarkParseStringToUint64(b *testing.B) {
	testInputs := []string{"123", "123.456", "1.23e5", "-78.9e-3"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = ParseStringToUint64(testInputs[i%len(testInputs)])
	}
}

func TestParseStringToUint64_Quick(t *testing.T) {
	f := func(s string) bool {
		if !regexp.MustCompile(`^-?\d+(\.\d+)?(eE-?\d+)?$`).MatchString(s) {
			return true
		}
		_, _, _, err := ParseStringToUint64(s)
		// Only verify that there is no panic, without verifying the result (suitable for preliminary checks).
		return err == nil || strings.Contains(err.Error(), "out of range") || strings.Contains(err.Error(), "overflow")
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
