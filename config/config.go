package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/fraudion/utils"
)

const (
	MINIMUM_TRIGGER_CHECK_PERIOD      = "5m"
	MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD = "5m"
)

var (
	SUPPORTED_SOFTWARE     = []string{"*ast_el_2.3", "*ast_alone_1.8"}
	SUPPORTED_CDRS_SOURCES = []string{"*db_mysql"}
)

// FraudionConfig ...
type FraudionConfig struct {
	General      General
	Triggers     Triggers
	Actions      Actions
	ActionChains ActionChains
	Contacts     Contacts
}

// General ...
type General struct {
	MonitoredSoftware                     string
	CdrsSource                            string
	DefaultTriggerCheckPeriod             time.Duration
	DefaultActionChainSleepPeriod         time.Duration
	DefaultActionChainRunCount            uint32
	DefaultMinimumDestinationNumberLength uint32
	// TODO Default Hit Threshold
}

// Triggers ...
type Triggers struct {
	SimultaneousCalls     triggerSimultaneousCalls
	DangerousDestinations triggerDangerousDestinations
	ExpectedDestinations  triggerExpectedDestinations
	SmallDurationCalls    triggerSmallCallDurations
}

type triggerSimultaneousCalls struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
}

type triggerDangerousDestinations struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
	PrefixList          []string
}

type triggerExpectedDestinations struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
	PrefixList          []string
}

type triggerSmallCallDurations struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
	DurationThreshold   time.Duration
}

// Actions ...
type Actions struct {
	Email         actionEmail
	Call          actionCall
	HTTP          actionHTTP
	LocalCommands map[string]string
}

type actionEmail struct {
	DefaultMessage string
}

type actionCall struct {
	DefaultMessage string
}

type actionHTTP struct {
	DefaultURL        string
	DefaultMethod     string
	DefaultParameters map[string]string
}

type actionLocalCommands struct {
	List map[string]string
}

// ActionChains ...
type ActionChains struct {
	List map[string][]actionChainAction
}

type actionChainAction struct {
	Action   string
	Contacts []string
	Command  string
}

// Contacts ...
type Contacts struct {
	List map[string]contact
}

type contact struct {
	ForActions     []string
	PhoneNumber    string
	Email          string
	Message        string
	HTTPURL        string
	HTTPMethod     string
	HTTPParameters map[string]string
}

