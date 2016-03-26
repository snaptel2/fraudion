package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/fraudion/utils"
)

const (
	constMinimumTriggerCheckPeriod     = "5m"
	constMaximumActionChainSleepPeriod = "5m"
)

var (
	constSupportedSoftware = []string{"*ast_elastix_2.3", "*ast_alone_1.8"}
	constSupportCDRSources = []string{"*db_mysql"}
)

// CheckJSONSanityAndLoadConfigs ...
func (fraudionConfig *FraudionConfig) CheckJSONSanityAndLoadConfigs(ConfigsJSON *FraudionConfigJSON /*, validateOnly bool*/) error {

	// ** General Section
	// * Monitored Software
	if err := errorOnValueNotFound(ConfigsJSON.General.MonitoredSoftware == nil, "General", "MonitoredSoftware", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfMonitoredSoftware, hasCorrectType := ConfigsJSON.General.MonitoredSoftware.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "MonitoredSoftware", "string", reflect.TypeOf(ConfigsJSON.General.MonitoredSoftware)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	found := utils.StringInStringsSlice(valueOfMonitoredSoftware, constSupportedSoftware)
	if err := errorOnIncorrectValue(!found, "General", "MonitoredSoftware", fmt.Sprintf("One of: %s", constSupportedSoftware), valueOfMonitoredSoftware); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.MonitoredSoftware = valueOfMonitoredSoftware

	// * CDRs Source
	if err := errorOnValueNotFound(ConfigsJSON.General.CDRsSource == nil, "General", "CDRsSource", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfCDRsSource, hasCorrectType := ConfigsJSON.General.CDRsSource.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "CDRsSource", "string", reflect.TypeOf(ConfigsJSON.General.CDRsSource)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	found = utils.StringInStringsSlice(valueOfCDRsSource, constSupportCDRSources)
	if err := errorOnIncorrectValue(!found, "General", "CDRsSource", fmt.Sprintf("One of: %s", constSupportCDRSources), valueOfCDRsSource); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.CDRsSource = valueOfCDRsSource

	// * DefaultTriggerCheckPeriod
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultTriggerCheckPeriod == nil, "General", "DefaultTriggerCheckPeriod", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultTriggerCheckPeriod, hasCorrectType := ConfigsJSON.General.DefaultTriggerCheckPeriod.(string)
	fmt.Println(valueOfDefaultTriggerCheckPeriod)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultTriggerCheckPeriod", "string", reflect.TypeOf(ConfigsJSON.General.DefaultTriggerCheckPeriod)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfDefaultTriggerCheckPeriod, errDuration := time.ParseDuration(valueOfDefaultTriggerCheckPeriod)
	if err := errorOnIncorrectValue(errDuration != nil, "General", "DefaultTriggerCheckPeriod", "parseable Duration (e.g. 5\"s\", \"m\" or \"h\")", valueOfDefaultTriggerCheckPeriod); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfConstMinimumTriggerCheckPeriod, errMinimumCheckPeriodDuration := time.ParseDuration(constMinimumTriggerCheckPeriod)
	errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", constMinimumTriggerCheckPeriod)
	if err := errorOnIncorrectValue(durationValueOfDefaultTriggerCheckPeriod < durationValueOfConstMinimumTriggerCheckPeriod, "General", "DefaultTriggerCheckPeriod", errorMessageSuffix, ""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if errMinimumCheckPeriodDuration != nil {
		return utils.DebugLogAndGetError("Internal: There seems to be an issue with the definition of constMinimumTriggerCheckPeriod", true)
	}
	fraudionConfig.General.DefaultTriggerCheckPeriod = durationValueOfDefaultTriggerCheckPeriod

	// DefaultHitThreshold
	// TODO !

	// * DefaultMinimumDestinationNumberLength
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultMinimumDestinationNumberLength == nil, "General", "DefaultActionChainRunCount", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultMinimumDestinationNumberLength, hasCorrectType := ConfigsJSON.General.DefaultMinimumDestinationNumberLength.(float64)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultMinimumDestinationNumberLength", "int", reflect.TypeOf(ConfigsJSON.General.DefaultMinimumDestinationNumberLength)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if err := errorOnIncorrectValue(valueOfDefaultMinimumDestinationNumberLength < 0, "General", "DefaultMinimumDestinationNumberLength", "Must be >= 0", strconv.FormatFloat(valueOfDefaultMinimumDestinationNumberLength, 'f', 6, 64)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.DefaultActionChainRunCount = uint32(valueOfDefaultMinimumDestinationNumberLength)

	// * DefaultActionChainSleepPeriod
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultActionChainSleepPeriod == nil, "General", "DefaultActionChainSleepPeriod", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultActionChainSleepPeriod, hasCorrectType := ConfigsJSON.General.DefaultActionChainSleepPeriod.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultActionChainSleePeriod", "string", reflect.TypeOf(ConfigsJSON.General.DefaultActionChainSleepPeriod)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfdefaultActionChainSleepPeriodDuration, errDuration := time.ParseDuration(valueOfDefaultActionChainSleepPeriod)
	if err := errorOnIncorrectValue(errDuration != nil, "General", "DefaultActionChainSleepPeriod", "parseable Duration (e.g. 5\"s\", \"m\" or \"h\")", valueOfDefaultActionChainSleepPeriod); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfmaximumActionChainSleepPeriodDuration, errMaximumChainSleepPeriod := time.ParseDuration(constMaximumActionChainSleepPeriod)
	errorMessageSuffix = fmt.Sprintf("Maximum value is \"%s\"", constMaximumActionChainSleepPeriod)
	if err := errorOnIncorrectValue(durationValueOfdefaultActionChainSleepPeriodDuration > durationValueOfmaximumActionChainSleepPeriodDuration, "General", "DefaultActionChainSleepPeriod", errorMessageSuffix, ""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if errMaximumChainSleepPeriod != nil {
		return utils.DebugLogAndGetError("Internal: There seems to be an issue with the definition of constMaximumActionChainSleepPeriod", true)
	}
	fraudionConfig.General.DefaultActionChainSleepPeriod = durationValueOfdefaultActionChainSleepPeriodDuration

	// * DefaultActionChainRunCount
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultActionChainRunCount == nil, "General", "DefaultActionChainRunCount", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultActionChainRunCount, hasCorrectType := ConfigsJSON.General.DefaultActionChainRunCount.(float64)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultActionChainRunCount", "int", reflect.TypeOf(ConfigsJSON.General.DefaultActionChainRunCount)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if err := errorOnIncorrectValue(valueOfDefaultActionChainRunCount <= 0, "General", "DefaultActionChainRunCount", "Must be > 0", strconv.FormatFloat(valueOfDefaultActionChainRunCount, 'f', 6, 64)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.DefaultActionChainRunCount = uint32(valueOfDefaultActionChainRunCount)

	// ** Triggers Section
	// * SimultaneousCalls
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.SimultaneousCalls == nil, "Triggers", "SimultaneousCalls", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false

	} else {

		mapOfTriggerSimultaneousCallsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.SimultaneousCalls.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "SimultaneousCalls", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.SimultaneousCalls)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerSimultaneousCallsConfigJSON, "SimultaneousCalls")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = valueEnabled

		if valueEnabled {

			// TODO This is optional. So if this is not found use the default!
			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerSimultaneousCallsConfigJSON, "SimultaneousCalls", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerSimultaneousCallsConfigJSON, "Simultaneous Calls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SimultaneousCalls.HitThreshold = valueHitThreshold

			// TODO This is optional. So if this is not found use the default!
			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerSimultaneousCallsConfigJSON, "Simultaneous Calls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SimultaneousCalls.MinimumNumberLength = valueMinimumNumberLength

		}

	}

	// * DangerousDestinations
	fmt.Println(ConfigsJSON.Triggers.DangerousDestinations)
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.DangerousDestinations == nil, "Triggers", "DangerousDestinations", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.DangerousDestinations.Enabled = false

	} else {

		mapOfTriggerDangerousDestinationsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.DangerousDestinations.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "DangerousDestinations", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.SimultaneousCalls)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.DangerousDestinations.Enabled = valueEnabled

		if valueEnabled {

			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.DangerousDestinations.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.DangerousDestinations.HitThreshold = valueHitThreshold

			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.DangerousDestinations.MinimumNumberLength = valueMinimumNumberLength

			valuePrefixList, err := checkJSONSanityOfPrefixListConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.DangerousDestinations.PrefixList = valuePrefixList

		}

	}

	// * ExpectedDestinations
	fmt.Println(ConfigsJSON.Triggers.ExpectedDestinations)
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.ExpectedDestinations == nil, "Triggers", "ExpectedDestinations", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.ExpectedDestinations.Enabled = false

	} else {

		mapOfTriggerExpectedDestinationsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.ExpectedDestinations.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "ExpectedDestinations", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.ExpectedDestinations)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.ExpectedDestinations.Enabled = valueEnabled

		if valueEnabled {

			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.ExpectedDestinations.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.ExpectedDestinations.HitThreshold = valueHitThreshold

			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.ExpectedDestinations.MinimumNumberLength = valueMinimumNumberLength

			valuePrefixList, err := checkJSONSanityOfPrefixListConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.ExpectedDestinations.PrefixList = valuePrefixList

		}

	}

	// * SmallDurationCalls
	fmt.Println(ConfigsJSON.Triggers.SmallDurationCalls)
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.SmallDurationCalls == nil, "Triggers", "SmallDurationCalls", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.SmallDurationCalls.Enabled = false

	} else {

		mapOfTriggerSmallDurationCallsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.SmallDurationCalls.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "SmallDurationCalls", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.SmallDurationCalls)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.SmallDurationCalls.Enabled = valueEnabled

		if valueEnabled {

			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SmallDurationCalls.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SmallDurationCalls.HitThreshold = valueHitThreshold

			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SmallDurationCalls.MinimumNumberLength = valueMinimumNumberLength

			if err := errorOnValueNotFound(!utils.StringKeyInMap("DurationThreshold", mapOfTriggerSmallDurationCallsConfigJSON), "Triggers", "SmallDurationCalls", "map[string]interface{}", "\"nothing\""); err != nil {
				// TODO If this trigger is enabled, this value should be mandatory!
				utils.DebugLogAndGetError(err.Error(), true)
			}

			valueDurationString, hasCorrectType := mapOfTriggerSmallDurationCallsConfigJSON["DurationThreshold"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "SmallDurationCalls", "DurationThreshold", "time.Duration", reflect.TypeOf(mapOfTriggerSmallDurationCallsConfigJSON["DurationThreshold"])); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}

			valueDuration, errDuration := time.ParseDuration(valueDurationString)
			if err := errorOnIncorrectValue(errDuration != nil, "Small Duration Calls", "duration_threshold", "Must be a valid duration", ""); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}

			fraudionConfig.Triggers.SmallDurationCalls.DurationThreshold = valueDuration

		}

	}

	/*

		// Actions Section
		actionHTTPJSONConfig, actionIsPresent := JSONConfigs.Actions.HTTP.(map[string]interface{})
		fmt.Println("**", actionHTTPJSONConfig)
		fmt.Println(reflect.TypeOf(actionHTTPJSONConfig))
		fmt.Println(actionIsPresent)
		fmt.Println()

		valueDefaultMethod, hasCorrectType := actionHTTPJSONConfig["default_method"].(string)
		fmt.Println(valueDefaultMethod)
		fmt.Println(hasCorrectType)
		valueDefaultURL, hasCorrectType := actionHTTPJSONConfig["default_url"].(string)
		fmt.Println(valueDefaultURL)
		fmt.Println(hasCorrectType)

		// Doesn't work like this!
		valueDefaultParameters2, hasCorrectType := actionHTTPJSONConfig["default_parameters"].(map[string]string)
		//fmt.Println("**", reflect.ValueOf(actionHTTPJSONConfig["default_parameters"]).Interface())
		fmt.Println("**", valueDefaultParameters2)
		fmt.Println(reflect.TypeOf(valueDefaultParameters2))
		fmt.Println(hasCorrectType)

		// Works like this!
		valueDefaultParameters, hasCorrectType := actionHTTPJSONConfig["default_parameters"].(map[string]interface{})
		fmt.Println("**", reflect.ValueOf(actionHTTPJSONConfig["default_parameters"]).Interface())
		fmt.Println(valueDefaultParameters)
		fmt.Println(reflect.TypeOf(valueDefaultParameters))
		fmt.Println(hasCorrectType)
		for k := range valueDefaultParameters {
			fmt.Println(k)
			valueDefaultParameters, hasCorrectType := valueDefaultParameters[k].(string)
			if !hasCorrectType {
				fmt.Println("WHAAAAAAAAAAAAAAAAAAAAAAAT??")
			}
			fmt.Println(valueDefaultParameters)
			fmt.Println(reflect.TypeOf(valueDefaultParameters))
		}

	*/

	return nil
}

