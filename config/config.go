package config

import (
	"fmt"
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
	CheckPeriod      time.Duration
	MaxCallThreshold uint32
}

type triggerDangerousDestinations struct {
	CheckPeriod  time.Duration
	PrefixList   []string
	HitThreshold uint32
}

type triggerExpectedDestinations struct {
	CheckPeriod  time.Duration
	PrefixList   []string
	HitThreshold uint32
}

type triggerSmallCallDurations struct {
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
	/*
		   Simultaneous_calls     map[string]interface{}
		   Dangerous_destinations map[string]interface{}
		   Expected_destinations  map[string]interface{}
		   Small_duration_calls   map[string]interface{}

			 e.g.

			 "expected_destinations": {
			 	"check_period": "6m", // [Optional] Defaults to what's in "general" section
			 	"prefix_list": ["244"],
			 	"hit_threshold": 5 // "this" number of calls in "check_period" will trigger!
			 }

	*/

	// Actions Section

	// ActionChains Section

	// Contacts Section

	return nil
}
