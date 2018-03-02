package storage

import "os"

// Exists returns a bool of wether or not a path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// stringSize returns the number of bytes to allocate for this string
func stringSize(s string) int {
	return len(s)
}
