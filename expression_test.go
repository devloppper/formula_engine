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

type SayHello struct {
}

func (sh SayHello) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	str := args[0].getStringValue()
	return newStringToken(fmt.Sprintf("Hello %s", str)), nil
}

func TestNewExpression2(t *testing.T) {
	wb := &WrapperBuilder{}
	wb.AddFormula(NewFormulaEnv("SAYHELLO", "", []string{fmt.Sprintf("%s[LOCK]", ArgStringType)}, ArgStringType), SayHello{})
	w := wb.Build()
	expression, err := NewExpression("SAYHELLO(My)")
	if err != nil {
		panic(err)
	}
	expression.WithOtherWrapper(w)
	s := &scope{
		data: map[string]interface{}{
			"My": "You",
		},
	}
	invoke, err := expression.Invoke(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(invoke)
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
