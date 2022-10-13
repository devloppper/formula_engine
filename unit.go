package formula_engine

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

// unit 单元格
type unit struct {
	*unitFormula         // 公式单元 指向一个公式
	*Token               // 值表示
	params       []*unit // 参数
}

// newUnit 新建一个单元
// 实则是给一个Root的Unit添加Param
func newUnit(str string, u *unit) error {
	tmpStr := ""
	continueNum := -1
	for index, r := range str {
		if index <= continueNum {
			continue
		}
		// 肯定是一个运算符号
		if calcSignDict[r] == true && len(tmpStr) > 0 {
			t1 := newToken(tmpStr)
			u.params = append(u.params, &unit{Token: t1})
			t2 := newToken(string(r))
			u.params = append(u.params, &unit{Token: t2})
			tmpStr = ""
			continue
		}

		// 公式
		if r == '(' && len(tmpStr) > 0 {
			lastPos := findLastMatchClosure(str[index:])
			if lastPos <= 0 {
				return errors.New("invalid bracket: not have right bracket")
			}
			fu := &unit{
				unitFormula: &unitFormula{FormulaName: strings.TrimSpace(strings.ToUpper(tmpStr))},
			}
			continueNum = index + lastPos
			tmpStr = ""
			subStr := str[index:][1:lastPos]
			if err := newUnit(subStr, fu); err != nil {
				return err
			} else {
				u.params = append(u.params, fu)
			}
			continue
		}
		if r == '(' || r == ')' || r == ',' {
			if tmpStr != "" {
				t1 := newToken(tmpStr)
				u.params = append(u.params, &unit{Token: t1})
				tmpStr = ""
			}
			t2 := newToken(string(r))
			u.params = append(u.params, &unit{Token: t2})
			continue
		}
		tmpStr = tmpStr + string(r)
	}
	if tmpStr != "" {
		t := newToken(tmpStr)
		u.params = append(u.params, &unit{Token: t})
		tmpStr = ""
	}
	return nil
}

// calc 计算节点值
func (u unit) calc(w *wrapper) (*Token, error) {
	if u.unitFormula != nil {
		exp := make([]*Token, 0)
		if len(u.params) > 0 {
			for _, param := range u.params {
				if ignoreFormulaParamCalc[u.FormulaName] == true {
					exp = append(exp, param.Token)
					continue
				}
				if param.unitFormula != nil {
					tmpResult, err := param.calc(w)
					if err != nil {
						return nil, err
					}
					exp = append(exp, tmpResult)
				} else {
					exp = append(exp, param.Token)
				}
			}
		} else {
			return nil, errors.New(fmt.Sprintf("Formula:%s, in fact, all formulas need args", u.FormulaName))
		}
		// 合并参数
		tempTokens := make([]*Token, 0)
		paramsToken := make([]*Token, 0)
		if ignoreFormulaParamCalc[u.FormulaName] == false {
			for index, t := range exp {
				if t.TokenType == Separator {
					r, err := baseCalc(tempTokens...)
					if err != nil {
						return nil, err
					}
					paramsToken = append(paramsToken, r)
					tempTokens = make([]*Token, 0)
					continue
				}
				tempTokens = append(tempTokens, t)
				if index == len(exp)-1 {
					r, err := baseCalc(tempTokens...)
					if err != nil {
						return nil, err
					}
					paramsToken = append(paramsToken, r)
				}
			}
		} else {
			paramsToken = exp
		}
		// 真实公式计算
		return u.unitFormula.calc(w, paramsToken...)
	}

	if len(u.params) > 0 {
		exp := make([]*Token, 0)
		for _, param := range u.params {
			if param.unitFormula != nil {
				tmpResult, err := param.calc(w)
				if err != nil {
					return nil, err
				}
				exp = append(exp, tmpResult)
			} else {
				exp = append(exp, param.Token)
			}
		}
		if len(exp) == 1 {
			return exp[0], nil
		}
		return baseCalc(exp...)
	}

	if u.Token != nil {
		return u.Token, nil
	}
	return nil, errors.New("unknown unit type")
}

// unitFormula 公式单元
type unitFormula struct {
	FormulaName string // 公式名称
}

// calc 公式计算
func (uf unitFormula) calc(w *wrapper, args ...*Token) (*Token, error) {
	// 查询环境变量
	fEnv := w.getFormulaEnv(uf.FormulaName)
	fFunc := w.getFormulaFunc(uf.FormulaName)
	if fEnv == nil || fFunc == nil {
		return nil, errors.New(fmt.Sprintf("formula:%s not config or implement", uf.FormulaName))
	}
	// 参数检查
	if err := uf.checkParams(fEnv, args...); err != nil {
		return nil, err
	}
	// String类型为q潜在变量
	originDict := map[*Token]interface{}{}
	if w.env != nil {
		for _, arg := range args {
			if arg.TokenType == String {
				originDict[arg] = arg.Value
				v := fmt.Sprintf("%v", arg.Value)
				if w.env.GetEnvValue(v) != nil {
					arg.Value = w.env.GetEnvValue(v)
				}
			}
		}
	}
	// 计算
	result, err := fFunc.invoke(w, args...)
	if err != nil {
		return nil, err
	}
	// 计算完后还原变量
	for t, v := range originDict {
		t.Value = v
	}
	if compareArgType(fEnv.ReturnType, result.TokenType) == false {
		return nil, errors.New(fmt.Sprintf("formula:%s return %s but actual it is %s", uf.FormulaName, fEnv.ReturnType, result.TokenType.getStr()))
	}
	return result, nil
}

