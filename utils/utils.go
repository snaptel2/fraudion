package utils

import (
	"fmt"

	"github.com/andmar/fraudion/logger"
)

// StringInStringsSlice ...
func StringInStringsSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// StringKeyInMap ...
func StringKeyInMap(theKey string, theMap map[string]interface{}) bool {
	for key := range theMap {
		if key == theKey {
			return true
		}
	}
	return false
}

// DebugLogAndGetError ...
func DebugLogAndGetError(errorMessage string, getError bool) error {

	logger.Log.Write(logger.ConstLoggerLevelError, errorMessage, false)

	if getError {
		return fmt.Errorf(errorMessage)
	}

	return nil

}
