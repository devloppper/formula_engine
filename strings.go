package formula_engine

const (
	LeftClosure  = '('
	RightClosure = ')'
)

// findLastMatchClosure 找到匹配的最后一个 ) 的位置
func findLastMatchClosure(str string) int {
	if len(str) <= 0 || str[0] != '(' {
		return 0
	}
	count := 0
	for index, r := range str {
		if r == '(' {
			count++
		}
		if r == ')' {
			count--
			if count == 0 {
				return index
			}
		}
	}
	return -1
}
