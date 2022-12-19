package formula_engine

import (
	"strings"
)

const (
	commonSign    = '%'
	commonSignStr = "%"
)

type fLike struct{}

func (f fLike) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	pattern := args[0].getStringValue()
	origin := args[1].getStringValue()
	if pattern == "" {
		if origin == "" {
			return newBoolToken(true), nil
		} else {
			return newBoolToken(false), nil
		}
	}
	realPattern := strings.ReplaceAll(pattern, commonSignStr, "")
	startPos := 0
	endPos := len(pattern) - 1
	result := true
	if pattern[startPos] == commonSign && pattern[endPos] == commonSign {
		result = strings.Contains(origin, realPattern)
	} else if pattern[startPos] == commonSign {
		index := strings.LastIndex(origin, realPattern)
		result = index == len(origin)-len(realPattern)
	} else if pattern[endPos] == commonSign {
		// 正序匹配字符串
		for index, c := range realPattern {
			if c != rune(origin[index]) {
				result = false
			}
		}
	} else {
		result = pattern == origin
	}
	return newBoolToken(result), nil
}

type fHasSubStr struct{}

func (f fHasSubStr) Invoke(env *Wrapper, args ...*Token) (*Token, error) {
	bigStr := args[0].getStringValue()
	smallStr := args[1].getStringValue()
	return newBoolToken(strings.Contains(bigStr, smallStr)), nil
}

/*

realPattern := strings.ReplaceAll(pattern, CommonSignStr, "")
	startPos := 0
	endPos := len(pattern) - 1
	if pattern[startPos] == CommonSign && pattern[endPos] == CommonSign {
		return strings.Contains(origin, realPattern)
	}
	if pattern[startPos] == CommonSign {
		index := strings.LastIndex(origin, realPattern)
		return index == len(origin)-len(realPattern)
	}
	if pattern[endPos] == CommonSign {
		// 正序匹配字符串
		for index, c := range realPattern {
			if c != rune(origin[index]) {
				return false
			}
		}
		return true
	}
	// 无通配符
	return origin == pattern

*/
