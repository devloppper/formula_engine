package formula_engine

import (
	"fmt"
	"testing"
)

func TestNewExpression(t *testing.T) {
	str := "(1+2 + 3 * (4,2,3))"
	fmt.Println(len(str) - 1)
	fmt.Println(findLastMatchClosure(str))
}

type scope struct {
	data map[string]interface{}
}

func (s scope) GetEnvValue(str string) interface{} {
	if s.data == nil {
		return nil
	}
	return s.data[str]
}

func TestIsBlank(t *testing.T) {
	str := "ISBLANK(ATTR_VALUE,TRUE)"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	env := scope{
		data: map[string]interface{}{
			"ATTR_VALUE": " 1",
		},
	}
	shouldBeTrue, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	fmt.Println(shouldBeTrue)
	env2 := scope{
		data: map[string]interface{}{
			"ATTR_VALUE": "1",
		},
	}
	shouldBeFalse, err := expression.Invoke(env2)
	if err != nil {
		panic(err)
	}
	fmt.Println(shouldBeFalse)
}

/*
$inside_attr:ENTITY(IS_SAP:EQ(ATTR_VALUE,YD),
IS_SAP:EQ(ATTR_VALUE,YF)
&PARENTH1:NEQ(ATTR_VALUE,TOTAL)
&PARENTH1:NEQ(ATTR_VALUE,ENV_JT))
*/
func TestEq(t *testing.T) {

}

func TestHasSubStr(t *testing.T) {
	str := "HASSUBSTR(a,b)"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	env2 := scope{
		data: map[string]interface{}{
			"a": "hello world",
			"b": "ell",
		},
	}
	shouldBeFalse, err := expression.Invoke(env2)
	if err != nil {
		panic(err)
	}
	fmt.Println(shouldBeFalse)
}

func TestIncludeStr(t *testing.T) {
	str := "INCLUDESTR(ATTR_VALUE,1000,I_ENV_HKI)"
	w2 := NewWrapperEnv(&wrapperEnv{
		data: map[string]interface{}{},
	})
	w2.AddEnv("ATTR_VALUE", "1000")
	expression, _ := NewExpression(str)
	result, _ := expression.Invoke(w2)
	fmt.Println(result)
}

type wrapperEnv struct {
	parentEnv Environment
	data      map[string]interface{}
}

// NewWrapperEnv 新建一个包装环境
func NewWrapperEnv(pEnv Environment) *wrapperEnv {
	return &wrapperEnv{
		parentEnv: pEnv,
		data:      map[string]interface{}{},
	}
}

// AddEnv 添加一个环境
func (we *wrapperEnv) AddEnv(key string, v interface{}) {
	we.data[key] = v
}

// GetEnvValue 获取环境值
func (we wrapperEnv) GetEnvValue(str string) interface{} {
	var value interface{}
	if we.data != nil {
		value = we.data[str]
	}
	if value == nil && we.parentEnv != nil {
		value = we.parentEnv.GetEnvValue(str)
	}
	return value
}
