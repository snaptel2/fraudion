package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fraudion/utils"
)

const (
	constMinimumTriggerCheckPeriod     = "5m"
	constMaximumActionChainSleepPeriod = "5m"
)

var (
	constSupportedSoftware = []string{"*ast_el_2.3", "*ast_alone_1.8"}
	constSupportCDRSources = []string{"*db_mysql"}
)

// CheckJSONSanityAndLoadConfigs ...
func (fraudionConfig *FraudionConfig) CheckJSONSanityAndLoadConfigs(JSONConfigs *FraudionConfigJSON, validateOnly bool) error {

	// ** General Section
	// Monitored Software
	valueOfMonitoredSoftware, hasCorrectType := JSONConfigs.General.MonitoredSoftware.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "MonitoredSoftware", "string", reflect.TypeOf(valueOfMonitoredSoftware)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	found := utils.StringInStringsSlice(valueOfMonitoredSoftware, constSupportedSoftware)
	if err := errorOnValueNotFound(!found, "General", "MonitoredSoftware", fmt.Sprintf("%s", constSupportedSoftware), valueOfMonitoredSoftware, true); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.MonitoredSoftware = valueOfMonitoredSoftware

	// CDRs Source
	valueOfCDRsSource, hasCorrectType := JSONConfigs.General.CDRsSource.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "CDRsSource", "string", reflect.TypeOf(valueOfCDRsSource)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	found = utils.StringInStringsSlice(valueOfMonitoredSoftware, constSupportCDRSources)
	if err := errorOnValueNotFound(!found, "General", "CDRsSource", fmt.Sprintf("%s", constSupportCDRSources), valueOfCDRsSource, true); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.CDRsSource = valueOfCDRsSource

	// DefaultTriggerCheckPeriod
	valueOfDefaultTriggerCheckPeriod, hasCorrectType := JSONConfigs.General.DefaultTriggerCheckPeriod.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultTriggerCheckPeriod", "string", reflect.TypeOf(valueOfDefaultTriggerCheckPeriod)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfDefaultTriggerCheckPeriod, errDuration := time.ParseDuration(valueOfDefaultTriggerCheckPeriod)
	if err := errorOnIncorrectValue(errDuration != nil, "General", "DefaultTriggerCheckPeriod", "parseable Duration (e.g. 5\"s\", \"m\" or \"h\")", valueOfDefaultTriggerCheckPeriod); err != nil {
		return err
	}

	durationValueOfConstMinimumTriggerCheckPeriod, errMinimumCheckPeriodDuration := time.ParseDuration(constMinimumTriggerCheckPeriod)
	errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", constMinimumTriggerCheckPeriod)
	if err := errorOnIncorrectValue(durationValueOfDefaultTriggerCheckPeriod < durationValueOfConstMinimumTriggerCheckPeriod, "General", "DefaultTriggerCheckPeriod", errorMessageSuffix, ""); err != nil {
		return err
	}

	if errMinimumCheckPeriodDuration != nil {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: Internal: There seems to be an issue with the definition of the constant constMinimumTriggerCheckPeriod")
		return fmt.Errorf("internal: there seems to be an issue with the definition of the constant constMinimumTriggerCheckPeriod")
	}

	fraudionConfig.General.DefaultTriggerCheckPeriod = durationValueOfDefaultTriggerCheckPeriod

	/*//  DefaultActionChainSleePeriod
	valueOfDefaultActionChainSleePeriod, hasCorrectType := JSONConfigs.General.DefaultActionChainSleePeriod.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "valueOfDefaultActionChainSleePeriod", "string", reflect.TypeOf(valueOfDefaultActionChainSleePeriod)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	defaultActionChainSleepPeriodDuration, err := time.ParseDuration(valueOfDefaultActionChainSleePeriod)
	if err != nil {
		// TODO:LOG Log this to Syslog
		fmt.Println("ERROR: It seems that the value in \"default_action_chain_sleep_period\" is not a valid duration. :(")
		return fmt.Errorf("it seems that the value in \"default_action_chain_sleep_period\" is not a valid duration")
	}

	MaximumActionChainSleepPeriodDuration, errMaximumChainSleepPeriod := time.ParseDuration(constMaximumActionChainSleepPeriod)
	if defaultActionChainSleepPeriodDuration > MaximumActionChainSleepPeriodDuration {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: \"default_trigger_check_period\" is too big. Maximum value is \"%s\". :(", constMaximumActionChainSleepPeriod)
		return fmt.Errorf("\"default_trigger_check_period\" is too big. maximum value is \"%s\"", constMaximumActionChainSleepPeriod)
	}

	if errMaximumChainSleepPeriod != nil {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: Internal: There seems to be an issue with the definition of the constant constMaximumActionChainSleepPeriod")
		return fmt.Errorf("internal: there seems to be an issue with the definition of the constant constMaximumActionChainSleepPeriod")
	}

	fraudionConfig.General.DefaultActionChainSleepPeriod = defaultActionChainSleepPeriodDuration

	//  Default_action_chain_run_count
	// TODO Add some kind of limits here?
	fraudionConfig.General.DefaultActionChainRunCount = JSONConfigs.General.DefaultActionChainRunCount

	//  Default_minimum_destination_number_length
	// TODO Add some kind of limits here?
	fraudionConfig.General.DefaultMinimumDestinationNumberLength = JSONConfigs.General.DefaultMinimumDestinationNumberLength

	// **** Triggers Section
	//  Simultaneous calls
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent := JSONConfigs.Triggers.SimultaneousCalls.(map[string]interface{})
	if triggerIsPresent {

		// ** Search "enabled" Key [Optional], defaults to "enabled"
		valueEnabled, err := checkJSONSanityOfEnabledConfigs(simCallsTriggerJSONConfig, "Simultaneous Calls")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = valueEnabled

		// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		valueCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Simultaneous Calls", minimumTriggerCheckPeriodDuration)
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = valueCheckPeriod

		// ** Search "hit_threshold" Key [Mandatory]
		valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Simultaneous Calls")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SimultaneousCalls.HitThreshold = valueHitThreshold

		// ** Search "minimum_number_length" Key [Optional], defaults to "fraudionConfig.General.DefaultMinimumDestinationNumberLength"
		valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Simultaneous Calls")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SimultaneousCalls.MinimumNumberLength = valueMinimumNumberLength

	} else {
		// TODO:LOG Log this to Syslog
		fmt.Println("WARNING: \"Simultaneous Calls\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
	}

	fmt.Println("******", fraudionConfig.Triggers)

	//  Dangerous Destinations
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent = JSONConfigs.Triggers.Dangerous_destinations.(map[string]interface{})
	if triggerIsPresent {

		// ** Search "enabled" Key [Optional], defaults to "enabled"
		valueEnabled, err := checkJSONSanityOfEnabledConfigs(simCallsTriggerJSONConfig, "Dangerous Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.DangerousDestinations.Enabled = valueEnabled

		// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		valueCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Dangerous Destinations", minimumTriggerCheckPeriodDuration)
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.DangerousDestinations.CheckPeriod = valueCheckPeriod

		// ** Search "hit_threshold" Key [Mandatory]
		valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Dangerous Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.DangerousDestinations.HitThreshold = valueHitThreshold

		// ** Search "minimum_number_length" Key [Optional], defaults to "fraudionConfig.General.DefaultMinimumDestinationNumberLength"
		valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Dangerous Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.DangerousDestinations.MinimumNumberLength = valueMinimumNumberLength

		// ** Search "prefix_list" Key [Mandatory + not empty]
		valuePrefixList, err := checkJSONSanityOfPrefixListConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Dangerous Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.DangerousDestinations.PrefixList = valuePrefixList

	} else {
		// TODO:LOG Log this to Syslog
		fmt.Println("WARNING: \"Dangerous Destinations\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
	}

	fmt.Println("******", fraudionConfig.Triggers)

	//  Expected Destinations
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent = JSONConfigs.Triggers.Expected_destinations.(map[string]interface{})
	if triggerIsPresent {

		// ** Search "enabled" Key [Optional], defaults to "enabled"
		valueEnabled, err := checkJSONSanityOfEnabledConfigs(simCallsTriggerJSONConfig, "Expected Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.ExpectedDestinations.Enabled = valueEnabled

		// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		valueCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations", minimumTriggerCheckPeriodDuration)
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.ExpectedDestinations.CheckPeriod = valueCheckPeriod

		// ** Search "hit_threshold" Key [Mandatory]
		valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.ExpectedDestinations.HitThreshold = valueHitThreshold

		// ** Search "minimum_number_length" Key [Optional], defaults to "fraudionConfig.General.DefaultMinimumDestinationNumberLength"
		valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.ExpectedDestinations.MinimumNumberLength = valueMinimumNumberLength

		// ** Search "prefix_list" Key [Mandatory + not empty]
		valuePrefixList, err := checkJSONSanityOfPrefixListConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.ExpectedDestinations.PrefixList = valuePrefixList

	} else {
		// TODO:LOG Log this to Syslog
		fmt.Println("WARNING: \"Expected Destinations\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
	}

	fmt.Println("******", fraudionConfig.Triggers)

	//  Small Duration calls
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent = JSONConfigs.Triggers.Small_duration_calls.(map[string]interface{})
	if triggerIsPresent {

		// ** Search "enabled" Key [Optional], defaults to "enabled"
		valueEnabled, err := checkJSONSanityOfEnabledConfigs(simCallsTriggerJSONConfig, "Small Duration Calls")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SmallDurationCalls.Enabled = valueEnabled

		// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		valueCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Small Duration Calls", minimumTriggerCheckPeriodDuration)
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SmallDurationCalls.CheckPeriod = valueCheckPeriod

		// ** Search "hit_threshold" Key [Mandatory]
		valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Small Duration Calls")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SmallDurationCalls.HitThreshold = valueHitThreshold

		// ** Search "minimum_number_length" Key [Optional], defaults to "fraudionConfig.General.DefaultMinimumDestinationNumberLength"
		valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Small Duration Calls")
		if err != nil {
			return err
		}
		fraudionConfig.Triggers.SmallDurationCalls.MinimumNumberLength = valueMinimumNumberLength

		// ** Search "duration_threshold" Key [Mandatory]
		if utils.StringKeyInMap("duration_threshold", simCallsTriggerJSONConfig) {

			valueDurationString, hasCorrectType := simCallsTriggerJSONConfig["duration_threshold"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Small Duration Calls", "duration_threshold", "time.Duration", reflect.TypeOf(simCallsTriggerJSONConfig["duration_threshold"])); err != nil {
				return err
			}

			valueDuration, errDuration := time.ParseDuration(valueDurationString)
			if err := errorOnIncorrectValue(errDuration != nil, "Small Duration Calls", "duration_threshold", "Must be a valid duration"); err != nil {
				return err
			}

			fraudionConfig.Triggers.SmallDurationCalls.DurationThreshold = valueDuration

		} else {

			// TODO:LOG Log this to Syslog
			fmt.Printf("ERROR: \"duration_threshold\" value for \"Small Duration Calls\" Trigger not found :(\n")
			//fraudionConfig.Triggers.SmallDurationCalls.DurationThreshold = fraudionConfig.General.DefaultTriggerCheckPeriod
			return fmt.Errorf("ERROR: \"duration_threshold\" value for \"Small Duration Calls\" Trigger not found\n")

		}

	} else {
		// TODO:LOG Log this to Syslog
		fmt.Println("WARNING: \"Small Duration Calls\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SmallDurationCalls.Enabled = false
	}

	fmt.Println("******", fraudionConfig.Triggers)

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
	}*/
	/*valueDefaultMessage2, hasCorrectType := valueDefaultMessage.(map[string]string)
	fmt.Println("**", reflect.ValueOf(valueDefaultMessage))
	fmt.Println(valueDefaultMessage2)
	fmt.Println(reflect.TypeOf(valueDefaultMessage2))
	fmt.Println(hasCorrectType)*/
	//buh := map[string]string{
	//	"rsc": "3711",
	//	"r":   "2138"
	//"http_post_parameters_1_k": "http_post_parameters_1_v",
	//"http_post_parameters_2_k": "http_post_parameters_2_v"
	//}
	/*if triggerIsPresent {

	}*/
	/*
			valueDefaultMessage, hasCorrectType := simCallsActionsJSONConfig["default_message"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Email", "Default Message", "string", reflect.TypeOf(simCallsTriggerJSONConfig["email"])); err != nil {
				return err
			}
			fraudionConfig.Actions.Email.DefaultMessage = valueDefaultMessage

		} else {
			// TODO:LOG Log this to Syslog
			fmt.Println("WARNING: \"Expected Destinations\" Trigger config not present, disabling!")
			fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
		}

		simCallsActionsJSONConfig, triggerIsPresent = JSONConfigs.Actions.HTTP.(map[string]interface{})
		if triggerIsPresent {

			valueDefaultMethod, hasCorrectType := simCallsActionsJSONConfig["default_method"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "HTTP", "Default Method", "string", reflect.TypeOf(simCallsTriggerJSONConfig["email"])); err != nil {
				return err
			}
			fraudionConfig.Actions.HTTP.DefaultMethod = valueDefaultMethod

			valueDefaultURL, hasCorrectType := simCallsActionsJSONConfig["default_url"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "HTTP", "Default URL", "string", reflect.TypeOf(simCallsTriggerJSONConfig["email"])); err != nil {
				return err
			}
			fraudionConfig.Actions.HTTP.DefaultURL = valueDefaultURL

			valueInterfaceParameters, hasCorrectType := simCallsActionsJSONConfig["default_parameters"].(map[string]interface{})
			//fmt.Println(simCallsActionsJSONConfig["default_parameters"].(map[string]interface{}))
			if err := errorOnIncorrectType(hasCorrectType, "HTTP", "Default Parameters", "map[string]interface{}", reflect.TypeOf(simCallsTriggerJSONConfig["email"])); err != nil {
				return err
			}

			mapReflectValue := reflect.ValueOf(valueInterfaceParameters)
			var finalMap map[string]string
			for i := 0; i < mapReflectValue.Len(); i++ {

				// TODO This is not needed because the previous check catches this, because ints implement no interface? Don't know but it catches when the slice has ints...
				valueOfMapItem, hasCorrectType := mapReflectValue.In.Interface().(map[string]string)
				if err := errorOnIncorrectType(hasCorrectType, "HTTP", "default_parameters", "[]", reflect.TypeOf(mapReflectValue.Interface())); err != nil {
					return err
				}

				fmt.Println(valueOfMapItem)

				/*isNumericString, errCheckingValue := regexp.MatchString("^[0-9]+$", valueOfSliceItem)
				if errCheckingValue != nil {
					fmt.Printf("ERROR: Internal: There seems to be an issue with checking if the \"prefix_list\" string values are numeric\n")
					return *new([]string), fmt.Errorf("internal: there seems to be an issue with checking if the \"prefix_list\" string values are numeric")
				}
				if err := errorOnIncorrectValue(!isNumericString, triggerName, "prefix_list", "Values must be Numeric Strings"); err != nil {
					return *new([]string), err
				}*/

	//finalMap = append(finalMap, valueOfMapItem)
	//}

	//fraudionConfig.Actions.HTTP.DefaultParameters = finalMap

	/*
		valueDefaultMessage, hasCorrectType := simCallsActionsJSONConfig["default_message"].(string)
		if err := errorOnIncorrectType(hasCorrectType, "Email", "Default Message", "string", reflect.TypeOf(simCallsTriggerJSONConfig["email"])); err != nil {
			return err
		}
		fraudionConfig.Actions.HTTP.DefaultURL = valueDefaultMessage

		valueDefaultMessage, hasCorrectType := simCallsActionsJSONConfig["default_message"].(string)
		if err := errorOnIncorrectType(hasCorrectType, "Email", "Default Message", "string", reflect.TypeOf(simCallsTriggerJSONConfig["email"])); err != nil {
			return err
		}
		fraudionConfig.Actions.HTTP.DefaultParameters = valueDefaultMessage

		/*	// ** Search "enabled" Key [Optional], defaults to "enabled"
			valueEnabled, err := checkJSONSanityOfEnabledConfigs(simCallsTriggerJSONConfig, "Expected Destinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.ExpectedDestinations.Enabled = valueEnabled

			// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
			valueCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations", minimumTriggerCheckPeriodDuration)
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.ExpectedDestinations.CheckPeriod = valueCheckPeriod

			// ** Search "hit_threshold" Key [Mandatory]
			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.ExpectedDestinations.HitThreshold = valueHitThreshold

			// ** Search "minimum_number_length" Key [Optional], defaults to "fraudionConfig.General.DefaultMinimumDestinationNumberLength"
			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.ExpectedDestinations.MinimumNumberLength = valueMinimumNumberLength

			// ** Search "prefix_list" Key [Mandatory + not empty]
			valuePrefixList, err := checkJSONSanityOfPrefixListConfigs(fraudionConfig, simCallsTriggerJSONConfig, "Expected Destinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.ExpectedDestinations.PrefixList = valuePrefixList
	*/
	/*} else {
		// TODO:LOG Log this to Syslog
		fmt.Println("WARNING: \"Expected Destinations\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
	}*/

	//simCallsTriggerJSONConfig, triggerIsPresent = JSONConfigs.Actions.HTTP.(map[string]interface{})

	//simCallsTriggerJSONConfig, triggerIsPresent = JSONConfigs.Actions.Call.(map[string]interface{})

	//simCallsTriggerJSONConfig, triggerIsPresent = JSONConfigs.Actions.Local_commands.(map[string]string)

	// ActionChains Section

	// Contacts Section

	return nil
}

