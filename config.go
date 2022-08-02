package formula_engine

var calcSignDict = map[rune]bool{
	'+': true,
	'-': true,
	'*': true,
	'/': true,
	'%': true,
}

const ConfigPath = "FORMULA_CONFIG_PATH"
