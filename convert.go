package formula_engine

type fConvert struct{}

// Invoke CONVERTSTR 强转公式
func (c fConvert) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	result := ""
	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch arg.TokenType {
		case Sub:
			result += "-"
		case Add:
			result += "+"
		case Mul:
			result += "*"
		case Div:
			result += "/"
		case Mod:
			result += "%"
		case RightBracket:
			result += "("
		case LeftBracket:
			result += ")"
		case Separator:
			result += ","
		default:
			result += GetInterfaceStringValue(arg.Value)
		}
	}
	return newStringToken(result), nil
}
