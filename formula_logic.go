package formula_engine

import "github.com/shopspring/decimal"

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

func (g fGTE) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], gte)
}

type fLTE struct{}

func (l fLTE) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], lte)
}

type fLT struct{}

func (f fLT) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], lt)
}

type fGT struct{}

func (f fGT) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], gt)
}

type fEQ struct{}

func (f fEQ) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], eq)
}

type fNEQ struct{}

func (f fNEQ) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	return compare(args[0], args[1], neq)
}

type fIF struct{}

func (f fIF) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	if args[0].getBoolValue() == true {
		return args[1].copy(), nil
	}
	return args[2].copy(), nil
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
