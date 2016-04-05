package config

import (
	"fmt"
	"os"

	"encoding/json"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
	"github.com/andmar/fraudion/utils"
)

const (
	constDefaultJSONConfigFilename2 = "fraudion_2.json"
)

// LoadConfigFromJSONFile2 ...
func (fraudionJSONConfig *FraudionConfigJSON2) LoadConfigFromJSONFile2(configDir string) error {

	configFileName := constDefaultJSONConfigFilename2

	fmt.Printf("Loading JSON from config file \"%s\"...\n", filepath.Join(configDir, configFileName))

	// ** JSON config file to map[string] to Raw JSON
	JSONconfigFile, err := os.Open(filepath.Join(configDir, configFileName))
	if err != nil {
		customErrorMessage := fmt.Sprintf("There was an error opening the JSON config file (\"%s\")", err.Error())
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}
	defer JSONconfigFile.Close()

	var RawJSON map[string]*json.RawMessage // NOTE: Better than using the Lib example's Empty interface... https://tour.golang.org/methods/14
	JSONConfigFileReader := JsonConfigReader.New(JSONconfigFile)

	err = json.NewDecoder(JSONConfigFileReader).Decode(&RawJSON) // NOTE: Reads the JSON file to JSONConfigReader as a map[string]<Raw JSON that has to be decoded further!>
	if err != nil {
		customErrorMessage := fmt.Sprintf("There was an error doing the initial parsing of the JSON config file (\"%s\")", err.Error())
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// ** General Section
	sectionName := "general"
	rawGeneralJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configGeneralJSON := new(GeneralJSON2)
	if err := json.Unmarshal(*rawGeneralJSON, configGeneralJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	fraudionJSONConfig.General = *configGeneralJSON
	fmt.Println("General:", configGeneralJSON)

	// ** Triggers Section
	sectionName = "triggers"
	rawTriggersJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configTriggersJSON := new(TriggersJSON2)
	if err := json.Unmarshal(*rawTriggersJSON, configTriggersJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	fraudionJSONConfig.Triggers = *configTriggersJSON
	fmt.Println("Triggers:", configTriggersJSON)

	// ** Actions Section
	sectionName = "actions"
	rawActionsJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configActionsJSON := new(ActionsJSON2)
	if err := json.Unmarshal(*rawActionsJSON, configActionsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	fraudionJSONConfig.Actions = *configActionsJSON
	fmt.Println("Actions:", configActionsJSON)

	// ** Actions Chains Section
	sectionName = "action_chains"
	rawActionChainsJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configActionChainsJSON := new(ActionChainsJSON2)
	if err := json.Unmarshal(*rawActionChainsJSON, configActionChainsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	fraudionJSONConfig.ActionChains = *configActionChainsJSON
	fmt.Println("Action Chains:", configActionChainsJSON)

	// ** Data Groups Section
	sectionName = "data_groups"
	rawDataGroupsJSON, hasKey := RawJSON[sectionName]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	configDataGroupsJSON := new(DataGroupsJSON2)
	if err := json.Unmarshal(*rawDataGroupsJSON, configDataGroupsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"%s\" section in config JSON", sectionName)
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	fraudionJSONConfig.DataGroups = *configDataGroupsJSON
	fmt.Println("Data Groups:", configDataGroupsJSON)

	fmt.Println("All Aboard:", fraudionJSONConfig)

	return nil

}
