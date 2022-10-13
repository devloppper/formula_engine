package formula_engine

import (
	"fmt"
	"math/big"
)

// GetInterfaceStringValue 获取Interface的String值
func GetInterfaceStringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%0.8f", v)
	case string:
		return fmt.Sprintf("%s", v)
	case []uint8:
		vInner, _ := v.([]uint8)
		byteList := make([]byte, len(vInner))
		for index, uintV := range vInner {
			byteList[index] = uintV
		}
		return string(byteList)
	case *big.Rat:
		floatV, _ := v.(*big.Rat).Float64()
		return fmt.Sprintf("%f", floatV)
	default:
		return fmt.Sprintf("%v", v)
	}
}
