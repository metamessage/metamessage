package core

import (
	"errors"
	"fmt"
)

// v0.5.0 全尾部定界：编码器一趟正向直接产出尾定界 wire（body 先落盘、定界信息
// 大端长度 + sign 追加在尾部），无需事后反转、无需回写长度。档位/内联语义与
// v0.3.0 完全一致，仅物理落盘顺序尾化，节点宽度不变。
//
// 解码入口 Untailify 从节点尾部反向定界，把尾定界 wire 还原成「sign 在前」的
// 头部形态，再交给原有正向解码器——解码按从后往前的顺序剥节点，逻辑与原正向
// 解码器解耦。
//
// 各类型 wire 物理布局（与 encode_*.go 原语输出一一对应）：
//   - Simple / int 内联值 / float 内联（exp==-1 且 mantissa<=7）：单字节节点；
//   - int：[body 原序][sign]（长度由 sign 档位反推，无档位大端字节）
//   - float：[exp][mantissa 原序][sign]
//   - str/bytes/container：[body/payload 原序][长度大端 BE][sign]（BE 紧邻 sign 左侧；
//     内联档无 BE 字节，长度值在 sign 低 5 位）
//   - tag 包装段：[inner wire][reverse(meta 正序)][meta_len 槽][total_len 大端 BE][sign]
//
// meta_len 槽三档（末字节为标记或直接值，反向可读）：
//   - <=253：[len]；<=255：[len u8][254]；<=65535：[hi][lo][255]

// 各类型 sign 的 (档位字节宽度, 内联长度/值)；两个返回值来自原 constants.go 辅助函数。
// w==0 表示内联档（无独立档位长度字节）：Simple/int/float 内联为单字节节点，
// str/bytes/container/tag 内联档长度值在 sign 低 5 位，body/payload 依然紧跟 sign。
func tierOf(sign byte) (w int, inline int, ok bool) {
	switch GetPrefix(sign) {
	case Simple:
		return 0, 0, true
	case PositiveInt, NegativeInt:
		w, l := IntLen(sign)
		return w, l, true
	case PrefixFloat:
		w, l := FloatLen(sign)
		return w, l, true
	case PrefixString:
		w, l := StringLen(sign)
		return w, l, true
	case PrefixBytes:
		w, l := BytesLen(sign)
		return w, l, true
	case Container:
		w, l := ContainerLen(sign)
		return w, l, true
	case PrefixTag:
		w, l := TagLen(sign)
		return w, l, true
	}
	return 0, 0, false
}

func tierBytesBE(length, w int) []byte {
	switch w {
	case 0:
		return nil
	case 1:
		return []byte{byte(length)}
	case 2:
		return []byte{byte(length >> 8), byte(length)}
	case 3:
		return []byte{byte(length >> 16), byte(length >> 8), byte(length)}
	default:
		return []byte{byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)}
	}
}

func readTierBE(data []byte, end, w int) (int, error) {
	if end-w < 0 {
		return 0, errors.New("truncated: tier bytes out of range")
	}
	v := 0
	for i := end - w; i < end; i++ {
		v = v<<8 | int(data[i])
	}
	return v, nil
}

// nodeTailWidth 计算 wire 中 [start,end) 单节点的定界（返回节点起始偏移）。
// wire 节点与 head 节点等宽：首字节即该节点 sign，档位信息原位可读。
func nodeTailStart(wire []byte, end int) (int, error) {
	if end < 1 || end > len(wire) {
		return 0, fmt.Errorf("truncated: node end %d out of range", end)
	}
	// wire 尾 sign 在节点末端；宽度需要按 sign 计算，而 sign 在 head 形态是首字节。
	// 等宽性允许我们先猜起点，但反向定界必须以 sign 为准：从尾读。
	sign := wire[end-1]
	w, inline, ok := tierOf(sign)
	if !ok {
		return 0, fmt.Errorf("invalid prefix %d", GetPrefix(sign))
	}
	switch GetPrefix(sign) {
	case Simple:
		return end - 1, nil
	case PositiveInt:
		fallthrough
	case NegativeInt:
		if w == 0 {
			return end - 1, nil
		}
		// int head [sign][body w 字节] → wire [body w][sign]
		return end - 1 - w, nil
	case PrefixFloat:
		if w == 0 {
			return end - 1, nil
		}
		// float head [sign][exp][mant w 字节]：body 实宽 w+1 → wire [exp][mant][sign]
		return end - 2 - w, nil
	case PrefixString:
		fallthrough
	case PrefixBytes:
		fallthrough
	case Container:
		if w == 0 {
			// 内联档：head [sign][body inline] → wire [body inline][sign]
			return end - 1 - inline, nil
		}
		l, err := readTierBE(wire, end-1, w)
		if err != nil {
			return 0, err
		}
		return end - 1 - w - l, nil
	case PrefixTag:
		total := inline
		if w > 0 {
			var err error
			total, err = readTierBE(wire, end-1, w)
			if err != nil {
				return 0, err
			}
		}
		slotEnd := end - 1 - w
		ml, slotW, err := parseFooterSlot(wire, slotEnd)
		if err != nil {
			return 0, err
		}
		// total = meta 前缀槽宽 + meta + inner（与头部定界定义一致）；
		// wire 尾段 = reverse(meta)+slot，inner 区在 slot 左侧。
		// 反推用实际尾槽宽 slotW：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
		// 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
		inner := total - ml - slotW
		if inner < 0 {
			return 0, errors.New("tag total underflow")
		}
		return slotEnd - slotW - ml - inner, nil
	}
	return 0, fmt.Errorf("invalid prefix")
}

