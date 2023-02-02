package common

import (
	"errors"
	"os"
)

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func IsDigitsOnly(s string) bool {
	isDot := false
	for _, c := range s {
		if c == '.' {
			if isDot {
				return false
			}
			isDot = true
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
