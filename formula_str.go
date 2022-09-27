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
