package mm

import (
	"encoding/json"
	"testing"

	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/jsonc/parser"
	"github.com/metamessage/metamessage/internal/jsonc/scanner"
)

type X64Breakpoint struct {
	Address  int    `mm:"desc=address"`
	Name     string `mm:"desc=name"`
	Enabled  bool   `mm:"type=bool"`
	HitCount int    `mm:"desc=hit count"`
}

type X64DebugState struct {
	Running     bool   `mm:"type=bool"`
	Paused      bool   `mm:"type=bool"`
	Pid         int    `mm:"desc=process id"`
	ThreadCount int    `mm:"desc=thread count"`
}

func TestX64BreakpointEncode(t *testing.T) {
	bp := X64Breakpoint{
		Address:  401000,
		Name:     "main",
		Enabled:  true,
		HitCount: 5,
	}

	bs, err := FromValue(bp, "")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	if len(bs) == 0 {
		t.Fatal("encoded bytes is empty")
	}

	node, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	t.Logf("Encoded breakpoint: %x", bs)
	t.Logf("JSONC: %s", jsonc.ToJSONC(node))
}

func TestX64DebugStateEncode(t *testing.T) {
	state := X64DebugState{
		Running:     false,
		Paused:      true,
		Pid:         1234,
		ThreadCount: 2,
	}

	bs, err := FromValue(state, "")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	if len(bs) == 0 {
		t.Fatal("encoded bytes is empty")
	}

	node, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if state.Paused != true {
		t.Error("state should be paused")
	}

	t.Logf("Encoded debug state: %d bytes", len(bs))
	t.Logf("JSONC: %s", jsonc.ToJSONC(node))
}

func TestX64JSONCParse(t *testing.T) {
	jsoncInput := `{
  // Debug state for x64dbg
  "running": true,
  "paused": false,
  "pid": 1234,
  "breakpoints": [
    {
      "address": 4198400,
      "name": "main",
      "enabled": true
    }
  ]
}`

	sc := scanner.New(jsoncInput)
	tokens := sc.ScanAll()
	p := parser.New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	t.Logf("Parsed JSONC: %s", jsonc.ToJSONC(node))
}

func TestX64JSONCComplex(t *testing.T) {
	jsoncInput := `{
  // Complete debug state for x64dbg plugin communication
  "running": true,
  "paused": true,
  "pid": 5678,
  "breakpoints": [
    {
      "address": 4198400,
      "name": "entry",
      "enabled": true
    },
    {
      "address": 4200448,
      "name": "function",
      "enabled": false
    }
  ],
  /* Memory regions
     mm:required
  */
  "memory": [
    {
      "base": 4194304,
      "size": 4096,
      "name": "test.exe"
    }
  ]
}`

	sc := scanner.New(jsoncInput)
	tokens := sc.ScanAll()
	p := parser.New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	t.Logf("Parsed complex JSONC: %s", jsonc.ToJSONC(node))
}

func TestX64RegisterState(t *testing.T) {
	type RegisterState struct {
		Rax uint64 `json:"rax"`
		Rbx uint64 `json:"rbx"`
		Rip uint64 `json:"rip"`
		Rsp uint64 `json:"rsp"`
		Rbp uint64 `json:"rbp"`
	}

	regs := RegisterState{
		Rax: 0x123456789ABCDEF0,
		Rip: 0x00401000,
		Rsp: 0x7FFFFFFF0000,
		Rbp: 0x00123456,
	}

	b, err := json.Marshal(regs)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	var decoded RegisterState
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if decoded.Rax != regs.Rax {
		t.Errorf("Rax mismatch: got %x, want %x", decoded.Rax, regs.Rax)
	}
	if decoded.Rip != regs.Rip {
		t.Errorf("Rip mismatch: got %x, want %x", decoded.Rip, regs.Rip)
	}

	t.Logf("Register state JSON: %s", string(b))
	t.Logf("Rax=0x%X, Rip=0x%X, Rsp=0x%X", decoded.Rax, decoded.Rip, decoded.Rsp)
}

func TestX64MemoryRegion(t *testing.T) {
	type MemoryRegion struct {
		Base    int    `json:"base"`
		Size    int    `json:"size"`
		Name    string `json:"name"`
		Content []byte `json:"content"`
	}

	mem := MemoryRegion{
		Base:    0x400000,
		Size:    4096,
		Name:    "test.exe",
		Content: []byte{0x4D, 0x5A, 0x90, 0x00},
	}

	b, err := json.Marshal(mem)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	var decoded MemoryRegion
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if decoded.Base != mem.Base {
		t.Errorf("Base mismatch: got %x, want %x", decoded.Base, mem.Base)
	}
	if len(decoded.Content) != len(mem.Content) {
		t.Errorf("Content length mismatch: got %d, want %d", len(decoded.Content), len(mem.Content))
	}

	t.Logf("Memory region JSON: %s", string(b))
}

func TestX64ThreadInfo(t *testing.T) {
	type ThreadInfo struct {
		TID    int    `json:"tid"`
		Handle int    `json:"handle"`
		Name   string `json:"name"`
		Entry  int    `json:"entry"`
	}

	threads := []ThreadInfo{
		{TID: 1234, Handle: 0xDEADBEEF, Name: "MainThread", Entry: 0x401000},
		{TID: 5678, Handle: 0xCAFEBABE, Name: "Worker", Entry: 0x402000},
	}

	b, err := json.Marshal(threads)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	var decoded []ThreadInfo
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if len(decoded) != len(threads) {
		t.Errorf("Thread count mismatch: got %d, want %d", len(decoded), len(threads))
	}

	t.Logf("Threads JSON: %s", string(b))
}

func BenchmarkX64DebugStateEncode(b *testing.B) {
	state := X64DebugState{
		Running:     true,
		Pid:         1234,
		ThreadCount: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromValue(state, "")
	}
}

func BenchmarkX64DebugStateDecode(b *testing.B) {
	state := X64DebugState{
		Running:     true,
		Pid:         1234,
		ThreadCount: 1,
	}

	bs, _ := FromValue(state, "")

	var decoded X64DebugState
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node, _ := Decode(bs)
		Bind(node, &decoded)
	}
}
