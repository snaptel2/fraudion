package config

import (
	"time"
)

// Types for Loaded Config

// FraudionConfig2 ...
type FraudionConfig2 struct {
	General      General2
	Triggers     Triggers
	Actions      Actions2
	ActionChains ActionChains2
	DataGroups   DataGroups2
}

// General2 ...
type General2 struct {
	MonitoredSoftware                     string
	CDRsSource                            string
	DefaultTriggerCheckPeriod             time.Duration
	DefaultHitThreshold                   uint32
	DefaultMinimumDestinationNumberLength uint32
	DefaultActionChainHoldoffPeriod       time.Duration
	DefaultActionChainRunCount            uint32
}

// Triggers2 ...
type Triggers2 struct {
	SimultaneousCalls     triggerSimultaneousCalls2
	DangerousDestinations triggerDangerousDestinations2
	ExpectedDestinations  triggerExpectedDestinations2
	SmallDurationCalls    triggerSmallCallDurations2
}

type triggerSimultaneousCalls2 struct {
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	ActionChainName        string
	LastActionChainRunTime time.Time
}

type triggerDangerousDestinations2 struct {
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	PrefixList             []string
	ActionChainName        string
	LastActionChainRunTime time.Time
}

type triggerExpectedDestinations2 struct {
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	PrefixList             []string
	ActionChainName        string
	LastActionChainRunTime time.Time
}

type triggerSmallCallDurations2 struct {
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	ActionChainName        string
	DurationThreshold      time.Duration
	LastActionChainRunTime time.Time
}

// Actions2 ...
type Actions2 struct {
	Email         actionEmail2
	Call          actionCall2
	HTTP          actionHTTP2
	LocalCommands actionLocalCommands2
}

type actionEmail2 struct {
	Enabled  bool
	Username string
	Password string
	Message  string
}

type actionCall2 struct {
	Enabled bool
}

type actionHTTP2 struct {
	Enabled bool
}

type actionLocalCommands2 struct {
	Enabled bool
}

// ActionChains2 ...
type ActionChains2 struct {
	List map[string][]actionChainAction2
}

type actionChainAction2 struct {
	ActionName     string
	DataGroupNames []string
}

// DataGroups2 ...
type DataGroups2 struct {
	List map[string]DataGroup2
}

// DataGroup2 ...
type DataGroup2 struct {
	ForActions       []string // TODO: Maybe this could serve to validate which of the fields bellow should be confirmed to exist
	PhoneNumber      string
	EmailAddress     string
	HTTPURL          string
	HTTPMethod       string
	HTTPParameters   map[string]string
	CommandName      string
	CommandArguments string
}

// Types for JSON Config Unmarshaling

// FraudionConfigJSON2 ...
type FraudionConfigJSON2 struct {
	General      GeneralJSON2
	Triggers     TriggersJSON2
	Actions      ActionsJSON2
	ActionChains ActionChainsJSON2
	DataGroups   DataGroupsJSON2
}

// GeneralJSON2 ...
type GeneralJSON2 struct {
	MonitoredSoftware                     string `json:"monitored_software"`
	CDRsSource                            string `json:"cdrs_source"`
	DefaultTriggerExecuteInterval         string `json:"default_trigger_execute_interval"`
	DefaultHitThreshold                   uint32 `json:"default_hit_threshold"`
	DefaultMinimumDestinationNumberLength uint32 `json:"default_minimum_destination_number_length"`
	DefaultActionChainHoldoffPeriod       string `json:"default_action_chain_holdoff_period"`
	DefaultActionChainRunCount            uint32 `json:"default_action_chain_run_count"`
}

// TriggersJSON2 ...
type TriggersJSON2 struct {
	SimultaneousCalls     triggerSimultaneousCallsJSON2     `json:"simultaneous_calls"`
	DangerousDestinations triggerDangerousDestinationsJSON2 `json:"dangerous_destinations"`
	ExpectedDestinations  triggerExpectedDestinationsJSON2  `json:"expected_destinations"`
	SmallDurationCalls    triggerSmallCallDurationsJSON2    `json:"small_duration_calls"`
}

