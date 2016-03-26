package utils

import (
	"fmt"
	//"strings"
)

const (
	DEBUG = true
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

	customErrorMessage := fmt.Sprintf("ERROR: %s :(\n", errorMessage)

	// TODO Log this to Syslog

	if DEBUG {
		fmt.Printf(customErrorMessage)
	}

	if getError {
		//return fmt.Errorf(strings.ToLower(customErrorMessage))
		return fmt.Errorf(customErrorMessage)
	}

	return nil

}