func errorOnIncorrectType(condition bool, sectionName string, valueNameOrPath string, typeExpected string, typeFound interface{}) error {
	if !condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" has an incorrect type (expected %s, found %s)", valueNameOrPath, sectionName, typeExpected, typeFound)
		return fmt.Errorf("%s", strings.ToLower(errorMessage))
	}
	return nil
}

func errorOnIncorrectValue(condition bool, sectionName string, valueNameOrPath string, valueExpected string, valueFound string) error {
	if condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" is incorrect type (expected %s, found %s)", valueNameOrPath, sectionName, valueExpected, valueFound)
		return fmt.Errorf("%s", strings.ToLower(errorMessage))
	}
	return nil
}

func errorOnValueNotFound(condition bool, sectionName string, valueNameOrPath string, valueExpected string, valueFound string, mandatory bool) error {
	if condition {
		var mandatoryString string
		if mandatory {
			mandatoryString = "yes"
		} else {
			mandatoryString = "no"
		}
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" not found (mandatory? %s)", valueNameOrPath, sectionName, mandatoryString)
		return fmt.Errorf("%s", strings.ToLower(errorMessage))
	}
	return nil
}

func checkJSONSanityOfEnabledConfigs(simCallsTriggerJSONConfig map[string]interface{}, triggerName string) (bool, error) {

	if utils.StringKeyInMap("enabled", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["enabled"].(bool)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "enabled", "boolean", reflect.TypeOf(simCallsTriggerJSONConfig["enabled"])); err != nil {
			return false, err
		}

		return value, nil

	}

	// TODO:LOG Log this to Syslog
	fmt.Printf("WARNING: \"enabled\" value for \"%s\" Trigger not found, assuming true\n", triggerName)

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
