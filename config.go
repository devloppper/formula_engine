package formula_engine

var calcSignDict = map[rune]bool{
	'+': true,
	'-': true,
	'*': true,
	'/': true,
	// '%': true, % ζδΈζ―ζ
}

const ConfigPath = "FORMULA_CONFIG_PATH"
