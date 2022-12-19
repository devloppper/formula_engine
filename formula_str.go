package formula_engine

import "strings"

type fISBLANK struct{}

// fISBLANK(待判空的字符串)
func (f *fISBLANK) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	arg := args[0]
	strValue := arg.getStringValue()
	if len(args) > 1 {
		if args[1].getBoolValue() == true {
			// 去除前后空格
			strValue = strings.TrimSpace(strValue)
		}
	}
	isBlank := strValue == ""
	return newBoolToken(isBlank), nil
}

// fINCLUDESTR判断字符串是否属于数组
type fINCLUDESTR struct{}

func (f *fINCLUDESTR) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	if len(args) <= 1 {
		return newBoolToken(false), nil
	}
	strDict := map[string]bool{}
	for i := 1; i < len(args); i++ {
		strDict[args[i].getStringValue()] = true
	}
	standardStr := args[0].getStringValue()
	return newBoolToken(strDict[standardStr]), nil
}

// fNINCLUDESTR fINCLUDESTR判断字符串是否不属于数组
type fNINCLUDESTR struct{}

func (f *fNINCLUDESTR) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	if len(args) <= 1 {
		return newBoolToken(true), nil
	}
	strDict := map[string]bool{}
	for i := 1; i < len(args); i++ {
		strDict[args[i].getStringValue()] = true
	}
	standardStr := args[0].getStringValue()
	return newBoolToken(!strDict[standardStr]), nil
}

// fMid MID(字符串,起始位置,取子符串位数) 起始位置包含 小于1按1处理
type fMid struct{}

func (f *fMid) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	str := args[0].getStringValue()
	if str == "" {
		return newStringToken(str), nil
	}
	startPos := args[1].getIntValue()
	if startPos < 1 {
		// 小于1按照1处理
		startPos = 1
	}
	startPos--
	if startPos > len(str)-1 {
		return newStringToken(""), nil
	}
	subLen := args[2].getIntValue()
	if subLen <= 0 {
		return newStringToken(""), nil
	}
	if startPos+subLen >= len(str) {
		return newStringToken(str[startPos:]), nil
	} else {
		return newStringToken(str[startPos : startPos+subLen]), nil
	}
}

// fLeft LEFT(字符串,取子符串位数)
type fLeft struct{}

func (f *fLeft) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	str := args[0].getStringValue()
	if str == "" {
		return newStringToken(""), nil
	}
	subLen := args[1].getIntValue()
	if subLen < 1 {
		return newStringToken(""), nil
	}
	if subLen > len(str) {
		return newStringToken(str), nil
	}
	return newStringToken(str[:subLen]), nil
}

type fRight struct{}

func (f *fRight) invoke(env *Wrapper, args ...*Token) (*Token, error) {
	str := args[0].getStringValue()
	if str == "" {
		return newStringToken(""), nil
	}
	subLen := args[1].getIntValue()
	if subLen < 1 {
		return newStringToken(""), nil
	}
	if subLen > len(str) {
		return newStringToken(str), nil
	}
	return newStringToken(str[len(str)-subLen:]), nil
}