// untailifyNode：wire 形态 [start,end) → head 形态。
func untailifyNode(wire []byte, start, end int) ([]byte, error) {
	if start < 0 || end > len(wire) || start >= end {
		return nil, fmt.Errorf("untailify: bad range [%d,%d)", start, end)
	}
	sign := wire[end-1]
	w, inline, ok := tierOf(sign)
	if !ok {
		return nil, fmt.Errorf("untailify: invalid prefix %d", GetPrefix(sign))
	}
	_ = inline
	switch GetPrefix(sign) {
	case Simple:
		return []byte{sign}, nil
	case PositiveInt, NegativeInt, PrefixFloat:
		if w == 0 {
			return []byte{sign}, nil
		}
		out := append([]byte{sign}, wire[start:end-1]...)
		return out, nil
	case PrefixString, PrefixBytes, Container:
		l := inline
		if w > 0 {
			var err error
			l, err = readTierBE(wire, end-1, w)
			if err != nil {
				return nil, err
			}
		}
		body := wire[start : end-1-w]
		if len(body) != l {
			return nil, fmt.Errorf("untailify: payload %d, tier says %d", len(body), l)
		}
		out := append([]byte{sign}, tierBytesBE(l, w)...)
		if GetPrefix(sign) == Container {
			child, err := untailifyChildren(body)
			if err != nil {
				return nil, err
			}
			return append(out, child...), nil
		}
		return append(out, body...), nil
	case PrefixTag:
		total := inline
		if w > 0 {
			var err error
			total, err = readTierBE(wire, end-1, w)
			if err != nil {
				return nil, err
			}
		}
		slotEnd := end - 1 - w
		ml, slotW, err := parseFooterSlot(wire, slotEnd)
		if err != nil {
			return nil, err
		}
		innerLen := total - ml - slotW
		if innerLen < 0 {
			return nil, errors.New("untailify tag: total underflow")
		}
		metaEnd := slotEnd - slotW
		meta := reverseBytes(wire[metaEnd-ml : metaEnd])
		innerWire := wire[start : metaEnd-ml]
		if len(innerWire) != innerLen {
			return nil, fmt.Errorf("untailify tag: inner %d, want %d", len(innerWire), innerLen)
		}
		var inner []byte
		if innerLen > 0 {
			v, err := untailifyNode(innerWire, 0, len(innerWire))
			if err != nil {
				return nil, err
			}
			inner = v
		}
		prefix, err := metaPrefixBytes(ml)
		if err != nil {
			return nil, err
		}
		out := append([]byte{sign}, tierBytesBE(total, w)...)
		out = append(out, prefix...)
		out = append(out, meta...)
		return append(out, inner...), nil
	}
	return nil, errors.New("untailify: unreachable")
}

// metaPrefixBytes 还原正向 meta 前缀槽：与尾部槽同一套 canonical 三档，
// 宽度恒等于 footerSlotBytes 槽宽（ml<254→1，254..255→2，256..65535→3）。
func metaPrefixBytes(ml int) ([]byte, error) {
	switch {
	case ml < 254:
		return []byte{byte(ml)}, nil
	case ml <= Max1Byte:
		return []byte{FooterSlotMarker1, byte(ml)}, nil
	case ml <= Max2Byte:
		return []byte{FooterSlotMarker2, byte(ml >> 8), byte(ml)}, nil
	default:
		return nil, fmt.Errorf("meta too long: %d", ml)
	}
}

// untailifyChildren：wire 容器 payload 内顺序 wire 子节点逐个还原。
// wire 子节点以尾 sign 结尾，从区域末端反向定界扫描，还原后按正序拼接。
func untailifyChildren(payload []byte) ([]byte, error) {
	nodes := make([][]byte, 0, 8)
	for end := len(payload); end > 0; {
		start, err := nodeTailStart(payload, end)
		if err != nil {
			return nil, err
		}
		if start < 0 || start >= end {
			return nil, fmt.Errorf("untailify children: bad boundary %d..%d", start, end)
		}
		node, err := untailifyNode(payload, start, end)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
		end = start
	}
	out := make([]byte, 0, len(payload))
	for i := len(nodes) - 1; i >= 0; i-- {
		out = append(out, nodes[i]...)
	}
	return out, nil
}

// Untailify 解码入口：整帧 wire 形态 → head 形态。
func Untailify(wire []byte) ([]byte, error) {
	if len(wire) == 0 {
		return nil, errors.New("empty wire")
	}
	start, err := nodeTailStart(wire, len(wire))
	if err != nil {
		return nil, err
	}
	if start != 0 {
		return nil, fmt.Errorf("footer boundary mismatch: node start %d, want 0", start)
	}
	return untailifyNode(wire, start, len(wire))
}
