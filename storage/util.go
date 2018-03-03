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

// ByteSizeForVal determines how many bytes the given val will take up when
// saved in binary format in the RQL DBMS. Only int/string are currently
// supported but in the future other types may be added.
func ByteSizeForVal(val interface{}) int {
	switch val := val.(type) {
	case int:
		return IntSize
	case string:
		return stringSize(val) + IntSize
	default:
		return 0
	}
}
