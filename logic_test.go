package formula_engine

import (
	"fmt"
	"testing"
)

func TestIf(t *testing.T) {
	//str := "IF(EQ(LEFT($fs, 2), DE), LEFT($fs, 6), IF(EQ(LEFT($fs,2), BE), LEFT($fs, 7),$fs))"
	str := "IF(EQ(1,1),TRUE,FALSE)"
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

func TestMin(t *testing.T) {
	str := "LIKE(LEFT(ABC, LEN(ABC)), Z)"
	expression, err := NewExpression(str)
	if err != nil {
		panic(err)
	}
	env := NewWrapperEnv(nil)
	env.AddEnv("ABC", "HELLO")
	env.AddEnv("HELLO", "Z")

	result, err := expression.Invoke(env)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", result)
	s := "XREALNAME"[len("XREAL"):]
	println(s)
}
