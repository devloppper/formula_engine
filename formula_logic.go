package formula_engine

import (
	"github.com/shopspring/decimal"
	"sort"
)

type compareType byte

const (
	gte = compareType(iota)
	lte
	lt
	gt
	eq
	neq
)

type fGTE struct{}

func (g fGTE) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], gte)
}

type fLTE struct{}

func (l fLTE) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], lte)
}

type fLT struct{}

func (f fLT) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], lt)
}

type fGT struct{}

func (f fGT) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], gt)
}

type fEQ struct{}

func (f fEQ) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	// 如果其中一个是布尔类型，则就全部用布尔类型
	if args[0].TokenType == Bool || args[1].TokenType == Bool {
		return newBoolToken(args[0].getBoolValue() == args[1].getBoolValue()), nil
	}
	return compare(args[0], args[1], eq)
}

type fNEQ struct{}

func (f fNEQ) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	// 如果其中一个是布尔类型，则就全部用布尔类型
	if args[0].TokenType == Bool || args[1].TokenType == Bool {
		return newBoolToken(args[0].getBoolValue() != args[1].getBoolValue()), nil
	}
	return compare(args[0], args[1], neq)
}

type fIF struct{}

func (f fIF) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	if args[0].getBoolValue() == true {
		return args[1].copy(), nil
	}
	return args[2].copy(), nil
}

type fMin struct{}

func (f fMin) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	tokens := make([]*Token, 0)
	for _, token := range args {
		if token.IsArray == false {
			tokens = append(tokens, token)
		} else {
			tokens = append(tokens, token.splitArrayToken()...)
		}
	}
	sort.Slice(tokens, func(i, j int) bool {
		token, err := compare(tokens[i], tokens[j], lt)
		if err != nil {
			return false
		}
		return token.getBoolValue()
	})
	return tokens[0], nil
}

type fAnd struct{}

func (f fAnd) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	for _, arg := range args {
		if arg.getBoolValue() == false {
			return newBoolToken(false), nil
		}
	}
	return newBoolToken(true), nil
}

type fOr struct{}

func (f fOr) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	var result = false
	for _, arg := range args {
		result = result || arg.getBoolValue()
	}
	return newBoolToken(result), nil
}

// compare 比较
func compare(v1, v2 *Token, way compareType) (*Token, error) {
	// 如果其中有一个是String的就全部用String进行判断
	if v1.TokenType == String || v2.TokenType == String {
		vStr1 := v1.getStringValue()
		vStr2 := v2.getStringValue()
		return newBoolToken(compareStr(vStr1, vStr2, way)), nil
	}
	// 数值类型全部转换成Float64进行判断
	v1f := v1.getFloatValue()
	v2f := v2.getFloatValue()
	return newBoolToken(compareNumber(v1f, v2f, way)), nil
}

// compareStr 字符串大小比较
// 根据实际情况可注册不同的比较方式
func compareStr(vStr1, vStr2 string, way compareType) bool {
	switch way {
	case gte:
		return vStr1 >= vStr2
	case gt:
		return vStr1 > vStr2
	case lte:
		return vStr1 <= vStr2
	case lt:
		return vStr1 < vStr2
	case eq:
		return vStr1 == vStr2
	case neq:
		return vStr1 != vStr2
	}
	return false
}

// compareNumber 数值大小比较
func compareNumber(v1BigD, v2BigD decimal.Decimal, way compareType) bool {
	switch way {
	case gte:
		return v1BigD.GreaterThanOrEqual(v2BigD)
	case gt:
		return v1BigD.GreaterThan(v2BigD)
	case lte:
		return v1BigD.LessThanOrEqual(v2BigD)
	case lt:
		return v1BigD.LessThan(v2BigD)
	case eq:
		return v1BigD.Equal(v2BigD)
	case neq:
		return !v1BigD.Equal(v2BigD)
	}
	return false
}
