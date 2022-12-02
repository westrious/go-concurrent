package utils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func GoID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("cannot get groutine id: %v", err))
	}
	return int64(id)
}
