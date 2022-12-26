package formula_engine

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	ArgNumberType      = "Number"  // 数值类型
	ArgIntegerType     = "Integer" // 整数类型
	ArgStringType      = "String"  // 字符串类型
	ArgBoolType        = "Bool"    // 布尔类型
	ArgAnyType         = "Any"     // 任意类型
	ArgNumberLockType  = "Number"  // 数值类型 锁定
	ArgIntegerLockType = "Integer" // 整数类型 锁定
	ArgStringLockType  = "String"  // 字符串类型 锁定
	ArgBoolLockType    = "Bool"    // 布尔类型 锁定
	ArgAnyLockType     = "Any"     // 任意类型 锁定
)

var formulaConfigList []*formulaEnv

func init() {
	path := os.Getenv(ConfigPath)
	if path == "" {
		path = "./formula.json"
	}
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	if err = json.Unmarshal(configBytes, &formulaConfigList); err != nil {
		return
	}
}

// Environment 环境
type Environment interface {
	GetEnvValue(str string) interface{}
}

// formulaEnv 公式环境变量
type formulaEnv struct {
	FormulaName string   `json:"formula_name"`
	FormulaDesc string   `json:"formula_desc"`
	ArgsType    []string `json:"args_type"`
	ReturnType  string   `json:"return_type"`
	LockDim     []bool   `json:""`
}

// NewFormulaEnv 创建一个公式环境
func NewFormulaEnv(name, desc string, argsType []string, rt string) *formulaEnv {
	return &formulaEnv{
		FormulaName: name,
		FormulaDesc: desc,
		ArgsType:    argsType,
		ReturnType:  rt,
	}
}

// newWrapper 新建Wrapper
func newWrapper(e Environment) *Wrapper {
	w := &Wrapper{
		Env:   e,
		fEnv:  map[string]*formulaEnv{},
		fDict: formulaDict,
	}
	if len(formulaConfigList) > 0 {
		for _, env := range formulaConfigList {
			w.fEnv[strings.ToUpper(env.FormulaName)] = env
		}
	}

	return w
}

// Wrapper 包装环境
type Wrapper struct {
	Env   Environment            // 运行环境变量
	fEnv  map[string]*formulaEnv // 公式环境变量字典 map[FORMULA_NAME] -- >
	fDict map[string]formula     // 公式字典
}

// getFormulaEnv 查询公式的环境变量
func (w Wrapper) getFormulaEnv(formulaName string) *formulaEnv {
	if w.fEnv == nil {
		return nil
	}
	upFName := strings.ToUpper(formulaName)
	return w.fEnv[upFName]
}

// getFormulaFunc 获取公式字典
func (w Wrapper) getFormulaFunc(formulaName string) formula {
	if w.fDict == nil {
		return nil
	}
	upFName := strings.ToUpper(formulaName)
	return w.fDict[upFName]
}

// merge 合并环境变量
// 以w2为准
func (w *Wrapper) merge(w2 *Wrapper) {
	if w2.Env != nil {
		w.Env = w2.Env
	}
	if len(w2.fEnv) > 0 {
		if w.fEnv == nil {
			w.fEnv = map[string]*formulaEnv{}
		}
		for k, v := range w2.fEnv {
			w.fEnv[k] = v
		}
	}
	if len(w2.fDict) > 0 {
		if w.fDict == nil {
			w.fDict = map[string]formula{}
		}
		for k, v := range w2.fDict {
			w.fDict[k] = v
		}
	}
}

type WrapperBuilder struct {
	*Wrapper
}

// AddFormula 添加一个公式
func (wb *WrapperBuilder) AddFormula(fEnv *formulaEnv, f formula) {
	if wb.Wrapper == nil {
		wb.Wrapper = &Wrapper{
			fEnv:  map[string]*formulaEnv{},
			fDict: map[string]formula{},
		}
	}
	fEnv.FormulaName = strings.ToUpper(fEnv.FormulaName)
	wb.Wrapper.fEnv[fEnv.FormulaName] = fEnv
	wb.Wrapper.fDict[fEnv.FormulaName] = f
}

// Build 构建环境包裹
func (wb WrapperBuilder) Build() *Wrapper {
	if wb.Wrapper == nil {
		return &Wrapper{}
	}
	return wb.Wrapper
}
