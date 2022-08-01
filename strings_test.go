package excel

import (
	"fmt"
	"testing"
)

func TestNewExpression(t *testing.T) {
	str := "(1+2 + 3 * (4,2,3))"
	fmt.Println(len(str) - 1)
	fmt.Println(findLastMatchClosure(str))
}
