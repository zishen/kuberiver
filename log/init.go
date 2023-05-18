// Package zklog provides the capability of processing log rules.
package zklog

import (
	"fmt"
	"github.com/zishen/kuberiver/config"
	"github.com/zishen/kuberiver/util"
	"os"
	"path"
)

const (
	// DefaultFileMaxSize  the default maximum size of a single log file is 20 MB
	DefaultFileMaxSize = 20
	// DefaultMinSaveAge the minimum storage duration of backup logs is 7 days
	DefaultMinSaveAge = 7
	// DefaultMaxSaveAge the maximum storage duration of backup logs is 700 days
	DefaultMaxSaveAge = 700
	// DefaultMaxBackups the default number of backup log
	DefaultMaxBackups = 30
	// LogFileMode log file mode
	LogFileMode os.FileMode = 0640
	// LogDirMode log dir mode
	LogDirMode            = 0750
	bitsize               = 64
	stackDeep             = 3
	pathLen               = 2
	minLogLevel           = -1
	maxLogLevel           = 3
	maxEachLineLen        = 1024
	defaultMaxEachLineLen = 256
)

// LogConfig log module config
type LogConfig struct {
	// log file path
	LogFileName string
	// only write to std out, default value: false
	OnlyToStdout bool
	// only write to file, default value: false
	OnlyToFile bool
	// log level, -1-debug, 0-info, 1-warning, 2-error 3-critical default value: 0
	LogLevel int
	// size of a single log file (MB), default value: 20MB
	FileMaxSize int
	// MaxLineLength Max length of each log line, default value: 256
	MaxLineLength int
	// maximum number of backup log files, default value: 30
	MaxBackups int
	// maximum number of days for backup log files, default value: 7
	MaxAge int
	// whether backup files need to be compressed, default value: false
	IsCompress bool
	// expiration time for log cache, default value: 1s
	ExpiredTime int
	// Size of log cache space, default: 10240
	CacheSize int
}

var kLogConfig = &LogConfig{LogFileName: config.LogFile, ExpiredTime: DefaultExpiredTime,
	CacheSize: DefaultCacheSize, MaxBackups: 10, MaxAge: 10, LogLevel: -1}

type validateFunc func(config *LogConfig) error

func checkDir(fileDir string) error {
	if !util.IsExist(fileDir) {
		if err := os.MkdirAll(fileDir, LogDirMode); err != nil {
			return fmt.Errorf("create dirs failed")
		}
		return nil
	}
	if err := os.Chmod(fileDir, LogDirMode); err != nil {
		return fmt.Errorf("change log dir mode failed")
	}
	return nil
}

func createFile(filePath string) error {
	fileName := path.Base(filePath)
	if !util.IsExist(filePath) {
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, LogFileMode)
		if err != nil {
			return fmt.Errorf("create file(%s) failed", fileName)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("close file failed: %v\n", err)
				return
			}
		}()
	}
	return nil
}

func checkAndCreateLogFile(filePath string) error {
	if !util.IsFile(filePath) {
		return fmt.Errorf("config path is not file")
	}
	fileDir := path.Dir(filePath)
	if err := checkDir(fileDir); err != nil {
		return err
	}
	if err := createFile(filePath); err != nil {
		return err
	}
	return nil
}

func validateLogConfigFileMaxSize(config *LogConfig) error {
	if config.FileMaxSize == 0 {
		config.FileMaxSize = DefaultFileMaxSize
		return nil
	}
	if config.FileMaxSize < 0 || config.FileMaxSize > DefaultFileMaxSize {
		return fmt.Errorf("the size of a single log file range is (0, 20] MB")
	}

	return nil
}

func validateLogConfigBackups(config *LogConfig) error {
	if config.MaxBackups <= 0 || config.MaxBackups > DefaultMaxBackups {
		return fmt.Errorf("the number of backup log file range is (0, 30]")
	}
	return nil
}

func validateLogConfigMaxAge(config *LogConfig) error {
	if config.MaxAge < DefaultMinSaveAge || config.MaxAge > DefaultMaxSaveAge {
		return fmt.Errorf("the maxage of backup logs range is [7,700]")
	}
	return nil
}

func validateLogLevel(config *LogConfig) error {
	if config.LogLevel < minLogLevel || config.LogLevel > maxLogLevel {
		return fmt.Errorf("the log level range should be [-1, 3]")
	}
	return nil
}

func validateMaxLineLength(config *LogConfig) error {
	if config.MaxLineLength == 0 {
		config.MaxLineLength = defaultMaxEachLineLen
		return nil
	}
	if config.MaxLineLength < 0 || config.MaxLineLength > maxEachLineLen {
		return fmt.Errorf("the max length of each log line should be in the range (0, 1024]")
	}
	return nil
}

func getValidateFuncList() []validateFunc {
	var funcList []validateFunc
	funcList = append(funcList, validateLogConfigFileMaxSize, validateLogConfigBackups, validateMaxLineLength,
		validateLogConfigMaxAge, validateLogLevel, validateLogConfigLimiter)
	return funcList
}

func validateLogConfigFiled(config *LogConfig) error {
	if config.OnlyToStdout {
		return nil
	}
	if _, err := util.CheckPath(config.LogFileName); err != nil && err != os.ErrNotExist {
		return fmt.Errorf("config log path is not absolute path: %v", err)
	}

	if err := checkAndCreateLogFile(config.LogFileName); err != nil {
		return err
	}
	validateFuncList := getValidateFuncList()
	for _, vaFunc := range validateFuncList {
		if err := vaFunc(config); err != nil {
			return err
		}
	}

	return nil
}

func validateLogConfigLimiter(config *LogConfig) error {
	if config.ExpiredTime < 0 || config.ExpiredTime > MaxExpiredTime {
		return fmt.Errorf("the expired time of log cache range is [0, 3600], the value 0 disables the limiter")
	}
	if config.CacheSize < 0 || config.CacheSize > MaxCacheSize {
		return fmt.Errorf("the size of log cache range is [0, 102400], the value 0 disables the limiter")
	}
	return nil
}
