// Package zklog provides the capability of processing log rules.
package zklog

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
)

// printHelper helper function for log printing
func printHelper(lg *log.Logger, msg string, maxLogLength int, ctx ...context.Context) {
	str := getCallerInfo(ctx...)
	trimMsg := strings.Replace(msg, "\r", " ", -1)
	trimMsg = strings.Replace(trimMsg, "\n", " ", -1)
	runeArr := []rune(trimMsg)
	if length := len(runeArr); length > maxLogLength {
		trimMsg = string(runeArr[:maxLogLength])
	}
	lg.Println(str + trimMsg)
}

// getCallerInfo gets the caller's information
func getCallerInfo(ctx ...context.Context) string {
	var deep = stackDeep
	var userID interface{}
	var traceID interface{}
	for _, c := range ctx {
		if c == nil {
			deep++
			continue
		}
		userID = c.Value(UserID)
		traceID = c.Value(ReqID)
	}
	var funcName string
	pc, codePath, codeLine, ok := runtime.Caller(deep)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	p := strings.Split(codePath, "/")
	l := len(p)
	if l == pathLen {
		funcName = p[l-1]
	} else if l > pathLen {
		funcName = fmt.Sprintf("%s/%s", p[l-pathLen], p[l-1])
	}
	callerPath := fmt.Sprintf("%s:%d", funcName, codeLine)
	goroutineID := getGoroutineID()
	str := fmt.Sprintf("%-8s%s    ", goroutineID, callerPath)
	if userID != nil || traceID != nil {
		str = fmt.Sprintf("%s{%#v}-{%#v} ", str, userID, traceID)
	}
	return str
}

// getCallerGoroutineID gets the goroutineID
func getGoroutineID() string {
	b := make([]byte, bitsize, bitsize)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}
