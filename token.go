package formula_engine

import (
	"fmt"
	"github.com/shopspring/decimal"
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

	Variable
)

// Token 元素
type Token struct {
	TokenType
	IsArray   bool
	Value     interface{}
	ListValue []interface{}
	lockValue bool

	schemaV string
}

// newToken 新建Token
func newToken(str string) *Token {
	str = strings.TrimSpace(str)
	t := &Token{
		schemaV: str,
	}
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
	if strings.HasPrefix(str, ".") == true {
		// .开头的是字符串 不再是数值
		t.TokenType = String
		t.Value = str
		return t
	}
	// 先判断是否是整数值
	if intValue, err := strconv.ParseInt(str, 10, 64); err == nil {
		t.TokenType = Integer
		t.Value = intValue
		return t
	}
	if v, err := decimal.NewFromString(str); err == nil {
		t.TokenType = Number
		t.Value = v
		return t
	}
	// 再判断是否为bool值
	if strings.ToUpper(str) == "TRUE" || strings.ToUpper(str) == "FALSE" {
		t.TokenType = Bool
		t.Value = strings.ToUpper(str) == "TRUE"
		return t
	}
	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") && len(str) > 2 {
		t.TokenType = Variable
		t.Value = str[1 : len(str)-1]
	} else {
		t.TokenType = String
		t.Value = str
	}
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

// newStringToken 新建一个String的Token
func newStringToken(v string) *Token {
	return &Token{
		TokenType: String,
		Value:     v,
	}
}

// newArrayToken 新建数组Token
func newArrayToken(tokenType TokenType) *Token {
	return &Token{
		TokenType: tokenType,
		IsArray:   true,
	}
}

// getFloatValue 获取浮点数值
func (t Token) getFloatValue() decimal.Decimal {
	if t.Value == nil {
		return decimal.NewFromFloat(0)
	}
	if n, ok := t.Value.(int64); ok == true {
		return decimal.NewFromInt(n)
	}
	if n, ok := t.Value.(decimal.Decimal); ok == true {
		return n
	}
	return decimal.NewFromFloat(0)
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

// getStringValue 获取字符串值
func (t Token) getStringValue() string {
	if t.TokenType != String && t.schemaV != "" {
		return t.schemaV
	}
	if t.Value == nil {
		return ""
	}
	if v, ok := t.Value.(string); ok == false {
		if v2, ok2 := t.Value.(int64); ok2 == true {
			return fmt.Sprintf("%d", v2)
		}
		if v2, ok2 := t.Value.(float64); ok2 == true {
			return fmt.Sprintf("%f", v2)
		}
	} else {
		return v
	}
	return ""
}

// getStringList 获取字符串列表
func (t Token) getStringList() []string {
	if t.ListValue == nil {
		return nil
	}
	result := make([]string, 0)
	for _, value := range t.ListValue {
		if v, ok := value.(string); ok == false {
			if v2, ok2 := value.(int64); ok2 == true {
				result = append(result, fmt.Sprintf("%d", v2))
			}
			if v2, ok2 := value.(float64); ok2 == true {
				result = append(result, fmt.Sprintf("%f", v2))
			}
		} else {
			result = append(result, v)
		}
	}
	return nil
}

// getBoolValue 获取布尔值
func (t Token) getBoolValue() bool {
	if t.Value == nil {
		return false
	}
	if v, ok := t.Value.(bool); ok == false {
		return false
	} else {
		return v
	}
}

// splitArrayToken 解析数组类型token
func (t *Token) splitArrayToken() []*Token {
	if t.IsArray == false {
		return []*Token{t}
	} else {
		tokens := make([]*Token, 0)
		for _, v := range t.ListValue {
			tokens = append(tokens, &Token{
				TokenType: t.TokenType,
				Value:     v,
				lockValue: t.lockValue,
			})
		}
		return tokens
	}
}

func (t Token) copy() *Token {
	return &Token{
		TokenType: t.TokenType,
		Value:     t.Value,
	}
}

// compareArgType 比较参数类型
func compareArgType(argType string, tokenType TokenType, tokenIsArray bool) bool {
	if strings.HasPrefix(argType, "Array") {
		if tokenIsArray == false {
			return false
		}
		argType = argType[5:]
	}
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
