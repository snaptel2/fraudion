package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/fraudion/utils"
)

const (
	MINIMUM_TRIGGER_CHECK_PERIOD      = "5m"
	MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD = "5m"
)

var (
	SUPPORTED_SOFTWARE = []string{"*ast_el_2.3", "*ast_alone_1.8"}
	CDRS_SOURCES       = []string{"*db_mysql"}
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
		// TODO Log this to Syslog
		fmt.Println("Could not find configured \"monitored_software\" in supported software list. :(")
		return fmt.Errorf("could not find configured \"monitored_software\" in supported software list")
	}

	fraudionConfig.General.MonitoredSoftware = JSONConfigs.General.Monitored_software

	//  CDRs Source
	found = utils.StringInStringsSlice(JSONConfigs.General.Cdrs_source, CDRS_SOURCES)
	if found == false {
		// TODO Log this to Syslog
		fmt.Println("Could not find configured \"cdrs_source\" in supported software list. :(")
		return fmt.Errorf("could not find configured \"cdrs_source\" in supported software list")
	}

	fraudionConfig.General.CdrsSource = JSONConfigs.General.Cdrs_source

	//  Default_trigger_check_period
	defaultTriggerCheckPeriodDuration, err := time.ParseDuration(JSONConfigs.General.Default_trigger_check_period)
	if err != nil {
		// TODO Log this to Syslog
		fmt.Println("ERROR: It seems that the value in \"default_trigger_check_period\" is not a valid duration. :(")
		return fmt.Errorf("it seems that the value in \"default_trigger_check_period\" is not a valid duration")
	}

	MinimumTriggerCheckPeriodDuration, err := time.ParseDuration(MINIMUM_TRIGGER_CHECK_PERIOD)
	if defaultTriggerCheckPeriodDuration < MinimumTriggerCheckPeriodDuration {
		// TODO Log this to Syslog
		fmt.Printf("ERROR: \"default_trigger_check_period\" is too small. Minimum value is \"%s\". :(", MINIMUM_TRIGGER_CHECK_PERIOD)
		return fmt.Errorf("\"default_trigger_check_period\" is too small. minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
	}
	if err != nil {
		// TODO Log this to Syslog
		fmt.Printf("ERROR: Internal: There seems to be an issue with the definition of the constant MINIMUM_TRIGGER_CHECK_PERIOD.")
		return fmt.Errorf("internal: there seems to be an issue with the definition of the constant MINIMUM_TRIGGER_CHECK_PERIOD")
	}

	fraudionConfig.General.DefaultTriggerCheckPeriod = defaultTriggerCheckPeriodDuration

	//  Default_action_chain_sleep_period
	defaultActionChainSleepPeriodDuration, err := time.ParseDuration(JSONConfigs.General.Default_action_chain_sleep_period)
	if err != nil {
		// TODO Log this to Syslog
		fmt.Println("ERROR: It seems that the value in \"default_action_chain_sleep_period\" is not a valid duration. :(")
		return fmt.Errorf("it seems that the value in \"default_action_chain_sleep_period\" is not a valid duration")
	}

	MaximumActionChainSleepPeriodDuration, err := time.ParseDuration(MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
	if defaultActionChainSleepPeriodDuration > MaximumActionChainSleepPeriodDuration {
		// TODO Log this to Syslog
		fmt.Printf("ERROR: \"default_trigger_check_period\" is too big. Maximum value is \"%s\". :(", MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
		return fmt.Errorf("\"default_trigger_check_period\" is too big. maximum value is \"%s\"", MAXIMUM_ACTION_CHAIN_SLEEP_PERIOD)
	}
	if err != nil {
		// TODO Log this to Syslog
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

	// Triggers Section
	//  Simultaneous calls
	//   Check if Trigger is present!
	simCallsTriggerJSONConfig, triggerIsPresent := JSONConfigs.Triggers.Simultaneous_calls.(map[string]interface{})
	if triggerIsPresent {

		// Search "enabled" Key [Optional], defaults to "enabled"
		if utils.StringKeyInMap("enabled", simCallsTriggerJSONConfig) {

			value, hasCorrectType := simCallsTriggerJSONConfig["enabled"].(bool)
			if !hasCorrectType {
				// TODO Log this to Syslog
				fmt.Printf("ERROR: \"enabled\" value for \"Simultaneous Calls\" Trigger not of correct type, boolean expected found %s. :(", reflect.TypeOf(simCallsTriggerJSONConfig["enabled"]))
				return fmt.Errorf("\"enabled\" value for \"simultaneous calls\" trigger not of correct type, boolean expected found %s", reflect.TypeOf(simCallsTriggerJSONConfig["enabled"]))
			}

			fraudionConfig.Triggers.SimultaneousCalls.Enabled = value

		} else {
			// TODO Log this to Syslog
			fmt.Println("WARNING: \"Enabled\" value for \"Simultaneous Calls\" Trigger not found, assuming true")
			fraudionConfig.Triggers.SimultaneousCalls.Enabled = true
		}

		// Search "check_period" Key [Optional], defaults to "fraudionConfig.General.DefaultTriggerCheckPeriod"
		if utils.StringKeyInMap("check_period", simCallsTriggerJSONConfig) {

			valueString, hasCorrectType := simCallsTriggerJSONConfig["check_period"].(string)
			if !hasCorrectType {
				// TODO Log this to Syslog
				fmt.Printf("ERROR: \"check_period\" value for \"Simultaneous Calls\" Trigger not of correct type, string expected found %s. :(\n", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"]))
				return fmt.Errorf("\"check_period\" value for \"simultaneous calls\" trigger not of correct type, string expected found %s", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"]))
			}

			value, err := time.ParseDuration(valueString)
			if err != nil {
				// TODO Log this to Syslog
				fmt.Println("ERROR: \"check_period\" value for \"Simultaneous Calls\" Trigger is not a valid duration. :(")
				return fmt.Errorf("\"check_period\" value for \"simultaneous calls\" trigger is not a valid duration")
			}

			if value < MinimumTriggerCheckPeriodDuration {
				// TODO Log this to Syslog
				fmt.Printf("ERROR: \"check_period\" value for \"Simultaneous Calls\" Trigger is too small. Minimum value is \"%s\" :(\n", MINIMUM_TRIGGER_CHECK_PERIOD)
				return fmt.Errorf("\"check_period\" value for \"simultaneous calls\" trigger is too small. minimum value is \"%s\"", MINIMUM_TRIGGER_CHECK_PERIOD)
			}

			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = value

		} else {
			// TODO Log this to Syslog
			fmt.Printf("WARNING: \"check_period\" value for \"Simultaneous Calls\" Trigger not found, assuming %s\n", fraudionConfig.General.DefaultTriggerCheckPeriod)
			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = fraudionConfig.General.DefaultTriggerCheckPeriod
		}

		// Search "max_call_threshold" Key [Mandatory]
		if utils.StringKeyInMap("max_call_threshold", simCallsTriggerJSONConfig) {

			value, hasCorrectType := simCallsTriggerJSONConfig["max_call_threshold"].(float64)
			if !hasCorrectType {
				// TODO Log this to Syslog
				fmt.Printf("ERROR: \"max_call_threshold\" value for \"Simultaneous Calls\" Trigger not of correct type, float64 expected found %s. :(\n", reflect.TypeOf(simCallsTriggerJSONConfig["max_call_threshold"]))
				return fmt.Errorf("\"max_call_threshold\" value for \"simultaneous calls\" trigger not of correct type, float64 expected found %s", reflect.TypeOf(simCallsTriggerJSONConfig["max_call_threshold"]))
			}

			if value <= 0 {
				// TODO Log this to Syslog
				fmt.Printf("ERROR: \"max_call_threshold\" value for \"Simultaneous Calls\" Trigger must be > 0")
				return fmt.Errorf("\"max_call_threshold\" value for \"simultaneous calls\" trigger must be > 0")
			}

			fraudionConfig.Triggers.SimultaneousCalls.MaxCallThreshold = uint32(value)

		} else { // This value is mandatory!
			// TODO Log this to Syslog
			fmt.Printf("ERROR: \"max_call_threshold\" value for \"Simultaneous Calls\" Trigger not found")
			return fmt.Errorf("\"max_call_threshold\" value for \"simultaneous calls\" trigger not found")
		}

		// Search "minimum_number_length" Key [Optional], defaults to "fraudionConfig.General.DefaultMinimumDestinationNumberLength"
		if utils.StringKeyInMap("minimum_number_length", simCallsTriggerJSONConfig) {

			value, hasCorrectType := simCallsTriggerJSONConfig["minimum_number_length"].(float64)
			if !hasCorrectType {
				// TODO Log this to Syslog
				fmt.Printf("ERROR: \"minimum_number_length\" value for \"Simultaneous Calls\" Trigger not of correct type, float64 expected found %s. :(\n", reflect.TypeOf(simCallsTriggerJSONConfig["minimum_number_length"]))
				return fmt.Errorf("\"minimum_number_length\" value for \"simultaneous calls\" trigger not of correct type, float64 expected found %s", reflect.TypeOf(simCallsTriggerJSONConfig["minimum_number_length"]))
			}

			if value <= 0 {
				// TODO Log this to Syslog
				fmt.Printf("ERROR: \"minimum_number_length\" value for \"Simultaneous Calls\" Trigger must be > 0")
				return fmt.Errorf("\"minimum_number_length\" value for \"simultaneous calls\" trigger must be > 0")
			}

			fraudionConfig.Triggers.SimultaneousCalls.MinimumNumberLength = uint32(value)

		} else { // This value is mandatory!
			// TODO Log this to Syslog
			fmt.Printf("WARNING: \"minimum_number_length\" value for \"Simultaneous Calls\" Trigger not found, assuming %d\n", fraudionConfig.General.DefaultMinimumDestinationNumberLength)
			fraudionConfig.Triggers.SimultaneousCalls.MinimumNumberLength = fraudionConfig.General.DefaultMinimumDestinationNumberLength
		}

	} else {
		// TODO Log this to Syslog
		fmt.Println("WARNING: \"Simultaneous Calls\" Trigger config not present, disabling!")
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false
	}

	fmt.Println("******", fraudionConfig.Triggers.SimultaneousCalls)

	// Actions Section

	// ActionChains Section

	// Contacts Section

	return nil
}
