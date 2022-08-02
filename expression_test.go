package formula_engine

import (
	"fmt"
	"testing"
)

// fakeCell 模拟真实单元格 获取相关环境变量
type fakeCell struct {
	data map[string]string
}

func (fc fakeCell) GetEnvValue(str string) interface{} {
	if fc.data == nil {
		return nil
	}
	if v, ok := fc.data[str]; ok == false {
		return nil
	} else {
		return v
	}
}

func TestNewExpression2(t *testing.T) {
	// exp, err := NewExpression(" REPLACEB(XCELLENTITY, INT(1 + 2), 1, REPLACEB(HELLO, 1, 2, S))")
	exp, err := NewExpression("LTE(-10,2)")
	if err != nil {
		fmt.Println(err)
	}
	// 模拟单元格中的变量
	fc := &fakeCell{
		data: map[string]string{
			"XCELLENTITY": "EH1010",
		},
	}
	result, err := exp.Invoke(fc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result is:%v \n", result)

}
