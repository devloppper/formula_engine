package excel

import (
	"errors"
	"fmt"
)

var formulaDict = map[string]formula{
	"REPLACEB": &fREPLACEB{},
	"INT":      &fINT{},
}

// formula 公式
type formula interface {
	// 计算公式值
	invoke(env *wrapper, args ...*Token) (*Token, error)
}

type fREPLACEB struct{}

// REPLACEB（原字符串，开始位置，字节个数，新字符串)
func (r fREPLACEB) invoke(env *wrapper, args ...*Token) (*Token, error) {
	result := newToken("")
	if args[0].Value == nil {
		return result, nil
	}
	startPos := args[1].getIntValue() - 1
	if startPos < 0 {
		return nil, errors.New(fmt.Sprintf("formula REPLACEB start position should be greater than 0"))
	}
	count := args[2].getIntValue()
	str := fmt.Sprintf("%v", args[0].Value)
	targetStr := fmt.Sprintf("%v", args[3].Value)
	if startPos > len(str) {
		return newToken(str), nil
	}
	prefix := str[0:startPos]
	if startPos+count >= len(str) {
		return newToken(prefix + targetStr), nil
	}
	suffix := str[startPos+count:]
	return newToken(prefix + targetStr + suffix), nil
}

type fINT struct{}

func (i fINT) invoke(env *wrapper, args ...*Token) (*Token, error) {
	arg := args[0]
	result := int64(arg.getFloatValue())
	return newIntToken(result), nil
}
