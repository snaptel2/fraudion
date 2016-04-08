package types

import (
	"io"
	"log"
	"time"
)

// FraudionGlobal The "Type"
type FraudionGlobal struct {
	StartUpTime time.Time
	Debug       bool
	LogTrace    *log.Logger
	LogInfo     *log.Logger
	LogWarning  *log.Logger
	LogError    *log.Logger
	Configs     *FraudionConfig
	State       *FraudionState
}

// SetupLogging ...
func (fraudion *FraudionGlobal) SetupLogging(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {

	fraudion.LogTrace = log.New(traceHandle, "FRAUDION TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	fraudion.LogInfo = log.New(infoHandle, "FRAUDION INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	fraudion.LogWarning = log.New(warningHandle, "FRAUDION WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	fraudion.LogError = log.New(errorHandle, "FRAUDION ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

}

// SetupState ...
func (fraudion *FraudionGlobal) SetupState() {
	fraudion.State = new(FraudionState)
}

// SetupConfigs ...
func (fraudion *FraudionGlobal) SetupConfigs() {
	fraudion.Configs = new(FraudionConfig)
}

// Fraudion ...
var Fraudion *FraudionGlobal

// Types for Loaded Config

// FraudionConfig ...
type FraudionConfig struct {
	General      General
	CDRsSources  interface{}
	Triggers     Triggers
	Actions      Actions
	ActionChains ActionChains
	DataGroups   DataGroups
}

// FraudionState ...
type FraudionState struct {
	StateTriggers StateTriggers
}

// General ...
type General struct {
	MonitoredSoftware                     string
	CDRsSource                            string
	DefaultTriggerExecuteInterval         time.Duration
	DefaultHitThreshold                   uint32
	DefaultMinimumDestinationNumberLength uint32
	DefaultActionChainHoldoffPeriod       time.Duration
	DefaultActionChainRunCount            uint32
}

// Triggers ...
type Triggers struct {
	SimultaneousCalls     triggerSimultaneousCalls
	DangerousDestinations triggerDangerousDestinations
	ExpectedDestinations  triggerExpectedDestinations
	SmallDurationCalls    triggerSmallCallDurations
}

// StateTriggers ...
type StateTriggers struct {
	StateSimultaneousCalls     stateSimultaneousCalls
	StateDangerousDestinations stateDangerousDestinations
	StateExpectedDestinations  stateExpectedDestinations
	StateSmallDurationCalls    stateSmallCallDurations
}

type triggerSimultaneousCalls struct {
	Enabled                  bool
	ExecuteInterval          time.Duration
	HitThreshold             uint32
	MinimumNumberLength      uint32
	ActionChainName          string
	ActionChainHoldoffPeriod uint32
	MaxActionChainRunCount   uint32
}

type stateSimultaneousCalls struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type triggerDangerousDestinations struct {
	Enabled                  bool
	ExecuteInterval          time.Duration
	HitThreshold             uint32
	MinimumNumberLength      uint32
	ActionChainName          string
	ActionChainHoldoffPeriod uint32
	MaxActionChainRunCount   uint32
	ConsiderCDRsFromLast     string
	PrefixList               []string
	MatchRegex               string
	IgnoreRegex              string
}

type stateDangerousDestinations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type triggerExpectedDestinations struct {
	Enabled                  bool
	ExecuteInterval          time.Duration
	HitThreshold             uint32
	MinimumNumberLength      uint32
	ActionChainName          string
	ActionChainHoldoffPeriod uint32
	MaxActionChainRunCount   uint32
	ConsiderCDRsFromLast     string
	PrefixList               []string
	MatchRegex               string
	IgnoreRegex              string
}

type stateExpectedDestinations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type triggerSmallCallDurations struct {
	Enabled                  bool
	ExecuteInterval          time.Duration
	HitThreshold             uint32
	MinimumNumberLength      uint32
	ActionChainName          string
	ActionChainHoldoffPeriod uint32
	MaxActionChainRunCount   uint32
	ConsiderCDRsFromLast     string
	DurationThreshold        time.Duration
}

type stateSmallCallDurations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

// Actions ...
type Actions struct {
	Email         actionEmail
	Call          actionCall
	HTTP          actionHTTP
	LocalCommands actionLocalCommands
}

type actionEmail struct {
	Enabled  bool
	Username string
	Password string
	Message  string
}

type actionCall struct {
	Enabled bool
}

type actionHTTP struct {
	Enabled bool
}

type actionLocalCommands struct {
	Enabled bool
}

// ActionChains ...
type ActionChains struct {
	List map[string][]actionChainAction
}

type actionChainAction struct {
	ActionName     string   `json:"action"`
	DataGroupNames []string `json:"data_groups"`
}

// DataGroups ...
type DataGroups struct {
	List map[string]DataGroup
}

// DataGroup ...
type DataGroup struct {
	PhoneNumber      string            `json:"phone_number"`
	EmailAddress     string            `json:"email_address"`
	HTTPURL          string            `json:"http_url"`
	HTTPMethod       string            `json:"http_method"`
	HTTPParameters   map[string]string `json:"http_parameters"`
	CommandName      string            `json:"command_name"`
	CommandArguments string            `json:"command_arguments"`
}

// Types for JSON Config Unmarshaling

// FraudionConfigJSON ...
type FraudionConfigJSON struct {
	General      GeneralJSON
	CDRsSources  map[string]map[string]string `json:"cdrs_sources"`
	Triggers     TriggersJSON
	Actions      ActionsJSON
	ActionChains ActionChainsJSON
	DataGroups   DataGroupsJSON
}

// GeneralJSON ...
type GeneralJSON struct {
	MonitoredSoftware                     string `json:"monitored_software"`
	CDRsSource                            string `json:"cdrs_source"`
	DefaultTriggerExecuteInterval         string `json:"default_trigger_execute_interval"`
	DefaultHitThreshold                   uint32 `json:"default_hit_threshold"`
	DefaultMinimumDestinationNumberLength uint32 `json:"default_minimum_destination_number_length"`
	DefaultActionChainHoldoffPeriod       string `json:"default_action_chain_holdoff_period"`
	DefaultActionChainRunCount            uint32 `json:"default_action_chain_run_count"`
}

// TriggersJSON ...
type TriggersJSON struct {
	SimultaneousCalls     triggerSimultaneousCallsJSON     `json:"simultaneous_calls"`
	DangerousDestinations triggerDangerousDestinationsJSON `json:"dangerous_destinations"`
	ExpectedDestinations  triggerExpectedDestinationsJSON  `json:"expected_destinations"`
	SmallDurationCalls    triggerSmallCallDurationsJSON    `json:"small_duration_calls"`
}

type triggerSimultaneousCallsJSON struct {
	Enabled                  bool   `json:"enabled"`
	ExecuteInterval          string `json:"execute_interval"`
	HitThreshold             uint32 `json:"hit_threshold"`
	MinimumNumberLength      uint32 `json:"minimum_number_length"`
	ActionChainName          string `json:"action_chain_name"`
	ActionChainHoldoffPeriod uint32 `json:"action_chain_holdoff_period"`
	MaxActionChainRunCount   uint32 `json:"action_chain_run_count"`
}

type triggerDangerousDestinationsJSON struct {
	Enabled                  bool     `json:"enabled"`
	ExecuteInterval          string   `json:"execute_interval"`
	HitThreshold             uint32   `json:"hit_threshold"`
	MinimumNumberLength      uint32   `json:"minimum_number_length"`
	ActionChainName          string   `json:"action_chain_name"`
	ActionChainHoldoffPeriod uint32   `json:"action_chain_holdoff_period"`
	MaxActionChainRunCount   uint32   `json:"action_chain_run_count"`
	ConsiderCDRsFromLast     string   `json:"consider_cdrs_from_last"`
	PrefixList               []string `json:"prefix_list"`
	MatchRegex               string   `json:"match_regex"`
	IgnoreRegex              string   `json:"ignore_regex"`
}

type triggerExpectedDestinationsJSON struct {
	Enabled                  bool     `json:"enabled"`
	ExecuteInterval          string   `json:"execute_interval"`
	HitThreshold             uint32   `json:"hit_threshold"`
	MinimumNumberLength      uint32   `json:"minimum_number_length"`
	ActionChainName          string   `json:"action_chain_name"`
	ActionChainHoldoffPeriod uint32   `json:"action_chain_holdoff_period"`
	MaxActionChainRunCount   uint32   `json:"action_chain_run_count"`
	ConsiderCDRsFromLast     string   `json:"consider_cdrs_from_last"`
	PrefixList               []string `json:"prefix_list"`
	MatchRegex               string   `json:"match_regex"`
	IgnoreRegex              string   `json:"ignore_regex"`
}

type triggerSmallCallDurationsJSON struct {
	Enabled                  bool   `json:"enabled"`
	ExecuteInterval          string `json:"execute_interval"`
	HitThreshold             uint32 `json:"hit_threshold"`
	MinimumNumberLength      uint32 `json:"minimum_number_length"`
	ActionChainName          string `json:"action_chain_name"`
	ActionChainHoldoffPeriod uint32 `json:"action_chain_holdoff_period"`
	MaxActionChainRunCount   uint32 `json:"action_chain_run_count"`
	ConsiderCDRsFromLast     string `json:"consider_cdrs_from_last"`
	DurationThreshold        string `json:"duration_threshold"`
}

// ActionsJSON ...
type ActionsJSON struct {
	Email         actionEmailJSON         `json:"email"`
	Call          actionCallJSON          `json:"call"`
	HTTP          actionHTTPJSON          `json:"http"`
	LocalCommands actionLocalCommandsJSON `json:"local_commands"`
}

type actionEmailJSON struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"gmail_username"`
	Password string `json:"gmail_password"`
	Message  string `json:"message"`
}

type actionCallJSON struct {
	Enabled bool `json:"enabled"`
}

type actionHTTPJSON struct {
	Enabled bool `json:"enabled"`
}

type actionLocalCommandsJSON struct {
	Enabled bool `json:"enabled"`
}

// ActionChainsJSON ...
type ActionChainsJSON struct {
	List map[string][]actionChainAction `json:"list"`
}

// DataGroupsJSON ...
type DataGroupsJSON struct {
	List map[string]DataGroup `json:"list"`
}