type triggerSimultaneousCallsJSON2 struct {
	Enabled             bool   `json:"enabled"`
	ExecuteInterval     string `json:"execute_interval"`
	HitThreshold        uint32 `json:"hit_threshold"`
	MinimumNumberLength uint32 `json:"minimum_number_length"`
	ActionChainName     string `json:"action_chain_name"`
}

type triggerDangerousDestinationsJSON2 struct {
	Enabled              bool     `json:"enabled"`
	ExecuteInterval      string   `json:"execute_interval"`
	HitThreshold         uint32   `json:"hit_threshold"`
	MinimumNumberLength  uint32   `json:"minimum_number_length"`
	ActionChainName      string   `json:"action_chain_name"`
	ConsiderCDRsFromLast string   `json:"consider_cdrs_from_last"`
	PrefixList           []string `json:"prefix_list"`
	MatchRegex           string   `json:"match_regex"`
	IgnoreRegex          string   `json:"ignore_regex"`
}

type triggerExpectedDestinationsJSON2 struct {
	Enabled              bool     `json:"enabled"`
	CheckPeriod          string   `json:"execute_interval"`
	HitThreshold         uint32   `json:"hit_threshold"`
	MinimumNumberLength  uint32   `json:"minimum_number_length"`
	ActionChainName      string   `json:"action_chain_name"`
	ConsiderCDRsFromLast string   `json:"consider_cdrs_from_last"`
	PrefixList           []string `json:"prefix_list"`
	MatchRegex           string   `json:"match_regex"`
	IgnoreRegex          string   `json:"ignore_regex"`
}

type triggerSmallCallDurationsJSON2 struct {
	Enabled              bool   `json:"enabled"`
	CheckPeriod          string `json:"execute_interval"`
	HitThreshold         uint32 `json:"hit_threshold"`
	MinimumNumberLength  uint32 `json:"minimum_number_length"`
	ActionChainName      string `json:"action_chain_name"`
	ConsiderCDRsFromLast string `json:"consider_cdrs_from_last"`
	DurationThreshold    string `json:"duration_threshold"`
}

// ActionsJSON2 ...
type ActionsJSON2 struct {
	Email         actionEmailJSON2         `json:"email"`
	Call          actionCallJSON2          `json:"call"`
	HTTP          actionHTTPJSON2          `json:"http"`
	LocalCommands actionLocalCommandsJSON2 `json:"local_commands"`
}

type actionEmailJSON2 struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"gmail_username"`
	Password string `json:"gmail_password"`
	Message  string `json:"message"`
}

type actionCallJSON2 struct {
	Enabled bool `json:"enabled"`
}

type actionHTTPJSON2 struct {
	Enabled bool `json:"enabled"`
}

type actionLocalCommandsJSON2 struct {
	Enabled bool `json:"enabled"`
}

// ActionChainsJSON2 ...
type ActionChainsJSON2 struct {
	List map[string][]actionChainActionJSON2 `json:"list"`
}

type actionChainActionJSON2 struct {
	ActionName     string   `json:"action"`
	DataGroupNames []string `json:"data_groups"`
}

// DataGroupsJSON2 ...
type DataGroupsJSON2 struct {
	List map[string]DataGroupJSON2 `json:"list"`
}

// DataGroupJSON2 ...
type DataGroupJSON2 struct {
	ForActions       []string          `json:"for_actions"` // TODO: Maybe this could serve to validate which of the fields bellow should be confirmed to exist
	PhoneNumber      string            `json:"phone_number"`
	EmailAddress     string            `json:"email_address"`
	HTTPURL          string            `json:"http_url"`
	HTTPMethod       string            `json:"http_method"`
	HTTPParameters   map[string]string `json:"http_parameters"`
	CommandName      string            `json:"command_name"`
	CommandArguments string            `json:"command_arguments"`
}
