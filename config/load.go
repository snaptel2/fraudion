package config

import (
	"fmt"
	"time"

	"github.com/andmar/fraudion/logger"
	"github.com/andmar/fraudion/utils"
)

const (
	constMinimumTriggerExecuteInterval   = "1m"
	constMaximumActionChainHoldoffPeriod = "5m"
)

var (
	constSupportedSoftware2   = []string{"*ast_elastix_2.3", "*ast_1.8"}
	constSupportedCDRSources2 = []string{"*db_mysql"}
)

// Load ...
func Load(configsJSON *FraudionConfigJSON) (*FraudionConfig, error) {

	configs := new(FraudionConfig)

	fmt.Println(configs)

	logger.Log.Write(logger.ConstLoggerLevelInfo, "Validating and Loading configurations...", false)

	// ** General Section

	// * MonitoredSoftware
	if configsJSON.General.MonitoredSoftware == "" {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"monitored_software\" value in section \"general\" missing OR is empty."), true)
	}
	found := utils.StringInStringsSlice(configsJSON.General.MonitoredSoftware, constSupportedSoftware2)
	if !found {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"monitored_software\" value in section \"general\" must be one of %s", constSupportedSoftware2), true)
	}
	configs.General.MonitoredSoftware = configsJSON.General.MonitoredSoftware

	// * DefaultTriggerExecuteInterval
	if configsJSON.General.DefaultTriggerExecuteInterval == "" {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_trigger_execute_interval\" value in section \"general\" missing OR is empty."), true)
	}
	durationDefaultTriggerExecuteInterval, err := time.ParseDuration(configsJSON.General.DefaultTriggerExecuteInterval)
	if err != nil {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_trigger_execute_interval\" value in section \"general\" is not a valid duration. Must be a parseable duration, a number followed by one of: \"s\", \"m\" or \"h\" for \"seconds\", \"minutes\" and \"hours\" respectively"), true)
	}
	durationConstMinimumTriggerExecuteInterval, err := time.ParseDuration(constMinimumTriggerExecuteInterval)
	if err != nil {
		return nil, utils.DebugLogAndGetError("(Internal) There seems to be an issue with the definition of constMinimumTriggerExecuteInterval2", true)
	}
	if durationDefaultTriggerExecuteInterval < durationConstMinimumTriggerExecuteInterval {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_trigger_execute_interval\" value in section \"general\" is too small. Value must be > %s", constMinimumTriggerExecuteInterval), true)
	}

	configs.General.DefaultTriggerExecuteInterval = durationDefaultTriggerExecuteInterval

	// * DefaultHitThreshold
	if configsJSON.General.DefaultHitThreshold == 0 {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_hit_threshold\" value in section \"general\" missing OR is 0."), true)
	}

	configs.General.DefaultHitThreshold = configsJSON.General.DefaultHitThreshold

	// * DefaultMinimumDestinationNumberLength
	if configsJSON.General.DefaultMinimumDestinationNumberLength == 0 {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_minimum_destination_number_length\" value in section \"general\" missing OR is 0."), true)
	}

	configs.General.DefaultMinimumDestinationNumberLength = configsJSON.General.DefaultMinimumDestinationNumberLength

	// * DefaultActionChainHoldoffPeriod
	if configsJSON.General.DefaultActionChainHoldoffPeriod == "" {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_holdoff_period\" value in section \"general\" missing OR empty."), true)
	}
	durationDefaultActionChainHoldoffPeriod, err := time.ParseDuration(configsJSON.General.DefaultActionChainHoldoffPeriod)
	if err != nil {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_holdoff_period\" value in section \"general\" is not a valid duration. Must be a parseable duration, a number followed by one of: \"s\", \"m\" or \"h\" for \"seconds\", \"minutes\" and \"hours\" respectively"), true)
	}
	durationConstMaximumActionChainHoldoffPeriod, err := time.ParseDuration(constMaximumActionChainHoldoffPeriod)
	if err != nil {
		return nil, utils.DebugLogAndGetError("(Internal) There seems to be an issue with the definition of constMaximumActionChainHoldoffPeriod", true)
	}
	if durationDefaultActionChainHoldoffPeriod > durationConstMaximumActionChainHoldoffPeriod {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_holdoff_period\" value in section \"general\" is too small. Value must be > %s", constMaximumActionChainHoldoffPeriod), true)
	}

	configs.General.DefaultActionChainHoldoffPeriod = durationDefaultActionChainHoldoffPeriod

	// * DefaultActionChainRunCount
	if configsJSON.General.DefaultActionChainRunCount == 0 {
		return nil, utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_run_count\" value in section \"general\" missing OR 0."), true)
	}

	configs.General.DefaultActionChainRunCount = configsJSON.General.DefaultActionChainRunCount

	// TODO: From this point the validateOnly flag is not yet checked and used to do something

	// ** CDRsSource Section

	configs.CDRsSources = configsJSON.CDRsSources

	// ** Triggers Section

	// * SimultaneousCalls
	configs.Triggers.SimultaneousCalls.Enabled = configsJSON.Triggers.SimultaneousCalls.Enabled
	if configs.Triggers.SimultaneousCalls.Enabled {
		configs.Triggers.SimultaneousCalls.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.SimultaneousCalls.ExecuteInterval)
		if err != nil {
			return nil, utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.SimultaneousCalls.HitThreshold = configsJSON.Triggers.SimultaneousCalls.HitThreshold
		configs.Triggers.SimultaneousCalls.MinimumNumberLength = configsJSON.Triggers.SimultaneousCalls.MinimumNumberLength
		configs.Triggers.SimultaneousCalls.ActionChainName = configsJSON.Triggers.SimultaneousCalls.ActionChainName
		configs.Triggers.SimultaneousCalls.MaxActionChainRunCount = configsJSON.Triggers.SimultaneousCalls.MaxActionChainRunCount
		//state.StateTriggers.StateSimultaneousCalls.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// * DangerousDestinations
	configs.Triggers.DangerousDestinations.Enabled = configsJSON.Triggers.DangerousDestinations.Enabled
	if configs.Triggers.DangerousDestinations.Enabled {
		configs.Triggers.DangerousDestinations.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.DangerousDestinations.ExecuteInterval)
		if err != nil {
			return nil, utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.DangerousDestinations.HitThreshold = configsJSON.Triggers.DangerousDestinations.HitThreshold
		configs.Triggers.DangerousDestinations.MinimumNumberLength = configsJSON.Triggers.DangerousDestinations.MinimumNumberLength
		configs.Triggers.DangerousDestinations.ActionChainName = configsJSON.Triggers.DangerousDestinations.ActionChainName
		configs.Triggers.DangerousDestinations.MaxActionChainRunCount = configsJSON.Triggers.DangerousDestinations.MaxActionChainRunCount
		configs.Triggers.DangerousDestinations.PrefixList = configsJSON.Triggers.DangerousDestinations.PrefixList
		configs.Triggers.DangerousDestinations.ConsiderCDRsFromLast = configsJSON.Triggers.DangerousDestinations.ConsiderCDRsFromLast
		configs.Triggers.DangerousDestinations.MatchRegex = configsJSON.Triggers.DangerousDestinations.MatchRegex
		configs.Triggers.DangerousDestinations.IgnoreRegex = configsJSON.Triggers.DangerousDestinations.IgnoreRegex
		//state.StateTriggers.StateDangerousDestinations.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// * ExpectedDestinations
	configs.Triggers.ExpectedDestinations.Enabled = configsJSON.Triggers.ExpectedDestinations.Enabled
	if configs.Triggers.ExpectedDestinations.Enabled {
		configs.Triggers.ExpectedDestinations.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.ExpectedDestinations.ExecuteInterval)
		if err != nil {
			return nil, utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.ExpectedDestinations.HitThreshold = configsJSON.Triggers.ExpectedDestinations.HitThreshold
		configs.Triggers.ExpectedDestinations.MinimumNumberLength = configsJSON.Triggers.ExpectedDestinations.MinimumNumberLength
		configs.Triggers.ExpectedDestinations.ActionChainName = configsJSON.Triggers.ExpectedDestinations.ActionChainName
		configs.Triggers.ExpectedDestinations.MaxActionChainRunCount = configsJSON.Triggers.ExpectedDestinations.MaxActionChainRunCount
		configs.Triggers.ExpectedDestinations.PrefixList = configsJSON.Triggers.ExpectedDestinations.PrefixList
		configs.Triggers.ExpectedDestinations.ConsiderCDRsFromLast = configsJSON.Triggers.ExpectedDestinations.ConsiderCDRsFromLast
		configs.Triggers.ExpectedDestinations.MatchRegex = configsJSON.Triggers.ExpectedDestinations.MatchRegex
		configs.Triggers.ExpectedDestinations.IgnoreRegex = configsJSON.Triggers.ExpectedDestinations.IgnoreRegex
		//state.StateTriggers.StateExpectedDestinations.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// * SmallDurationCalls
	configs.Triggers.SmallDurationCalls.Enabled = configsJSON.Triggers.ExpectedDestinations.Enabled
	if configs.Triggers.SmallDurationCalls.Enabled {
		configs.Triggers.SmallDurationCalls.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.ExpectedDestinations.ExecuteInterval)
		if err != nil {
			return nil, utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.SmallDurationCalls.HitThreshold = configsJSON.Triggers.SmallDurationCalls.HitThreshold
		configs.Triggers.SmallDurationCalls.MinimumNumberLength = configsJSON.Triggers.SmallDurationCalls.MinimumNumberLength
		configs.Triggers.SmallDurationCalls.ActionChainName = configsJSON.Triggers.SmallDurationCalls.ActionChainName
		configs.Triggers.SmallDurationCalls.MaxActionChainRunCount = configsJSON.Triggers.SmallDurationCalls.MaxActionChainRunCount
		configs.Triggers.SmallDurationCalls.ConsiderCDRsFromLast = configsJSON.Triggers.SmallDurationCalls.ConsiderCDRsFromLast
		configs.Triggers.SmallDurationCalls.DurationThreshold, err = time.ParseDuration(configsJSON.Triggers.SmallDurationCalls.DurationThreshold)
		if err != nil {
			return nil, utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		//state.StateTriggers.StateSmallDurationCalls.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// ** Actions Section

	// * Email
	configs.Actions.Email.Enabled = configsJSON.Actions.Email.Enabled
	if configs.Actions.Email.Enabled {
		configs.Actions.Email.Username = configsJSON.Actions.Email.Username
		configs.Actions.Email.Password = configsJSON.Actions.Email.Password
		configs.Actions.Email.Message = configsJSON.Actions.Email.Message
	}

	// * HTTP
	configs.Actions.HTTP.Enabled = configsJSON.Actions.HTTP.Enabled

	// * Call
	configs.Actions.Call.Enabled = configsJSON.Actions.Call.Enabled

	// * LocalCommands
	configs.Actions.LocalCommands.Enabled = configsJSON.Actions.LocalCommands.Enabled

	// ** ActionChains Section
	configs.ActionChains.List = configsJSON.ActionChains.List

	// ** DataGroups Section
	configs.DataGroups.List = configsJSON.DataGroups.List

	logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Loaded Configs: %v", configs), false)

	return configs, nil

}

// FraudionConfig ...
type FraudionConfig struct {
	General      General
	CDRsSources  interface{}
	Triggers     Triggers
	Actions      Actions
	ActionChains ActionChains
	DataGroups   DataGroups
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

type triggerSimultaneousCalls struct {
	Enabled                  bool
	ExecuteInterval          time.Duration
	HitThreshold             uint32
	MinimumNumberLength      uint32
	ActionChainName          string
	ActionChainHoldoffPeriod uint32
	MaxActionChainRunCount   uint32
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
