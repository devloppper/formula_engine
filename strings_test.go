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
	str := "NINCLUDESTR(He, A, B, C,HeA)"
	expression, _ := NewExpression(str)
	result, _ := expression.Invoke(nil)
	fmt.Println(result)
}
