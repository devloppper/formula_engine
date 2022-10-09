package formula_engine

import (
	"fmt"
	"testing"
)

func TestIf(t *testing.T) {
	str := "IF(EQ(LEFT($fs, 2), DE), LEFT($fs, 6), IF(EQ(LEFT($fs,2), BE), LEFT($fs, 7),$fs))"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	env := NewWrapperEnv(nil)
	env.AddEnv("$fs", "AE12345678")
	result, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", result)
}
