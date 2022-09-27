package formula_engine

import "strings"

type fISBLANK struct{}

// fISBLANK(待判空的字符串)
func (f *fISBLANK) invoke(env *wrapper, args ...*Token) (*Token, error) {
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

func (f *fINCLUDESTR) invoke(env *wrapper, args ...*Token) (*Token, error) {
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

func (f *fNINCLUDESTR) invoke(env *wrapper, args ...*Token) (*Token, error) {
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