// checkParams 检查参数
// 如果需要String的，则将参数类型转为String
func (uf unitFormula) checkParams(fEnv *formulaEnv, args ...*Token) error {
	minArgsCount := 0
	for _, argType := range fEnv.ArgsType {
		if strings.HasPrefix(argType, "...") == true {
			break
		}
		minArgsCount++
	}
	if minArgsCount > len(args) {
		return errors.New(fmt.Sprintf("formula:%s need at least %d args but actual it is %d", uf.FormulaName, len(fEnv.ArgsType), len(args)))
	}
	argIndex := 0
	for _, argType := range fEnv.ArgsType {
		flexArg := strings.HasPrefix(argType, "...")
		if flexArg == true {
			argType = argType[3:]
			if argIndex >= len(args) {
				break
			}
		}
		arg := args[argIndex]
		if flexArg == false {
			if argType == ArgStringType {
				arg.Value = arg.getStringValue()
				arg.TokenType = String
			} else {
				if result := compareArgType(argType, arg.TokenType); result != true {
					return errors.New(fmt.Sprintf("formula:%s need arg type:%s But actual it is %s", uf.FormulaName, argType, arg.TokenType.getStr()))
				}
			}
		} else {
			for i := argIndex; i < len(args); i++ {
				arg = args[i]
				if argType == ArgStringType {
					arg.Value = arg.getStringValue()
					arg.TokenType = String
				} else {
					if result := compareArgType(argType, arg.TokenType); result != true {
						return errors.New(fmt.Sprintf("formula:%s need arg type:%s But actual it is %s", uf.FormulaName, argType, arg.TokenType.getStr()))
					}
				}
			}
			break
		}
		argIndex++
	}
	return nil
}

// baseCalc 基本运算
func baseCalc(tokens ...*Token) (*Token, error) {
	stacks := map[int]*stack{}
	currentStackIndex := 0
	stacks[currentStackIndex] = NewStack()
	for _, t := range tokens {
		if t.TokenType == LeftBracket {
			currentStackIndex++
			stacks[currentStackIndex] = NewStack()
			continue
		}
		if t.TokenType == RightBracket {
			// 计算stacks
			result, err := calcStack(stacks[currentStackIndex])
			if err != nil {
				return nil, err
			}
			delete(stacks, currentStackIndex)
			currentStackIndex--
			stacks[currentStackIndex].Push(result)
			continue
		}
		stacks[currentStackIndex].Push(t)
	}
	if len(stacks) > 1 {
		return nil, errors.New("expression has some error")
	}
	return calcStack(stacks[0])
}

// calcStack 计算栈中元素
func calcStack(s *stack) (*Token, error) {
	s1 := NewStack()
	// 先算高优先级 * / %
	for s.len() > 0 {
		t := s.Pop()
		if t.TokenType == Mul || t.TokenType == Mod || t.TokenType == Div {
			lt := s.Pop()
			nt := s1.Pop()
			if (lt.TokenType != Number && lt.TokenType != Integer) || (nt.TokenType != Number && nt.TokenType != Integer) {
				return nil, errors.New("* / % need two Number")
			}
			tResult, err := calcToken(lt, nt, t)
			if err != nil {
				return nil, err
			}
			s.Push(tResult)
			continue
		}
		s1.Push(t)
	}
	s2 := NewStack()
	// 再算低优先级 + -
	for s1.len() > 0 {
		t := s1.Pop()
		if t.TokenType == Add || t.TokenType == Sub {
			lt := s1.Pop()
			nt := s2.Pop()
			if (lt.TokenType != Number && lt.TokenType != Integer) || (nt.TokenType != Number && nt.TokenType != Integer) {
				return nil, errors.New("* / % need two Number")
			}
			tResult, err := calcToken(lt, nt, t)
			if err != nil {
				return nil, err
			}
			s1.Push(tResult)
			continue
		}
		s2.Push(t)
	}
	return s2.Pop(), nil
}

// calcToken 计算
func calcToken(s1, s2, sign *Token) (*Token, error) {
	r := &Token{}
	switch sign.TokenType {
	case Mul:
		r.TokenType = Number
		n1 := decimal.NewFromFloat(s1.getFloatValue())
		n2 := decimal.NewFromFloat(s2.getFloatValue())
		r.Value, _ = n1.Mul(n2).Float64()
		return r, nil
	case Div:
		r.TokenType = Number
		n1 := decimal.NewFromFloat(s1.getFloatValue())
		n2 := decimal.NewFromFloat(s2.getFloatValue())
		r.Value, _ = n1.Div(n2).Float64()
		return r, nil
	case Add:
		r.TokenType = Number
		n1 := decimal.NewFromFloat(s1.getFloatValue())
		n2 := decimal.NewFromFloat(s2.getFloatValue())
		r.Value, _ = n1.Add(n2).Float64()
		return r, nil
	case Sub:
		r.TokenType = Number
		n1 := decimal.NewFromFloat(s1.getFloatValue())
		n2 := decimal.NewFromFloat(s2.getFloatValue())
		r.Value, _ = n1.Sub(n2).Float64()
		return r, nil
	case Mod:
		r.TokenType = Number
		n1 := decimal.NewFromFloat(s1.getFloatValue())
		n2 := decimal.NewFromFloat(s2.getFloatValue())
		r.Value, _ = n1.Mod(n2).Float64()
		return r, nil
	}
	return nil, errors.New("unknown calc type")
}
