// Package zklog provides the capability of processing log rules.
package zklog

import (
	"errors"
	"fmt"
)

// RunLog run logger
var RunLog *logger

// InitRunLogger initialize run logger
func InitRunLogger(config *LogConfig) error {
	if config == nil {
		return errors.New("run logger config is nil")
	}
	if RunLog != nil && RunLog.isInit() {
		RunLog.Warn("run logger is been initialized.")
		return nil
	}
	RunLog = new(logger)
	if RunLog == nil {
		return errors.New("malloc new logger flied")
	}
	if err := RunLog.setLogger(config); err != nil {
		return err
	}
	if !RunLog.isInit() {
		return errors.New("run logger init failed")
	}
	return nil
}

func InitZKLogger() error {
	if err := InitRunLogger(kLogConfig); err != nil {
		fmt.Printf("hwlog init failed, error is %#v", err)
		return err
	}
	return nil
}