func errorOnIncorrectType(condition bool, sectionName string, valueNameOrPath string, typeExpected string, typeFound interface{}) error {
	if !condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" has an incorrect type (expected %s, found %s)", valueNameOrPath, sectionName, typeExpected, typeFound)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func errorOnIncorrectValue(condition bool, sectionName string, valueNameOrPath string, valueExpected string, valueFound string) error {
	if condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" is incorrect (expected %s, found %s)", valueNameOrPath, sectionName, valueExpected, valueFound)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func errorOnValueNotFound(condition bool, sectionName string, valueNameOrPath string, valueExpected string, valueFound string) error {
	if condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" not found", valueNameOrPath, sectionName)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func checkJSONSanityOfEnabledConfigs(simCallsTriggerJSONConfig map[string]interface{}, triggerName string) (bool, error) {

	if utils.StringKeyInMap("Enabled", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["enabled"].(bool)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "enabled", "boolean", reflect.TypeOf(simCallsTriggerJSONConfig["enabled"])); err != nil {
			return false, err
		}

		return value, nil

	}

	// TODO:LOG Log this to Syslog
	fmt.Printf("WARNING: \"Enabled\" value for \"%s\" Trigger not found, assuming true\n", triggerName)

	return true, nil

}

func checkJSONSanityOfCheckPeriodConfigs(fraudionConfig *FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string, minimumTriggerCheckPeriodDuration time.Duration) (time.Duration, error) {

	if utils.StringKeyInMap("CheckPeriod", simCallsTriggerJSONConfig) {

		valueString, hasCorrectType := simCallsTriggerJSONConfig["CheckPeriod"].(string)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "CheckPeriod", "string", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"])); err != nil {
			return 0, err
		}

		valueDuration, errDuration := time.ParseDuration(valueString)
		if err := errorOnIncorrectValue(errDuration != nil, triggerName, "CheckPeriod", "Must be a valid duration", valueString); err != nil {
			return 0, err
		}

		errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", constMinimumTriggerCheckPeriod)
		if err := errorOnIncorrectValue(valueDuration < minimumTriggerCheckPeriodDuration, triggerName, "check_period", errorMessageSuffix, valueString); err != nil {
			return 0, err
		}

		return valueDuration, nil

	}

	// TODO:LOG Log this to Syslog
	fmt.Printf("WARNING: \"CheckPeriod\" value for \"%s\" Trigger not found, assuming %s\n", triggerName, fraudionConfig.General.DefaultTriggerCheckPeriod)
	return fraudionConfig.General.DefaultTriggerCheckPeriod, nil

}

func checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig *FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) (uint32, error) {

	if utils.StringKeyInMap("MinimumNumberLength", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["MinimumNumberLength"].(float64)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "MinimumNumberLength", "float64", reflect.TypeOf(simCallsTriggerJSONConfig["minimum_number_length"])); err != nil {
			return 0, err
		}

		if err := errorOnIncorrectValue(value <= 0, triggerName, "MinimumNumberLength", "Must be > 0", strconv.FormatFloat(value, 'f', 6, 64)); err != nil {
			return 0, err
		}

		return uint32(value), nil

	}

	// This value is mandatory!

	// TODO:LOG Log this to Syslog
	fmt.Printf("WARNING: \"minimum_number_length\" value for \"%s\" Trigger not found, assuming %d\n", triggerName, fraudionConfig.General.DefaultMinimumDestinationNumberLength)
	return fraudionConfig.General.DefaultMinimumDestinationNumberLength, nil

}

func checkJSONSanityOfHitThresholdConfigs(fraudionConfig *FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) (uint32, error) {

	if utils.StringKeyInMap("hit_threshold", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["hit_threshold"].(float64)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "hit_threshold", "float64", reflect.TypeOf(simCallsTriggerJSONConfig["max_call_threshold"])); err != nil {
			return 0, err
		}

		if err := errorOnIncorrectValue(value <= 0, triggerName, "HitThreshold", "integer > 0", strconv.FormatFloat(value, 'f', 6, 64)); err != nil {
			return 0, err
		}

		return uint32(value), nil

	}

	// This value is mandatory!
	// TODO:LOG Log this to Syslog
	fmt.Printf("ERROR: \"hit_threshold\" value for \"%s\" Trigger not found", triggerName)
	return 0, fmt.Errorf("\"hit_threshold\" value for \"%s\" trigger not found", triggerName)

}

