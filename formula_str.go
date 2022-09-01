package formula_engine

type fISBLANK struct{}

// fISBLANK(待判空的字符串)
func (f *fISBLANK) invoke(env *wrapper, args ...*Token) (*Token, error) {
	arg := args[0]
	isBlank := arg.getStringValue() == ""
	return newBoolToken(isBlank), nil
}