// CheckJSONSanityAndLoadConfigs ...
func (fraudionConfig *FraudionConfig) CheckJSONSanityAndLoadConfigs(JSONConfigs *FraudionJSONConfig) error {

	// General Section
	//  Monitored Software
	found := utils.StringInStringsSlice(JSONConfigs.General.Monitored_software, SUPPORTED_SOFTWARE)
	if found == false {
		// TODO:LOG Log this to Syslog
		fmt.Println("Could not find configured \"monitored_software\" in supported software list. :(")
		return fmt.Errorf("could not find configured \"monitored_software\" in supported software list")
	}

	fraudionConfig.General.MonitoredSoftware = JSONConfigs.General.Monitored_software

	//  CDRs Source
	found = utils.StringInStringsSlice(JSONConfigs.General.Cdrs_source, SUPPORTED_CDRS_SOURCES)
	if found == false {
		// TODO:LOG Log this to Syslog
		fmt.Println("Could not find configured \"cdrs_source\" in supported CDR sources list. :(")
		return fmt.Errorf("could not find configured \"cdrs_source\" in CDR sources list")
	}

	fraudionConfig.General.CdrsSource = JSONConfigs.General.Cdrs_source

	//  Default_trigger_check_period
	defaultTriggerCheckPeriodDuration, errDuration := time.ParseDuration(JSONConfigs.General.Default_trigger_check_period)
	if err := errorOnIncorrectValue(errDuration != nil, "General", "default_trigger_check_period", "Must be a valid duration"); err != nil {
		return err
	}

	minimumTriggerCheckPeriodDuration, errMinimumCheckPeriodDuration := time.ParseDuration(MINIMUM_TRIGGER_CHECK_PERIOD)
	errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
	if err := errorOnIncorrectValue(defaultTriggerCheckPeriodDuration < minimumTriggerCheckPeriodDuration, "General", "default_trigger_check_period", errorMessageSuffix); err != nil {
		return err
	}

	if errMinimumCheckPeriodDuration != nil {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: Internal: There seems to be an issue with the definition of the constant MINIMUM_TRIGGER_CHECK_PERIOD")
		return fmt.Errorf("internal: there seems to be an issue with the definition of the constant MINIMUM_TRIGGER_CHECK_PERIOD")
	}

	fraudionConfig.General.DefaultTriggerCheckPeriod = defaultTriggerCheckPeriodDuration

	//  Default_action_chain_sleep_period
	defaultActionChainSleepPeriodDuration, err := time.ParseDuration(JSONConfigs.General.Default_action_chain_sleep_period)
	if err != nil {
		// TODO:LOG Log this to Syslog
		fmt.Println("ERROR: It seems that the value in \"default_action_chain_sleep_period\" is not a valid duration. :(")
		return fmt.Errorf("it seems that the value in \"default_action_chain_sleep_period\" is not a valid duration")
	}

	MaximumActionChainSleepPeriodDuration, errMaximumChainSleepPeriod := time.ParseDuration(MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
	if defaultActionChainSleepPeriodDuration > MaximumActionChainSleepPeriodDuration {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: \"default_trigger_check_period\" is too big. Maximum value is \"%s\". :(", MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
		return fmt.Errorf("\"default_trigger_check_period\" is too big. maximum value is \"%s\"", MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
	}

	if errMaximumChainSleepPeriod != nil {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: Internal: There seems to be an issue with the definition of the constant MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD")
		return fmt.Errorf("internal: there seems to be an issue with the definition of the constant MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD")
	}

	fraudionConfig.General.DefaultActionChainSleepPeriod = defaultActionChainSleepPeriodDuration

	//  Default_action_chain_run_count
	// TODO Add some kind of limits here?
	fraudionConfig.General.DefaultActionChainRunCount = JSONConfigs.General.Default_action_chain_run_count

	//  Default_minimum_destination_number_length
	// TODO Add some kind of limits here?
	fraudionConfig.General.DefaultMinimumDestinationNumberLength = JSONConfigs.General.Default_minimum_destination_number_length

	// **** Triggers Section
	//  Simultaneous calls
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent := JSONConfigs.Triggers.Simultaneous_calls.(map[string]interface{})
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
	simCallsActionsJSONConfig, triggerIsPresent := JSONConfigs.Actions.HTTP.(map[string]interface{})
	fmt.Println("**", simCallsActionsJSONConfig)
	fmt.Println(reflect.TypeOf(simCallsActionsJSONConfig))
	fmt.Println(triggerIsPresent)
	fmt.Println()
	valueDefaultMessage, hasCorrectType := simCallsActionsJSONConfig["default_parameters"].(map[string]interface{})
	fmt.Println("**", reflect.ValueOf(simCallsActionsJSONConfig["default_parameters"]).Interface())
	fmt.Println(valueDefaultMessage)
	fmt.Println(reflect.TypeOf(valueDefaultMessage))
	fmt.Println(hasCorrectType)
	for k, _ := range valueDefaultMessage {
		fmt.Println(k)
		valueDefaultMessage, hasCorrectType := valueDefaultMessage[k].(string)
		if !hasCorrectType {
			fmt.Println("WHAAAAAAAAAAAAAAAAAAAAAAAT??")
		}
		fmt.Println(valueDefaultMessage)
		fmt.Println(reflect.TypeOf(valueDefaultMessage))
	}
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

func errorOnIncorrectType(condition bool, triggerName string, triggerValue string, expectedType string, foundType interface{}) error {
	if !condition {
		// TODO:LOG Log this to Syslog
		errorMessage := fmt.Sprintf("\"%s\" value for \"%s\" Trigger not of correct type, %s expected found %s", triggerValue, triggerName, expectedType, foundType)
		fmt.Printf("ERROR: %s :(\n", errorMessage)
		return fmt.Errorf("%s", strings.ToLower(errorMessage))
	}
	return nil
}

func errorOnIncorrectValue(condition bool, triggerName string, triggerValue string, whatItMustBeMessage string) error {

	if condition {
		// TODO:LOG Log this to Syslog
		errorMessage := fmt.Sprintf("\"%s\" value for \"%s\" in \"%s\" section is incorrect. ", triggerValue, triggerName, whatItMustBeMessage)
		fmt.Printf("ERROR: %s :(", errorMessage)
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

	if utils.StringKeyInMap("check_period", simCallsTriggerJSONConfig) {

		valueString, hasCorrectType := simCallsTriggerJSONConfig["check_period"].(string)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "check_period", "string", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"])); err != nil {
			return 0, err
		}

		valueDuration, errDuration := time.ParseDuration(valueString)
		if err := errorOnIncorrectValue(errDuration != nil, triggerName, "check_period", "Must be a valid duration"); err != nil {
			return 0, err
		}

		errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
		if err := errorOnIncorrectValue(valueDuration < minimumTriggerCheckPeriodDuration, triggerName, "check_period", errorMessageSuffix); err != nil {
			return 0, err
		}

		return valueDuration, nil

	}

	// TODO:LOG Log this to Syslog
	fmt.Printf("WARNING: \"check_period\" value for \"%s\" Trigger not found, assuming %s\n", triggerName, fraudionConfig.General.DefaultTriggerCheckPeriod)
	return fraudionConfig.General.DefaultTriggerCheckPeriod, nil

}

func checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig *FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) (uint32, error) {

	if utils.StringKeyInMap("minimum_number_length", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["minimum_number_length"].(float64)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "minimum_number_length", "float64", reflect.TypeOf(simCallsTriggerJSONConfig["minimum_number_length"])); err != nil {
			return 0, err
		}

		if err := errorOnIncorrectValue(value <= 0, triggerName, "minimum_number_length", "Must be > 0"); err != nil {
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

		if err := errorOnIncorrectValue(value <= 0, triggerName, "hit_threshold", "Must be > 0"); err != nil {
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
			if err := errorOnIncorrectValue(!isNumericString, triggerName, "prefix_list", "Values must be Numeric Strings"); err != nil {
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
