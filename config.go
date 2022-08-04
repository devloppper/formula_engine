package formula_engine

var calcSignDict = map[rune]bool{
	'+': true,
	'-': true,
	'*': true,
	'/': true,
	// '%': true, % 暂不支持
}

const ConfigPath = "FORMULA_CONFIG_PATH"
