package util

import "os"

func RemoveDir(path string) {
	os.RemoveAll(path)
}
