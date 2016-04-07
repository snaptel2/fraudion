package config

import (
	"fmt"
	"time"

	"github.com/andmar/fraudion/types"
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

// ValidateAndLoadConfigs ...
func ValidateAndLoadConfigs(configsJSON *types.FraudionConfigJSON, validateOnly bool) error {

	fraudion := types.Fraudion
	configs := fraudion.Configs

	fmt.Println(configs)

	types.Fraudion.LogInfo.Println("Validating and Loading configurations...")

	// ** General Section

	// * MonitoredSoftware
	if configsJSON.General.MonitoredSoftware == "" {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"monitored_software\" value in section \"general\" missing OR is empty."), true)
	}
	found := utils.StringInStringsSlice(configsJSON.General.MonitoredSoftware, constSupportedSoftware2)
	if !found {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"monitored_software\" value in section \"general\" must be one of %s", constSupportedSoftware2), true)
	}
	if !validateOnly {
		configs.General.MonitoredSoftware = configsJSON.General.MonitoredSoftware
	}

	// * CDRsSource
	if configsJSON.General.CDRsSource == "" {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"cdrs_source\" value in section \"general\" missing OR is empty."), true)
	}
	found = utils.StringInStringsSlice(configsJSON.General.CDRsSource, constSupportedCDRSources2)
	if !found {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"cdrs_source\" value in section \"general\" must be one of %s", constSupportedCDRSources2), true)
	}
	if !validateOnly {
		configs.General.CDRsSource = configsJSON.General.CDRsSource
	}

	// * DefaultTriggerExecuteInterval
	if configsJSON.General.DefaultTriggerExecuteInterval == "" {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_trigger_execute_interval\" value in section \"general\" missing OR is empty."), true)
	}
	durationDefaultTriggerExecuteInterval, err := time.ParseDuration(configsJSON.General.DefaultTriggerExecuteInterval)
	if err != nil {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_trigger_execute_interval\" value in section \"general\" is not a valid duration. Must be a parseable duration, a number followed by one of: \"s\", \"m\" or \"h\" for \"seconds\", \"minutes\" and \"hours\" respectively"), true)
	}
	durationConstMinimumTriggerExecuteInterval, err := time.ParseDuration(constMinimumTriggerExecuteInterval)
	if err != nil {
		return utils.DebugLogAndGetError("(Internal) There seems to be an issue with the definition of constMinimumTriggerExecuteInterval2", true)
	}
	if durationDefaultTriggerExecuteInterval < durationConstMinimumTriggerExecuteInterval {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_trigger_execute_interval\" value in section \"general\" is too small. Value must be > %s", constMinimumTriggerExecuteInterval), true)
	}
	if !validateOnly {
		configs.General.DefaultTriggerExecuteInterval = durationDefaultTriggerExecuteInterval
	}

	// * DefaultHitThreshold
	if configsJSON.General.DefaultHitThreshold == 0 {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_hit_threshold\" value in section \"general\" missing OR is 0."), true)
	}
	if !validateOnly {
		configs.General.DefaultHitThreshold = configsJSON.General.DefaultHitThreshold
	}

	// * DefaultMinimumDestinationNumberLength
	if configsJSON.General.DefaultMinimumDestinationNumberLength == 0 {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_minimum_destination_number_length\" value in section \"general\" missing OR is 0."), true)
	}
	if !validateOnly {
		configs.General.DefaultMinimumDestinationNumberLength = configsJSON.General.DefaultMinimumDestinationNumberLength
	}

	// * DefaultActionChainHoldoffPeriod
	if configsJSON.General.DefaultActionChainHoldoffPeriod == "" {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_holdoff_period\" value in section \"general\" missing OR empty."), true)
	}
	durationDefaultActionChainHoldoffPeriod, err := time.ParseDuration(configsJSON.General.DefaultActionChainHoldoffPeriod)
	if err != nil {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_holdoff_period\" value in section \"general\" is not a valid duration. Must be a parseable duration, a number followed by one of: \"s\", \"m\" or \"h\" for \"seconds\", \"minutes\" and \"hours\" respectively"), true)
	}
	durationConstMaximumActionChainHoldoffPeriod, err := time.ParseDuration(constMaximumActionChainHoldoffPeriod)
	if err != nil {
		return utils.DebugLogAndGetError("(Internal) There seems to be an issue with the definition of constMaximumActionChainHoldoffPeriod", true)
	}
	if durationDefaultActionChainHoldoffPeriod > durationConstMaximumActionChainHoldoffPeriod {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_holdoff_period\" value in section \"general\" is too small. Value must be > %s", constMaximumActionChainHoldoffPeriod), true)
	}
	if !validateOnly {
		configs.General.DefaultActionChainHoldoffPeriod = durationDefaultActionChainHoldoffPeriod
	}

	// * DefaultActionChainRunCount
	if configsJSON.General.DefaultActionChainRunCount == 0 {
		return utils.DebugLogAndGetError(fmt.Sprintf("\"default_action_chain_run_count\" value in section \"general\" missing OR 0."), true)
	}
	if !validateOnly {
		configs.General.DefaultActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// TODO: From this point the validateOnly flag is not yet checked and used to do something

	// ** Triggers Section

	// * SimultaneousCalls
	configs.Triggers.SimultaneousCalls.Enabled = configsJSON.Triggers.SimultaneousCalls.Enabled
	if configs.Triggers.SimultaneousCalls.Enabled {
		configs.Triggers.SimultaneousCalls.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.SimultaneousCalls.ExecuteInterval)
		if err != nil {
			return utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.SimultaneousCalls.HitThreshold = configsJSON.Triggers.SimultaneousCalls.HitThreshold
		configs.Triggers.SimultaneousCalls.MinimumNumberLength = configsJSON.Triggers.SimultaneousCalls.MinimumNumberLength
		configs.Triggers.SimultaneousCalls.ActionChainName = configsJSON.Triggers.SimultaneousCalls.ActionChainName
		//configs.Triggers.SimultaneousCalls.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// * DangerousDestinations
	configs.Triggers.DangerousDestinations.Enabled = configsJSON.Triggers.DangerousDestinations.Enabled
	if configs.Triggers.DangerousDestinations.Enabled {
		configs.Triggers.DangerousDestinations.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.DangerousDestinations.ExecuteInterval)
		if err != nil {
			return utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.DangerousDestinations.HitThreshold = configsJSON.Triggers.DangerousDestinations.HitThreshold
		configs.Triggers.DangerousDestinations.MinimumNumberLength = configsJSON.Triggers.DangerousDestinations.MinimumNumberLength
		configs.Triggers.DangerousDestinations.ActionChainName = configsJSON.Triggers.DangerousDestinations.ActionChainName
		configs.Triggers.DangerousDestinations.PrefixList = configsJSON.Triggers.DangerousDestinations.PrefixList
		configs.Triggers.DangerousDestinations.ConsiderCDRsFromLast = configsJSON.Triggers.DangerousDestinations.ConsiderCDRsFromLast
		configs.Triggers.DangerousDestinations.MatchRegex = configsJSON.Triggers.DangerousDestinations.MatchRegex
		configs.Triggers.DangerousDestinations.IgnoreRegex = configsJSON.Triggers.DangerousDestinations.IgnoreRegex
		//configs.Triggers.DangerousDestinations.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// * ExpectedDestinations
	configs.Triggers.ExpectedDestinations.Enabled = configsJSON.Triggers.ExpectedDestinations.Enabled
	if configs.Triggers.ExpectedDestinations.Enabled {
		configs.Triggers.ExpectedDestinations.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.ExpectedDestinations.ExecuteInterval)
		if err != nil {
			return utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.ExpectedDestinations.HitThreshold = configsJSON.Triggers.ExpectedDestinations.HitThreshold
		configs.Triggers.ExpectedDestinations.MinimumNumberLength = configsJSON.Triggers.ExpectedDestinations.MinimumNumberLength
		configs.Triggers.ExpectedDestinations.ActionChainName = configsJSON.Triggers.ExpectedDestinations.ActionChainName
		configs.Triggers.ExpectedDestinations.PrefixList = configsJSON.Triggers.ExpectedDestinations.PrefixList
		configs.Triggers.ExpectedDestinations.ConsiderCDRsFromLast = configsJSON.Triggers.ExpectedDestinations.ConsiderCDRsFromLast
		configs.Triggers.ExpectedDestinations.MatchRegex = configsJSON.Triggers.ExpectedDestinations.MatchRegex
		configs.Triggers.ExpectedDestinations.IgnoreRegex = configsJSON.Triggers.ExpectedDestinations.IgnoreRegex
		//configs.Triggers.ExpectedDestinations.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
	}

	// * SmallDurationCalls
	configs.Triggers.SmallDurationCalls.Enabled = configsJSON.Triggers.ExpectedDestinations.Enabled
	if configs.Triggers.SmallDurationCalls.Enabled {
		configs.Triggers.SmallDurationCalls.ExecuteInterval, err = time.ParseDuration(configsJSON.Triggers.ExpectedDestinations.ExecuteInterval)
		if err != nil {
			return utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		configs.Triggers.SmallDurationCalls.HitThreshold = configsJSON.Triggers.SmallDurationCalls.HitThreshold
		configs.Triggers.SmallDurationCalls.MinimumNumberLength = configsJSON.Triggers.SmallDurationCalls.MinimumNumberLength
		configs.Triggers.SmallDurationCalls.ActionChainName = configsJSON.Triggers.SmallDurationCalls.ActionChainName
		configs.Triggers.SmallDurationCalls.ConsiderCDRsFromLast = configsJSON.Triggers.SmallDurationCalls.ConsiderCDRsFromLast
		configs.Triggers.SmallDurationCalls.DurationThreshold, err = time.ParseDuration(configsJSON.Triggers.SmallDurationCalls.DurationThreshold)
		if err != nil {
			return utils.DebugLogAndGetError(fmt.Sprintf("Error message missing..."), true)
		}
		//configs.Triggers.SmallDurationCalls.ActionChainRunCount = configsJSON.General.DefaultActionChainRunCount
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

	types.Fraudion.LogInfo.Printf("Loaded Configs: %v", configs)

	return nil

}
