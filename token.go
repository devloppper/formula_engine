package formula_engine

import (
	"strconv"
	"strings"
)

type TokenType int

const (
	None = TokenType(iota)
	Number
	Integer
	String
	Bool
	Separator
	LeftBracket
	RightBracket
	Mul
	Sub
	Mod
	Div
	Add
)

// Token 元素
type Token struct {
	TokenType
	Value interface{}
}

// newToken 新建Token
func newToken(str string) *Token {
	str = strings.TrimSpace(str)
	t := &Token{}
	if str == "," {
		t.TokenType = Separator
		return t
	}
	if str == "(" {
		t.TokenType = LeftBracket
		return t
	}
	if str == ")" {
		t.TokenType = RightBracket
		return t
	}
	if str == "+" || str == "-" || str == "*" || str == "/" || str == "%" {
		switch str {
		case "+":
			t.TokenType = Add
		case "-":
			t.TokenType = Sub
		case "*":
			t.TokenType = Mul
		case "/":
			t.TokenType = Div
		case "%":
			t.TokenType = Mod
		}
		return t
	}
	// 先判断是否是整数值
	if intValue, err := strconv.ParseInt(str, 10, 64); err == nil {
		t.TokenType = Integer
		t.Value = intValue
		return t
	}
	if floatValue, err := strconv.ParseFloat(str, 64); err == nil {
		t.TokenType = Number
		t.Value = floatValue
		return t
	}
	// 再判断是否为bool值
	if strings.ToUpper(str) == "TRUE" || strings.ToUpper(str) == "FALSE" {
		t.TokenType = Bool
		t.Value = strings.ToUpper(str) == "TRUE"
		return t
	}
	t.TokenType = String
	t.Value = str
	return t
}

// newIntToken 新建一个整数Token
func newIntToken(v int64) *Token {
	return &Token{
		TokenType: Integer,
		Value:     v,
	}
}

// newBoolToken 新建一个bool的Token
func newBoolToken(v bool) *Token {
	return &Token{
		TokenType: Bool,
		Value:     v,
	}
}

// getFloatValue 获取浮点数值
func (t Token) getFloatValue() float64 {
	if t.Value == nil {
		return 0
	}
	if n, ok := t.Value.(int64); ok == true {
		return float64(n)
	}
	if n, ok := t.Value.(float64); ok == true {
		return n
	}
	return 0
}

// getIntValue 获取整数值
func (t Token) getIntValue() int {
	if t.Value == nil {
		return 0
	}
	if n, ok := t.Value.(int64); ok == true {
		return int(n)
	}
	return 0
}

// compareArgType 比较参数类型
func compareArgType(argType string, tokenType TokenType) bool {
	switch argType {
	case ArgNumberType:
		return tokenType == Number || tokenType == Integer
	case ArgIntegerType:
		return tokenType == Integer
	case ArgStringType:
		return tokenType == String
	case ArgBoolType:
		return tokenType == Bool
	case ArgAnyType:
		return true
	}
	return false
}

// getStr 获取类型字符串描述
func (ty TokenType) getStr() string {
	switch ty {
	case None:
		return "None"
	case Number:
		return "Number"
	case Integer:
		return "Integer"
	case String:
		return "String"
	case Bool:
		return "Bool"
	}
	return "Unknown"
}

// Sign 元算符号
type Sign struct {
}
