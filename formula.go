package formula_engine

import (
	"errors"
	"fmt"
)

var ignoreFormulaParamCalc = map[string]bool{
	"CONVERTSTR": true,
}

var formulaDict = map[string]formula{
	"REPLACEB":    &fREPLACEB{},
	"INT":         &fINT{},
	"GTE":         &fGTE{},
	"LTE":         &fLTE{},
	"GT":          &fGT{},
	"LT":          &fLT{},
	"EQ":          &fEQ{},
	"NEQ":         &fNEQ{},
	"LIKE":        &fLike{},
	"ISBLANK":     &fISBLANK{},
	"HASSUBSTR":   &fHasSubStr{},
	"INCLUDESTR":  &fINCLUDESTR{},
	"NINCLUDESTR": &fNINCLUDESTR{},
	"MID":         &fMid{},
	"LEFT":        &fLeft{},
	"RIGHT":       &fRight{},
	"IF":          &fIF{},
	"CONVERTSTR":  &fConvert{},
	"LEN":         &fLen{},
}

// formula 公式
type formula interface {
	// Invoke 计算公式值
	Invoke(env *Wrapper, args ...*Token) (*Token, error)
}

type fREPLACEB struct{}

// Invoke REPLACEB（原字符串，开始位置，字节个数，新字符串)
func (r fREPLACEB) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
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
		return newStringToken(str), nil
	}
	prefix := str[0:startPos]
	if startPos+count >= len(str) {
		return newStringToken(prefix + targetStr), nil
	}
	suffix := str[startPos+count:]
	return newStringToken(prefix + targetStr + suffix), nil
}

type fINT struct{}

func (i fINT) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	arg := args[0]
	result := arg.getFloatValue().IntPart()
	return newIntToken(result), nil
}
