package formula_engine

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

const separator = ','

var splitDict = map[rune]bool{
	'+': true,
	'-': true,
	'*': true,
	'/': true,
	'%': true,
}

type Expression struct {
	root *unit
	*wrapper
}

// NewExpression 新建一个运算表达式
// 10
// 10 + 20
// 10 + SUM(1,2)
// SUM(1, SUM(3,4)
func NewExpression(str string) (*Expression, error) {
	str = strings.TrimSpace(str)
	root := &unit{
		Token: newToken("#Value"),
	}
	if err := newUnit(str, root); err != nil {
		return nil, err
	}
	return &Expression{root: root}, nil
}

// WithOtherWrapper 添加额外的配置信息
func (exp *Expression) WithOtherWrapper(w *wrapper) {
	exp.wrapper = w
}

// Invoke 执行
func (exp Expression) Invoke(e Environment) (interface{}, error) {
	w := newWrapper(e)
	if exp.wrapper != nil {
		w.merge(exp.wrapper)
	}
	calc, err := exp.root.calc(w, e)
	if err != nil {
		return nil, err
	}
	if calc == nil || calc.Value == nil {
		return nil, nil
	}
	if calc.TokenType == String && e != nil {
		// 可能是潜在变量
		value := e.GetEnvValue(fmt.Sprintf("%v", calc.Value))
		if value != nil {
			return value, nil
		}
	}
	if calc.TokenType == Number {
		if v, ok := calc.Value.(decimal.Decimal); ok == true {
			return v.String(), nil
		}
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
