package excel

import "strings"

const separator = ','

var splitDict = map[rune]bool{
	'+': true,
	'-': true,
	'*': true,
	'/': true,
	'%': true,
}

type expression struct {
	root *unit
	*wrapper
}

// NewExpression 新建一个运算表达式
// 10
// 10 + 20
// 10 + SUM(1,2)
// SUM(1, SUM(3,4)
func NewExpression(str string) (*expression, error) {
	str = strings.TrimSpace(str)
	root := &unit{
		Token: newToken("#Value"),
	}
	if err := newUnit(str, root); err != nil {
		return nil, err
	}
	return &expression{root: root}, nil
}

// Invoke 执行
func (exp expression) Invoke(e environment) (interface{}, error) {
	w := newWrapper(e)
	calc, err := exp.root.calc(w)
	if err != nil {
		return nil, err
	}
	if calc == nil || calc.Value == nil {
		return nil, nil
	}
	return calc.Value, nil
}

// wrapperBaseFunc 包装BASE函数
// + ( 1 + 2)  ---> + BASE(1 + 2)
func wrapperBaseFunc(exp string) string {
	markPoint := make([]int, 0)
	lastIsCalcSign := false
	for index, r := range exp {
		if lastIsCalcSign == false && calcSignDict[r] == true {
			lastIsCalcSign = true
		}
		if lastIsCalcSign == true && r == '(' {
			markPoint = append(markPoint, index)
		}
	}
	for i := len(markPoint) - 1; i >= 0; i-- {
		tempPrefix := exp[0:i]
		tempSuffix := exp[i:]
		exp = tempPrefix + "BASE" + tempSuffix
	}
	return exp
}
