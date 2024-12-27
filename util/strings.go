// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"slices"
	"strings"
)

func SafeErrStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func OptionalString(val, def string) string {
	if len(val) != 0 {
		return val
	}
	return def
}

func NilOrEmptyString(str *string) bool {
	if str == nil {
		return true
	}
	return len(*str) == 0
}

func SplitAndTrim(s, sep string) []string {
	split := strings.Split(s, sep)
	for i, v := range split {
		split[i] = strings.TrimSpace(v)
	}
	return split
}

var trueVals = []string{"yes", "true", "y", "t", "1"}
var falseVals = []string{"no", "false", "n", "f", "0"}

func IsBool(val string) bool {
	return IsTrue(val) || IsFalse(val)
}

func IsTrue(val string) bool {
	return slices.Contains(trueVals, strings.ToLower(val))
}

func IsFalse(val string) bool {
	return slices.Contains(falseVals, strings.ToLower(val))
}
