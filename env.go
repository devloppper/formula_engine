package formula_engine

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	ArgNumberType  = "Number"  // 数值类型
	ArgIntegerType = "Integer" // 整数类型
	ArgStringType  = "String"  // 字符串类型
	ArgBoolType    = "Bool"    // 布尔类型
)

var formulaConfigList []*formulaEnv

func init() {
	configBytes, err := os.ReadFile("./formula.json")
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
}

// newWrapper 新建wrapper
func newWrapper(e Environment) *wrapper {
	w := &wrapper{
		env:   e,
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

// wrapper 包装环境
type wrapper struct {
	env   Environment            // 运行环境变量
	fEnv  map[string]*formulaEnv // 公式环境变量字典 map[FORMULA_NAME] -- >
	fDict map[string]formula     // 公式字典
}

// getFormulaEnv 查询公式的环境变量
func (w wrapper) getFormulaEnv(formulaName string) *formulaEnv {
	if w.fEnv == nil {
		return nil
	}
	upFName := strings.ToUpper(formulaName)
	return w.fEnv[upFName]
}

// getFormulaFunc 获取公式字典
func (w wrapper) getFormulaFunc(formulaName string) formula {
	if w.fDict == nil {
		return nil
	}
	upFName := strings.ToUpper(formulaName)
	return w.fDict[upFName]
}
