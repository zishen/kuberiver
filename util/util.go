package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// IsExist check whether the path exists, If the file is a symbolic link, the returned the final FileInfo
func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

// IsDir check whether the path is a directory.
func IsDir(path string) bool {
	if path == "" {
		return false
	}

	if !IsExist(path) {
		return path[len(path)-1:] == "/"
	}
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile check whether the path is a file
func IsFile(path string) bool {
	if path == "" {
		return false
	}
	return !IsDir(path)
}

// IsLexist check whether the path exists, If the file is a symbolic link, the returned FileInfo
// describes the symbolic link
func IsLexist(filePath string) bool {
	_, err := os.Lstat(filePath)
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

// CheckPath  validate given path and return resolved absolute path
func CheckPath(path string) (string, error) {
	if path == "" {
		return path, nil
	}
	origin := path
	for !IsLexist(path) {
		path = filepath.Dir(path)
		if path == "." {
			return "", os.ErrNotExist
		}
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", errors.New("get the absolute path failed")
	}
	resoledPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return "", os.ErrNotExist
		}
		return "", errors.New("get the symlinks path failed")
	}
	if absPath != resoledPath {
		return "", errors.New("can't support symlinks")
	}
	// get the original full path
	absOrigin, err := filepath.Abs(origin)
	if err != nil {
		return "", errors.New("get the absolute path failed")
	}
	return absOrigin, nil
}
