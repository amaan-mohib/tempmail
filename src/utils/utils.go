package utils

import (
	"log"
	"strconv"
)

func ParseInt(num string) int64 {
	result, err := strconv.ParseInt(num, 0, 64)
	if err != nil {
		log.Fatalln(err)
		return 0
	}
	return result
}

func InlineIf[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
