package config

import (
	"fmt"
	"reflect"
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
	MaxCallThreshold    uint32
	MinimumNumberLength uint32
}

type triggerDangerousDestinations struct {
	Enabled      bool
	CheckPeriod  time.Duration
	PrefixList   []string
	HitThreshold uint32
}

type triggerExpectedDestinations struct {
	Enabled      bool
	CheckPeriod  time.Duration
	PrefixList   []string
	HitThreshold uint32
}

type triggerSmallCallDurations struct {
	Enabled           bool
	CheckPeriod       time.Duration
	DurationThreshold time.Duration
	HitThreshold      uint32
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
	// TODO Do we need to check if the section is defined OR the JSON import fails if it is not?
	/*
		Monitored_software                        string
		Cdrs_source                               string
		Default_trigger_check_period              string
		Default_action_chain_sleep_period         string
		Default_action_chain_run_count            uint32
		Default_minimum_destination_number_length uint32
	*/
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
		fmt.Println("Could not find configured \"cdrs_source\" in supported software list. :(")
		return fmt.Errorf("could not find configured \"cdrs_source\" in supported software list")
	}

	fraudionConfig.General.CdrsSource = JSONConfigs.General.Cdrs_source

	//  Default_trigger_check_period
	defaultTriggerCheckPeriodDuration, err := time.ParseDuration(JSONConfigs.General.Default_trigger_check_period)
	if err != nil {
		// TODO:LOG Log this to Syslog
		fmt.Println("ERROR: It seems that the value in \"default_trigger_check_period\" is not a valid duration. :(")
		return fmt.Errorf("it seems that the value in \"default_trigger_check_period\" is not a valid duration")
	}

	MinimumTriggerCheckPeriodDuration, err := time.ParseDuration(MINIMUM_TRIGGER_CHECK_PERIOD)
	if defaultTriggerCheckPeriodDuration < MinimumTriggerCheckPeriodDuration {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: \"default_trigger_check_period\" is too small. Minimum value is \"%s\". :(", MINIMUM_TRIGGER_CHECK_PERIOD)
		return fmt.Errorf("\"default_trigger_check_period\" is too small. minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
	}
	if err != nil {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: Internal: There seems to be an issue with the definition of the constant MINIMUM_TRIGGER_CHECK_PERIOD.")
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

	MaximumActionChainSleepPeriodDuration, err := time.ParseDuration(MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
	if defaultActionChainSleepPeriodDuration > MaximumActionChainSleepPeriodDuration {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: \"default_trigger_check_period\" is too big. Maximum value is \"%s\". :(", MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
		return fmt.Errorf("\"default_trigger_check_period\" is too big. maximum value is \"%s\"", MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
	}
	if err != nil {
		// TODO:LOG Log this to Syslog
		fmt.Printf("ERROR: Internal: There seems to be an issue with the definition of the constant MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD.")
		return fmt.Errorf("internal: there seems to be an issue with the definition of the constant MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD")
	}

	fraudionConfig.General.DefaultActionChainSleepPeriod = defaultActionChainSleepPeriodDuration

	//  Default_action_chain_run_count
	// TODO Add some kind of limit/minimum here?
	fraudionConfig.General.DefaultActionChainRunCount = JSONConfigs.General.Default_action_chain_run_count

	//  Default_minimum_destination_number_length
	// TODO Add some kind of limit/minimum here?
	fraudionConfig.General.DefaultMinimumDestinationNumberLength = JSONConfigs.General.Default_minimum_destination_number_length

	// **** Triggers Section
	//  Simultaneous calls
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent := JSONConfigs.Triggers.Simultaneous_calls.(map[string]interface{})
	if triggerIsPresent {

		// ** Search "enabled" Key [Optional], defaults to "enabled"

		/*if utils.StringKeyInMap("enabled", simCallsTriggerJSONConfig) {

			value, hasCorrectType := simCallsTriggerJSONConfig["enabled"].(bool)
			if err := errorOnIncorrectType(hasCorrectType, "Simultaneous Calls", "enabled", "boolean", reflect.TypeOf(simCallsTriggerJSONConfig["enabled"])); err != nil {
				return err
			}

			fraudionConfig.Triggers.SimultaneousCalls.Enabled = value

		} else {
			// TODO:LOG Log this to Syslog
			fmt.Println("WARNING: \"Enabled\" value for \"Simultaneous Calls\" Trigger not found, assuming true")
			fraudionConfig.Triggers.SimultaneousCalls.Enabled = true
		}

		// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		if utils.StringKeyInMap("check_period", simCallsTriggerJSONConfig) {

			valueString, hasCorrectType := simCallsTriggerJSONConfig["check_period"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Simultaneous Calls", "check_period", "string", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"])); err != nil {
				return err
			}

			valueDuration, errDuration := time.ParseDuration(valueString)
			if err := errorOnIncorrectValue(errDuration != nil, "Simultaneous Calls", "check_period", "Must be a valid duration"); err != nil {
				return err
			}

			errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
			if err := errorOnIncorrectValue(valueDuration < MinimumTriggerCheckPeriodDuration, "Simultaneous Calls", "check_period", errorMessageSuffix); err != nil {
				return err
			}

			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = valueDuration

		} else {
			// TODO:LOG Log this to Syslog
			fmt.Printf("WARNING: \"check_period\" value for \"Simultaneous Calls\" Trigger not found, assuming %s\n", fraudionConfig.General.DefaultTriggerCheckPeriod)
			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = fraudionConfig.General.DefaultTriggerCheckPeriod
		}*/

		// ** Search "max_call_threshold" Key [Mandatory]
		if utils.StringKeyInMap("max_call_threshold", simCallsTriggerJSONConfig) {

			value, hasCorrectType := simCallsTriggerJSONConfig["max_call_threshold"].(float64)
			if err := errorOnIncorrectType(hasCorrectType, "Simultaneous Calls", "max_call_threshold", "float64", reflect.TypeOf(simCallsTriggerJSONConfig["max_call_threshold"])); err != nil {
				return err
			}

			if err := errorOnIncorrectValue(value <= 0, "Simultaneous Calls", "max_call_threshold", "Must be > 0"); err != nil {
				return err
			}

			fraudionConfig.Triggers.SimultaneousCalls.MaxCallThreshold = uint32(value)

		} else { // This value is mandatory!
			// TODO:LOG Log this to Syslog
			fmt.Printf("ERROR: \"max_call_threshold\" value for \"Simultaneous Calls\" Trigger not found")
			return fmt.Errorf("\"max_call_threshold\" value for \"simultaneous calls\" trigger not found")
		}

		// ** Search "minimum_number_length" Key [Optional], defaults to "fraudionConfig.General.DefaultMinimumDestinationNumberLength"
		if utils.StringKeyInMap("minimum_number_length", simCallsTriggerJSONConfig) {

			value, hasCorrectType := simCallsTriggerJSONConfig["minimum_number_length"].(float64)
			if err := errorOnIncorrectType(hasCorrectType, "Simultaneous Calls", "minimum_number_length", "float64", reflect.TypeOf(simCallsTriggerJSONConfig["minimum_number_length"])); err != nil {
				return err
			}

			if err := errorOnIncorrectValue(value <= 0, "Simultaneous Calls", "minimum_number_length", "Must be > 0"); err != nil {
				return err
			}

			fraudionConfig.Triggers.SimultaneousCalls.MinimumNumberLength = uint32(value)

		} else { // This value is mandatory!
			// TODO:LOG Log this to Syslog
			fmt.Printf("WARNING: \"minimum_number_length\" value for \"Simultaneous Calls\" Trigger not found, assuming %d\n", fraudionConfig.General.DefaultMinimumDestinationNumberLength)
			fraudionConfig.Triggers.SimultaneousCalls.MinimumNumberLength = fraudionConfig.General.DefaultMinimumDestinationNumberLength
		}

	} else {
		// TODO:LOG Log this to Syslog
		fmt.Println("WARNING: \"Simultaneous Calls\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
	}

	fmt.Println("******", fraudionConfig.Triggers.SimultaneousCalls)

	//  Dangerous Destinations
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent = JSONConfigs.Triggers.Dangerous_destinations.(map[string]interface{})
	if triggerIsPresent {

		// ** Search "enabled" Key [Optional], defaults to "enabled"
		if utils.StringKeyInMap("enabled", simCallsTriggerJSONConfig) {

			value, hasCorrectType := simCallsTriggerJSONConfig["enabled"].(bool)
			if err := errorOnIncorrectType(hasCorrectType, "Dangerous Destinations", "enabled", "boolean", reflect.TypeOf(simCallsTriggerJSONConfig["enabled"])); err != nil {
				return err
			}

			fraudionConfig.Triggers.DangerousDestinations.Enabled = value

		} else {
			// TODO:LOG Log this to Syslog
			fmt.Println("WARNING: \"Enabled\" value for \"Dangerous Destinations\" Trigger not found, assuming true")
			fraudionConfig.Triggers.DangerousDestinations.Enabled = true
		}

		// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		if utils.StringKeyInMap("check_period", simCallsTriggerJSONConfig) {

			valueString, hasCorrectType := simCallsTriggerJSONConfig["check_period"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Dangerous Destinations", "check_period", "string", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"])); err != nil {
				return err
			}

			valueDuration, errDuration := time.ParseDuration(valueString)
			if err := errorOnIncorrectValue(errDuration != nil, "Dangerous Destinations", "check_period", "Must be a valid duration"); err != nil {
				return err
			}

			errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
			if err := errorOnIncorrectValue(valueDuration < MinimumTriggerCheckPeriodDuration, "Dangerous Destinations", "check_period", errorMessageSuffix); err != nil {
				return err
			}

			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = valueDuration

		} else {
			// TODO:LOG Log this to Syslog
			fmt.Printf("WARNING: \"check_period\" value for \"Dangerous Destinations\" Trigger not found, assuming %s\n", fraudionConfig.General.DefaultTriggerCheckPeriod)
			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = fraudionConfig.General.DefaultTriggerCheckPeriod
		}

		// ** Search "prefix_list" Key [Mandatory + not empty]
		if utils.StringKeyInMap("prefix_list", simCallsTriggerJSONConfig) {

			valueString, hasCorrectType := simCallsTriggerJSONConfig["prefix_list"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Dangerous Destinations", "prefix_list", "string", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"])); err != nil {
				return err
			}

			valueDuration, errDuration := time.ParseDuration(valueString)
			if err := errorOnIncorrectValue(errDuration != nil, "Dangerous Destinations", "prefix_list", "Must be a valid duration"); err != nil {
				return err
			}

			errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
			if err := errorOnIncorrectValue(valueDuration < MinimumTriggerCheckPeriodDuration, "Dangerous Destinations", "prefix_list", errorMessageSuffix); err != nil {
				return err
			}

			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = valueDuration

		} else {
			// TODO:LOG Log this to Syslog
			fmt.Printf("WARNING: \"prefix_list\" value for \"Dangerous Destinations\" Trigger not found, assuming %s\n", fraudionConfig.General.DefaultTriggerCheckPeriod)
			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = fraudionConfig.General.DefaultTriggerCheckPeriod
		}

	} else {
		// TODO:LOG Log this to Syslog
		fmt.Println("WARNING: \"Dangerous Destinations\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
	}
	// Actions Section

	// ActionChains Section

	// Contacts Section

	return nil
}

