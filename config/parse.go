package config

import (
	"fmt"
	"os"

	"encoding/json"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"

	"github.com/andmar/fraudion/logger"
	"github.com/andmar/fraudion/utils"
)

const (
	constDefaultJSONConfigFilename = "fraudion.json"
)

// Parse ...
func Parse(configDir string) (*FraudionConfigJSON, error) {

	configsJSON := new(FraudionConfigJSON)

	configFileName := constDefaultJSONConfigFilename

	logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Parsing JSON from config file \"%s\"...", filepath.Join(configDir, configFileName)), false)

	// ** JSON config file to map[string] to Raw JSON
	JSONconfigFile, err := os.Open(filepath.Join(configDir, configFileName))
	if err != nil {
		customErrorMessage := fmt.Sprintf("There was an error (%s) opening the JSON config file", err.Error())
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}
	defer JSONconfigFile.Close()

	var RawJSON map[string]*json.RawMessage // NOTE: Better than using the Lib example's Empty interface... https://tour.golang.org/methods/14
	JSONConfigFileReader := JsonConfigReader.New(JSONconfigFile)

	err = json.NewDecoder(JSONConfigFileReader).Decode(&RawJSON) // NOTE: Reads the JSON file to JSONConfigReader as a map[string]<Raw JSON that has to be decoded further!>
	if err != nil {
		customErrorMessage := fmt.Sprintf("There was an error (%s) doing the initial parsing of the JSON config file", err.Error())
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// ** General Section
	sectionName := "general"
	rawGeneralJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configGeneralJSON := new(GeneralJSON)
	if err := json.Unmarshal(*rawGeneralJSON, configGeneralJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not (%s) Unmarshal \"%s\" section in config JSON", err, sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configsJSON.General = *configGeneralJSON
	//Fraudion.LogInfo.Println("General:", configGeneralJSON)

	// ** CDRsSources Section
	sectionName = "cdrs_sources"
	rawCDRsSourcesJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configCDRsSourcesJSON := new(map[string]map[string]string)
	if err := json.Unmarshal(*rawCDRsSourcesJSON, configCDRsSourcesJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not (%s) Unmarshal \"%s\" section in config JSON", err, sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configsJSON.CDRsSources = *configCDRsSourcesJSON
	//Fraudion.LogInfo.Println("CDRsSources:", *configCDRsSourcesJSON)

	// ** Triggers Section
	sectionName = "triggers"
	rawTriggersJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configTriggersJSON := new(TriggersJSON)
	if err := json.Unmarshal(*rawTriggersJSON, configTriggersJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configsJSON.Triggers = *configTriggersJSON
	//Fraudion.LogInfo.Println("Triggers:", configTriggersJSON)

	// ** Actions Section
	sectionName = "actions"
	rawActionsJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configActionsJSON := new(ActionsJSON)
	if err := json.Unmarshal(*rawActionsJSON, configActionsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configsJSON.Actions = *configActionsJSON
	//Fraudion.LogInfo.Println("Actions:", configActionsJSON)

	// ** Actions Chains Section
	sectionName = "action_chains"
	rawActionChainsJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configActionChainsJSON := new(ActionChainsJSON)
	if err := json.Unmarshal(*rawActionChainsJSON, configActionChainsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configsJSON.ActionChains = *configActionChainsJSON
	//Fraudion.LogInfo.Println("Action Chains:", configActionChainsJSON)

	// ** Data Groups Section
	sectionName = "data_groups"
	rawDataGroupsJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configDataGroupsJSON := new(DataGroupsJSON)
	if err := json.Unmarshal(*rawDataGroupsJSON, configDataGroupsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return nil, utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configsJSON.DataGroups = *configDataGroupsJSON
	//ypes.Fraudion.LogInfo.Println("Data Groups:", configDataGroupsJSON)

	logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Parsed Configs: %v", configsJSON), false)

	return configsJSON, nil

}

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
