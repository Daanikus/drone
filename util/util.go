package util

import "os"

func RemoveDir(path string) {
	os.RemoveAll(path)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !stat.IsDir() {
		return false
	}
	return true
}
