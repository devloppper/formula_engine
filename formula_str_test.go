package formula_engine

import (
	"fmt"
	"testing"
)

func Test_fAdd_Invoke(t *testing.T) {
	exp := "CONCAT(ADD(LEFT({{B}},4),1),.12)"
	env := scope{
		data: map[string]interface{}{
			"{B}": "2021.TOTAL",
		},
	}
	expression, err := NewExpression(exp)
	if err != nil {
		panic(err)
	}
	invoke, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	fmt.Println(invoke)
}
