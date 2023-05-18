// Package zklog provides the capability of processing log rules.
package zklog

import (
	"fmt"
	"sync"
	"time"
)

const (
	// MaxCacheSize indicates the maximum log cache size
	MaxCacheSize = 100 * 1024
	// MaxExpiredTime indicates the maximum log cache expired time
	MaxExpiredTime = 60 * 60
	// DefaultCacheSize indicates the default log cache size
	DefaultCacheSize = 100 * 1024
	// DefaultExpiredTime indicates the default log cache expired time
	DefaultExpiredTime = 1
)

// LogLimiter encapsulates Logs and provides the log traffic limiting capability
// to prevent too many duplicate logs.
type LogLimiter struct {
	// Logs is a log rotate instance
	Logs   *Logs
	logMu  sync.Mutex
	doOnce sync.Once

	logExpiredTime time.Duration
	// CacheSize indicates the size of log cache
	CacheSize int
	// ExpiredTime indicates the expired time of log cache
	ExpiredTime int
}

// Write implements io.Writer. It encapsulates the Write method of Los and uses
// the lru cache to prevent duplicate log writing.
func (l *LogLimiter) Write(d []byte) (int, error) {
	if l == nil {
		return 0, fmt.Errorf("log limiter pointer does not exist")
	}

	l.logMu.Lock()
	defer l.logMu.Unlock()

	if l.ExpiredTime == 0 || l.CacheSize == 0 {
		return l.Logs.Write(d)
	}

	l.doOnce.Do(func() {
		l.validateLimiterConf()
		l.logExpiredTime = time.Duration(int64(l.ExpiredTime) * int64(time.Second))
	})

	return l.Logs.Write(d)
}

// Close implements io.Closer. It encapsulates the Close method of Logs.
func (l *LogLimiter) Close() error {
	if l == nil {
		return fmt.Errorf("log limiter pointer does not exist")
	}

	l.logMu.Lock()
	defer l.logMu.Unlock()

	return l.Logs.Close()
}

// Flush encapsulates the Flush method of Logs.
func (l *LogLimiter) Flush() error {
	if l == nil {
		return fmt.Errorf("log limiter pointer does not exist")
	}

	l.logMu.Lock()
	defer l.logMu.Unlock()

	return l.Logs.Flush()
}

// validateLimiterConf verifies the external input parameters in the LogLimiter.
func (l *LogLimiter) validateLimiterConf() {
	if l.CacheSize < 0 || l.CacheSize > MaxCacheSize {
		l.CacheSize = DefaultCacheSize
	}
	if l.ExpiredTime < 0 || l.ExpiredTime > MaxExpiredTime {
		l.ExpiredTime = DefaultExpiredTime
	}
}
