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
	exp, err := NewExpression("{V1}+{V2}")
	if err != nil {
		fmt.Println(err)
	}
	// 模拟单元格中的变量
	s := &scope{
		data: map[string]interface{}{
			"V1": 123.4,
		},
	}
	result, err := exp.Invoke(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result is:%v \n", result)

}

func TestGt(t *testing.T) {
	exp, _ := NewExpression("GT(ATTR_VALUE,202012)")
	w := NewWrapperEnv(nil)
	w.AddEnv("ATTR_VALUE", " ")
	result, _ := exp.Invoke(w)
	fmt.Println(result)
}

func TestE(t *testing.T) {
	expression, err := NewExpression("num1*(-1)")
	if err != nil {
		panic(err)
	}
	env := NewWrapperEnv(nil)
	env.AddEnv("num1", 10.0)
	result, err := expression.Invoke(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
