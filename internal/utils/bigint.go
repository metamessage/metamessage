package utils

import (
	"fmt"
	"io"
	"strings"
)

func EncodeBigInt(w io.ByteWriter, s string) uint32 {
	if s == "" {
		return 0
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}

	var bits []int
	if neg {
		bits = append(bits, 1)
	} else {
		bits = append(bits, 0)
	}

	n := len(s)
	i := 0
	for ; i < n; i += 3 {
		rem := n - i
		if rem >= 3 {
			num := atoi(s[i : i+3])
			bits = append(bits, toBits(num, 10)...)
		} else if rem == 2 {
			num := atoi(s[i : i+2])
			bits = append(bits, toBits(num, 7)...)
		} else {
			num := atoi(s[i : i+1])
			bits = append(bits, toBits(num, 4)...)
		}
	}

	return writeBits(w, bits)
}

func DecodeBigInt(data []byte, n int) (string, error) {
	bits := bytesToBits(data)
	if len(bits) == 0 {
		return "", nil
	}

	neg := bits[0] == 1
	bits = bits[1:]

	var numStr strings.Builder
	for n > 0 {
		if n >= 3 && len(bits) >= 10 {
			num := fromBits(bits[:10])
			fmt.Fprintf(&numStr, "%03d", num)
			bits = bits[10:]
			n -= 3
		} else if n >= 2 && len(bits) >= 7 {
			num := fromBits(bits[:7])
			fmt.Fprintf(&numStr, "%02d", num)
			bits = bits[7:]
			n -= 2
		} else if n >= 1 && len(bits) >= 4 {
			num := fromBits(bits[:4])
			fmt.Fprint(&numStr, num)
			bits = bits[4:]
			n -= 1
		} else {
			break
		}
	}

	res := numStr.String()
	if neg {
		res = "-" + res
	}

	return res, nil
}

func EncodeBigIntPositive(w io.ByteWriter, s string) uint32 {
	if s == "" {
		return 0
	}

	var bits []int

	n := len(s)
	i := 0
	for ; i < n; i += 3 {
		rem := n - i
		if rem >= 3 {
			num := atoi(s[i : i+3])
			bits = append(bits, toBits(num, 10)...)
		} else if rem == 2 {
			num := atoi(s[i : i+2])
			bits = append(bits, toBits(num, 7)...)
		} else {
			num := atoi(s[i : i+1])
			bits = append(bits, toBits(num, 4)...)
		}
	}

	return writeBits(w, bits)
}

func DecodeBigIntPositive(data []byte, n int) (string, error) {
	bits := bytesToBits(data)
	if len(bits) == 0 {
		return "", nil
	}

	var numStr strings.Builder
	for n > 0 {
		if n >= 3 && len(bits) >= 10 {
			num := fromBits(bits[:10])
			fmt.Fprintf(&numStr, "%03d", num)
			bits = bits[10:]
			n -= 3
		} else if n >= 2 && len(bits) >= 7 {
			num := fromBits(bits[:7])
			fmt.Fprintf(&numStr, "%02d", num)
			bits = bits[7:]
			n -= 2
		} else if n >= 1 && len(bits) >= 4 {
			num := fromBits(bits[:4])
			fmt.Fprint(&numStr, num)
			bits = bits[4:]
			n -= 1
		} else {
			break
		}
	}

	res := numStr.String()

	return res, nil
}

func atoi(s string) int {
	v := 0
	for _, c := range s {
		v = v*10 + int(c-'0')
	}
	return v
}

func toBits(v int, n int) []int {
	b := make([]int, n)
	for i := range n {
		b[n-1-i] = (v >> i) & 1
	}
	return b
}

func fromBits(b []int) int {
	v := 0
	for _, bit := range b {
		v = v<<1 | bit
	}
	return v
}

func writeBits(w io.ByteWriter, bits []int) (n uint32) {
	bt := byte(0)
	bl := 0
	for _, b := range bits {
		bt = bt<<1 | byte(b)
		bl++
		if bl == 8 {
			w.WriteByte(bt)
			n++
			bt = 0
			bl = 0
		}
	}
	if bl > 0 {
		bt <<= (8 - bl)
		w.WriteByte(bt)
		n++
	}
	return
}

func bytesToBits(data []byte) []int {
	totalBits := len(data) * 8
	bits := make([]int, 0, totalBits)
	for _, bt := range data {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((bt>>i)&1))
		}
	}
	return bits
}

func CalcBigIntSize(s string) int {
	s = strings.TrimPrefix(s, "-")

	n := len(s)
	if n == 0 {
		return 0
	}

	totalBits := 1
	groups := n / 3
	rem := n % 3

	totalBits += groups * 10

	switch rem {
	case 1:
		totalBits += 4
	case 2:
		totalBits += 7
	}

	bytesSize := (totalBits + 7) / 8
	return bytesSize
}

func CalcBigIntSizePositive(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}

	totalBits := 0
	groups := n / 3
	rem := n % 3

	totalBits += groups * 10

	switch rem {
	case 1:
		totalBits += 4
	case 2:
		totalBits += 7
	}

	bytesSize := (totalBits + 7) / 8
	return bytesSize
}