func checkJSONSanityOfPrefixListConfigs(fraudionConfig *FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) ([]string, error) {

	if utils.StringKeyInMap("prefix_list", simCallsTriggerJSONConfig) {

		sliceOfInterface := simCallsTriggerJSONConfig["prefix_list"]

		// WARNING: Just checks if it's a slice basically (of anything really because interface{}) (verification of the type of the items of the slice is done below)
		valueInterfaceSlice, hasCorrectType := sliceOfInterface.([]interface{})
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "prefix_list", "[]string", reflect.TypeOf(sliceOfInterface)); err != nil {
			return *new([]string), err
		}

		sliceReflectValue := reflect.ValueOf(valueInterfaceSlice)
		var finalSlice []string
		for i := 0; i < sliceReflectValue.Len(); i++ {

			// TODO This is not needed because the previous check catches this, because ints implement no interface? Don't know but it catches when the slice has ints...
			valueOfSliceItem, hasCorrectType := sliceReflectValue.Index(i).Interface().(string)
			if err := errorOnIncorrectType(hasCorrectType, triggerName, "prefix_list", "[]", reflect.TypeOf(sliceOfInterface)); err != nil {
				return *new([]string), err
			}

			isNumericString, errCheckingValue := regexp.MatchString("^[0-9]+$", valueOfSliceItem)
			if errCheckingValue != nil {
				fmt.Printf("ERROR: Internal: There seems to be an issue with checking if the \"prefix_list\" string values are numeric\n")
				return *new([]string), fmt.Errorf("internal: there seems to be an issue with checking if the \"prefix_list\" string values are numeric")
			}
			if err := errorOnIncorrectValue(!isNumericString, triggerName, "prefix_list", "Values must be Numeric Strings", valueOfSliceItem); err != nil {
				return *new([]string), err
			}

			finalSlice = append(finalSlice, valueOfSliceItem)
		}

		return finalSlice, nil

	}

	// TODO:LOG Log this to Syslog
	fmt.Printf("ERROR: \"prefix_list\" value for \"%s\" Trigger not found\n", triggerName)
	return *new([]string), fmt.Errorf("\"prefix_list\" value for \"%s\" not found", triggerName)

}
