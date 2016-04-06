package utils

import (
	"fmt"

	"github.com/andmar/fraudion/types"
)

// StringInStringsSlice ...
func StringInStringsSlice(str string, list []string) bool {
	for _, v := range list {
		fmt.Println(v, str)
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

	customErrorMessage := fmt.Sprintf("ERROR: %s :(", errorMessage)

	// TODO Log this to Syslog

	if types.Globals.Debug {
		fmt.Printf("%s\n", customErrorMessage)
	}

	if getError {
		//return fmt.Errorf(strings.ToLower(customErrorMessage))
		return fmt.Errorf(customErrorMessage)
	}

	return nil

}
