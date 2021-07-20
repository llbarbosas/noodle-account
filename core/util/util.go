package util

import (
	"reflect"
	"runtime"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetKeys(mapStr map[string]string) []string {
	keys := make([]string, len(mapStr))

	for k := range mapStr {
		keys = append(keys, k)
	}

	return keys
}
