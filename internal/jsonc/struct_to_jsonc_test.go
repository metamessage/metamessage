package jsonc

import (
	"fmt"
	"testing"
	"time"
)

// go test -v -run TestStructToJsonc
//
// go test internal/jsonc/*.go -v -run TestStructToJsonc/empty_slice_int
func TestStructToJsonc(t *testing.T) {
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

	type C struct {
		B time.Time `mm:"type=time;desc=地址信息"`
	}

	type A struct {
		B time.Time `mm:"type=time;desc=地址信息"`
		C []*struct {
			C []*struct {
				D time.Time `mm:"type=time;desc=地址信息"`
			}
		}
		E []*[]*[]struct{}
	}

	type Some struct {
		ID       int64    `mm:"name=id; min=1,desc=用户ID"`                // 数值最小值+描述
		Name     string   `mm:"name=name; required,max_len=20,desc=用户名"` // 必填+字符串长度+描述
		Age      int      `mm:"name=age; min=18,max=120,desc=年龄"`        // 数值范围+描述
		IsActive bool     `mm:"name=is_active; default=true,desc=是否激活"`  // 默认值+描述
		Tags     []string `mm:"name=tags; max_items=10,desc=标签列表"`       // 切片元素限制+描述
		Addr     *A
	}

	type Datetime struct {
		Datetime  *time.Time `mm:"type=datetime; location=8"`
		Datetime2 []*[]*time.Time
		// Datetime3       string    `mm:"type=datetime"`
		// Datetime4       time.Time `mm:"type=str"`
		// Date            time.Time `mm:"type=date"`
		// Time            time.Time `mm:"type=time"`
		// DatetimePointer *time.Time
		// DataPointer     *time.Time `mm:"type=date"`
		// TimePointer     *time.Time `mm:"type=time"`
		// TimePointerArr  []*time.Time `mm:"type=time"`
		// TimePointerArr2 []*time.Time   `mm:"type=time"`
		// TimePointerArr3 [][]*time.Time `mm:"type=time"`
		// TimePointerArr4 [][]*time.Time `mm:""`
	}

	// type Datetime struct {
	// 	Datetime time.Time `mm:"type=datetime"`
	// 	// TimePointerArr  *time.Time       `mm:"type=time"`
	// 	// TimePointerArr2 []*time.Time     `mm:"type=time"`
	// 	// TimePointerArr3 [][]time.Time    `mm:"type=time"`
	// 	// TimePointerArr4 [][]*time.Time   `mm:"type=date"`
	// 	// TimePointerArr5 [][]*[]time.Time `mm:"type=time; desc=时间"`
	// 	// TimePointerArr6 [1][2]*[3]string `mm:"size=10; desc=天气"`
	// 	TimePointerArr7 map[string]map[string]*[3]string `mm:"size=10; desc=天气"`
	// }

	// type Datetime struct {
	// 	ArrArrPrtTime [][]*time.Time `mm:"child_type=time"`
	// }

	var now = time.Now().UTC()
	fmt.Println(now)

	testCases := []struct {
		name        string
		in          any
		expectedErr bool
	}{
		{
			name: "user struct",
			in: User{
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
			},
		},
		{
			name: "simple map",
			in: map[string]any{
				"v": 123,
				"s": "abc",
			},
		},
		{
			name: "nested slice",
			in: struct {
				Numbers []int `mm:"min_items=1,max_items=5"`
			}{
				Numbers: []int{1, 2, 3},
			},
		},
		{
			name:        "empty slice",
			in:          []any{},
			expectedErr: true,
		},
		{
			name:        "empty slice int",
			in:          []int{},
			expectedErr: false,
		},
		{
			name:        "empty slice int8",
			in:          []int8{},
			expectedErr: false,
		},
		{
			name:        "empty slice pointer int8",
			in:          []*int8{},
			expectedErr: false,
		},
		{
			name:        "empty slice string",
			in:          []string{},
			expectedErr: false,
		},
		{
			name:        "empty slice pointer string",
			in:          []*string{},
			expectedErr: false,
		},
		{
			name:        "empty slice pointer bool",
			in:          []*bool{},
			expectedErr: false,
		},
		{
			name:        "empty slice slice string",
			in:          [][]string{},
			expectedErr: false,
		},
		{
			name:        "empty slice slice slice string",
			in:          [][][]string{},
			expectedErr: false,
		},
		{
			name:        "empty slice slice uint8",
			in:          [][][]uint8{},
			expectedErr: false,
		},
		{
			name:        "empty slice slice pointer uint8",
			in:          [][][]*uint8{},
			expectedErr: false,
		},
		{
			name:        "empty slice pointer slice uint8",
			in:          [][]*[]uint8{},
			expectedErr: false,
		},
		{
			name:        "empty_array_slice_pointer_slice_uint8",
			in:          [2][]*[]uint8{},
			expectedErr: false,
		},
		{
			name:        "empty_slice_array_pointer_slice_uint8",
			in:          [][2]*[]uint8{},
			expectedErr: false,
		},
		{
			name:        "empty_slice_slice_pointer_array_uint8",
			in:          [][]*[2]uint8{},
			expectedErr: false,
		},
		{
			name:        "empty_slice_slice_pointer_array_map",
			in:          [][]*[2]map[string]string{},
			expectedErr: false,
		},
		{
			name:        "datetime",
			in:          Datetime{},
			expectedErr: false,
		},
		{
			name: "datetime_now",
			in:   Datetime{
				// Datetime: &now,
				// Datetime2:       now,
				// Date:            now,
				// Time:            now,
				// DatetimePointer: &now,
				// DataPointer:     &now,
				// TimePointer:     &now,
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			astNode, err := StructToJSONC(tc.in, "")
			t.Logf("err: %v\n", err)

			if (err != nil) != tc.expectedErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.expectedErr, err)
			}

			if !tc.expectedErr {
				if astNode == nil {
					t.Fatalf("expected astNode non-nil")
				}
				fmt.Println("=== AST Node ===")
				fmt.Println("AST Node", Json(astNode))

				// fmt.Println("\n=== go struct ===")
				// PrintGoStruct(astNode)

				fmt.Println("\n=== Print ===")
				Print(astNode)
			}

		})
	}
}

func TestStructToJsonc_InvalidInput(t *testing.T) {
	type Recursive struct {
		Self *Recursive
	}
	r := &Recursive{}
	r.Self = r

	_, err := StructToJSONC(r, "")
	if err == nil {
		t.Error("Expected error for recursive struct, got nil")
	}
}
