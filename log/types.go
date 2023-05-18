// Package zklog provides the capability of processing log rules.
package zklog

import "errors"

// ContextKey especially for context value
// to solve problem of "should not use basic type untyped string as key in context.WithValue"
type ContextKey string

// String  the implement of String method
func (c ContextKey) String() string {
	return string(c)
}

const (
	// UserID used for context value key of "ID"
	UserID ContextKey = "UserID"
	// ReqID used for context value key of "requestID"
	ReqID ContextKey = "RequestID"
)

// SelfLogWriter used this to replace some opensource log
type SelfLogWriter struct {
}

// Write  implement the interface of io.writer
func (l *SelfLogWriter) Write(p []byte) (int, error) {
	if RunLog == nil {
		return -1, errors.New("hwlog is not initialized")
	}
	RunLog.Info(string(p))
	return len(p), nil
}
