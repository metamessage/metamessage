package pkg

import (
	"fmt"
	"math/big"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/metamessage/metamessage/internal/gen"
	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/jsonc/ast"
)

//go test -v -run TestEncodeDecode
//
//go test -v -run TestEncodeDecode/pointer
//
//go test -bench=BenchmarkEncodeDecode -benchmem

type encodeDecodeTestCase struct {
	name        string
	input       any
	expectedOut any
	expectedErr string
}

func TestEncodeDecode(t *testing.T) {
	type Datetime struct {
		Datetime  []*time.Time `mm:"type=datetime; location=8"`
		Datetime2 time.Time
		// Datetime3       string    `mm:"type=datetime"`
		// Datetime4       time.Time `mm:"type=str"`
		// Date      time.Time `mm:"type=date; location=8"`
		// Time      time.Time `mm:"type=time; location=8"`
		// Undefined string    `mm:"type=underfined"`
		// DatetimePointer *time.Time
		// DataPointer     *time.Time `mm:"type=date"`
		// TimePointer     *time.Time `mm:"type=time"`
		// TimePointerArr  []*time.Time `mm:"type=time"`
		// TimePointerArr2 []*time.Time   `mm:"type=time"`
		// TimePointerArr3 [][]*time.Time `mm:"type=time"`
		// TimePointerArr4 [][]*time.Time `mm:""`
	}
	var now = time.Now().UTC()
	testCases := []encodeDecodeTestCase{
		{
			name: "datetime",
			input: Datetime{
				Datetime2: now,
			},
			expectedErr: "",
		},
		// {
		// 	name:        "pointer2",
		// 	input:       &y,
		// 	expectedOut: []byte{},
		// 	expectedErr: "",
		// },
		// {
		// 	name:        "Ordinary byte slice",
		// 	input:       []byte("hello world"),
		// 	expectedOut: []byte("hello world"),
		// 	expectedErr: "",
		// },
		// {
		// 	name:        "Empty slice ([]byte{})",
		// 	input:       []byte{},
		// 	expectedOut: []byte{},
		// 	expectedErr: "",
		// },
		// {
		// 	name:        "slice",
		// 	input:       []byte{0, 0, 0, 0},
		// 	expectedOut: []byte{0, 0, 0, 0},
		// 	expectedErr: "",
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bs []byte
			bs, err := EncodeFromStruct(tc.input, "")

			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Errorf("Expected error: %s, actual error: %v", tc.expectedErr, err)
				}
				return
			}
			fmt.Println("bs", bs)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			rs, _ := DecodeToJSONC(bs)
			fmt.Println(rs)
			// switch rs.GetType() {
			// case ast.NodeTypeArray:
			// case ast.NodeTypeObject:
			// case ast.NodeTypeValue:
			// 	if !reflect.DeepEqual(rs.(*ast.Value).Data, tc.expectedOut) {
			// 		t.Errorf("Expected output: %v %T, actual output: %v %T", tc.expectedOut, tc.expectedOut, bs2.(*ast.Value).Data, bs2.(*ast.Value).Data)
			// 	}
			// }
		})
	}
}

func BenchmarkEncodeDecode(b *testing.B) {
	data := []byte("benchmark test data 1234567890")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, _ := EncodeFromStruct(data, "")
		_, _ = DecodeToJSONC(n)
	}
}

func TestGenerateGoBasic(t *testing.T) {
	type User struct {
		String   string
		Bytes    []byte `mm:"type=bytes"`
		DateTime time.Time
		Uuid     string
		Ip       net.IP
		Url      url.URL
		Email    string `mm:"type=email"`
		Enum     string `mm:"enum=email|phone"`
		Int      int
		Int8     int8
		Bigint   big.Int
		Decimal  string
		Bool     bool
	}

	user := User{
		Int8: 8,
	}

	astNode, err := jsonc.StructToJSONC(user, "")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println("json", jsonc.Json(astNode))
	fmt.Println("=== astNode res ===")
	fmt.Println(gen.ToGo(astNode))
}

func TestGenerateGoStruct1(t *testing.T) {
	objectNode1 := &ast.Object{
		Tag: &ast.Tag{
			Name: "user_info",
			Type: ast.ValueTypeUnknown,
		},
		Fields: []*ast.Field{
			{
				Key: "user_name",
				Value: &ast.Array{
					Tag: &ast.Tag{
						Name: "ages",
						Type: ast.ValueTypeInt8,
					},
					Items: []ast.Node{
						&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
						&ast.Object{
							Tag: &ast.Tag{
								Name: "user_info1",
								Type: ast.ValueTypeUnknown,
							},
							Fields: []*ast.Field{
								{
									Key:   "user_name",
									Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "zhangsan"},
								},
								{
									Key:   "age",
									Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
								},
							},
						},
					},
				},
			},
			{
				Key:   "age",
				Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
			},
		},
	}

	fmt.Println("=== 嵌套 Object 生成结果 ===")
	fmt.Println(gen.ToGo(objectNode1))

	// 1. 定义测试结构体（包含嵌套、数组、tag）
	type Address struct {
		Province string `mm:"required; max_len=20; desc=省份"`
		City     string `mm:"required; max_len=20; desc=城市"`
		ZipCode  string `mm:"pattern=^\\d{6}$; desc=邮政编码"`
	}

	type User struct {
		ID       int64    `mm:"name=id,min=1,desc=用户ID"`      // 数值最小值+描述
		Name     string   `mm:"required,max_len=20,desc=用户名"` // 必填+字符串长度+描述
		Age      int      `mm:"min=18,max=120,desc=年龄"`       // 数值范围+描述
		IsActive bool     `mm:"default=true,desc=是否激活"`       // 默认值+描述
		Tags     []string `mm:"max_items=10,desc=标签列表"`       // 切片元素限制+描述
		Addr     Address  `mm:"required,desc=地址信息"`           // 必填+描述（嵌套结构体）
	}

	// 2. 创建结构体实例
	user := User{
		ID:       1001,
		Name:     "张三",
		Age:      28,
		IsActive: true,
		Tags:     []string{"golang", "ast", "jsonc"},
		Addr: Address{
			Province: "北京市",
			City:     "朝阳区",
			ZipCode:  "100000",
		},
	}

	// 3. 转换为AST
	astNode, err := jsonc.StructToJSONC(user, "user")
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return
	}
	fmt.Println("json", jsonc.Json(astNode))
	fmt.Println("=== astNode 生成结果 ===")
	fmt.Println(gen.ToGo(astNode))
}
