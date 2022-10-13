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
	str := "ISBLANK(ATTR_VALUE)"
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
	str := "INCLUDESTR(ATTR_VALUE,IL1-1-1-3)"
	w2 := NewWrapperEnv(&wrapperEnv{
		data: map[string]interface{}{},
	})
	w2.AddEnv("ATTR_VALUE", "I_BL2")
	expression, _ := NewExpression(str)
	result, _ := expression.Invoke(w2)
	fmt.Println(result)
	w2.AddEnv("ATTR_VALUE", "I_BL")
	result, _ = expression.Invoke(w2)
	fmt.Println(result)
}

func TestReplaceB(t *testing.T) {
	str := "IF(EQ(1,1), IF(EQ(2,3),TOM, JERRY),Teifi )"
	env := NewWrapperEnv(nil)
	env.AddEnv("ATTR_VALUE", " ")
	expression, _ := NewExpression(str)
	result, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("%v", result))
}

func TestMid(t *testing.T) {
	str := "MID($fs, 2, 10)"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	env := NewWrapperEnv(nil)
	env.AddEnv("$fs", "HELLO WORLD")
	result, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestLEFT(t *testing.T) {
	str := "LEFT($fs, 7)"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	env := NewWrapperEnv(nil)
	env.AddEnv("$fs", "ER-J009-02")
	result, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestRIGHT(t *testing.T) {
	str := "RIGHT($fs, 9)"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	env := NewWrapperEnv(nil)
	env.AddEnv("$fs", "HELLO WORLD")
	result, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestConvert(t *testing.T) {
	str := "CONVERTSTR(INT(IL1-1-1-3))"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	invoke, err := expression.Invoke(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(invoke)
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