func errorOnIncorrectType(condition bool, triggerName string, triggerValue string, expectedType string, foundType reflect.Type) error {
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
		errorMessage := fmt.Sprintf("\"%s\" value for \"%s\" Trigger %s is incorrect. ", triggerValue, triggerName, whatItMustBeMessage)
		fmt.Printf("ERROR: %s :(", errorMessage)
		return fmt.Errorf("%s", strings.ToLower(errorMessage))
	}
	return nil
}

func CheckJSONSanityAndLoadEnabledConfigs(fraudionConfig FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) error {
	if utils.StringKeyInMap("enabled", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["enabled"].(bool)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "enabled", "boolean", reflect.TypeOf(simCallsTriggerJSONConfig["enabled"])); err != nil {
			return err
		}

		fraudionConfig.Triggers.SimultaneousCalls.Enabled = value

	} else {
		// TODO:LOG Log this to Syslog
		fmt.Printf("WARNING: \"enabled\" value for \"%s\" Trigger not found, assuming true", triggerName)
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = true
	}

	return nil
}

/*func CheckJSONSanityAndLoadCheckPeriodConfigs()         {

		// ** Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		if utils.StringKeyInMap("check_period", simCallsTriggerJSONConfig) {

			valueString, hasCorrectType := simCallsTriggerJSONConfig["check_period"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Simultaneous Calls", "check_period", "string", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"])); err != nil {
				return err
			}

			valueDuration, errDuration := time.ParseDuration(valueString)
			if err := errorOnIncorrectValue(errDuration != nil, "Simultaneous Calls", "check_period", "Must be a valid duration"); err != nil {
				return err
			}

			errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
			if err := errorOnIncorrectValue(valueDuration < MinimumTriggerCheckPeriodDuration, "Simultaneous Calls", "check_period", errorMessageSuffix); err != nil {
				return err
			}

			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = valueDuration

		} else {
			// TODO:LOG Log this to Syslog
			fmt.Printf("WARNING: \"check_period\" value for \"Simultaneous Calls\" Trigger not found, assuming %s\n", fraudionConfig.General.DefaultTriggerCheckPeriod)
			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = fraudionConfig.General.DefaultTriggerCheckPeriod
		}
}
func CheckJSONSanityAndLoadMinimumNumberLengthConfigs() {}*/
