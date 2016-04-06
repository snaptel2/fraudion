package config

/*
import (
	"fmt"
	"os"
	//"reflect"

	"encoding/json"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
	"github.com/andmar/fraudion/types"
	"github.com/andmar/fraudion/utils"
)

const (
	constDefaultJSONConfigFilename = "fraudion.json"
)

// LoadConfigFromJSONFile ...
func LoadConfigFromJSONFile(configDir string) error {

	// TODO: Remove this print!
	//fmt.Println("** JSON file parsing start.")

	configFileName := constDefaultJSONConfigFilename
	//configFileFullPath := filepath.Join(configDir, constDefaultJSONConfigFilename)

	// ** JSON config file to map[string] to Raw JSON
	//fmt.Printf("Trying to open the JSON config file \"%s\" at \"%s\" (fullpath: \"%s\").\n", configFileName, configDir, configFileFullPath)
	configFileJSON, err := os.Open(filepath.Join(configDir, configFileName))
	if err != nil {
		customErrorMessage := fmt.Sprintf("There was an error opening the JSON config file (\"%s\")", err.Error())
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}
	defer configFileJSON.Close()

	var RawJSON map[string]*json.RawMessage // NOTE: Better than using the Lib example's Empty interface... https://tour.golang.org/methods/14
	JSONConfigReader := JsonConfigReader.New(configFileJSON)

	err = json.NewDecoder(JSONConfigReader).Decode(&RawJSON) // NOTE: Reads the JSON file to JSONConfigReader as a map[string]<Raw JSON that has to be decoded further!>
	if err != nil {
		customErrorMessage := fmt.Sprintf("There was an error doing the initial parsing of the JSON config file (\"%s\")", err.Error())
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// ** JSON config file main configuration sections to "objects" (general, triggers, actions, etc)
	// General Section
	rawGeneralJSON, hasKey := RawJSON["General"]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"General\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}
	configGeneralJSON := new(types.GeneralJSON)
	if err := json.Unmarshal(*rawGeneralJSON, configGeneralJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"General\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// TODO: Remove these prints!
	//fmt.Print(reflect.TypeOf(configGeneralJSON))
	//fmt.Println(" ", configGeneralJSON)

	fraudionJSONConfig.General = *types.configGeneralJSON

	// Triggers Sectionwefd
	rawTriggersJSON, hasKey := RawJSON["Triggers"]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"Triggers\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}
	configTriggersJSON := new(TriggersJSON)
	if err := json.Unmarshal(*rawTriggersJSON, configTriggersJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"Triggers\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// TODO: Remove these prints!
	//fmt.Print(reflect.TypeOf(configTriggersJSON))
	//fmt.Println(" ", configTriggersJSON)

	fraudionJSONConfig.Triggers = *configTriggersJSON

	// Actions Section
	rawActionsJSON, hasKey := RawJSON["Actions"]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"Actions\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}
	configActionsJSON := new(ActionsJSON)
	if err := json.Unmarshal(*rawActionsJSON, configActionsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"Actions\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// TODO: Remove these prints!
	//fmt.Print(reflect.TypeOf(configActionsJSON))
	//fmt.Println(" ", configActionsJSON)

	fraudionJSONConfig.Actions = *configActionsJSON

	// Action Chains Section
	rawActionChainsJSON, hasKey := RawJSON["ActionChains"]
	if hasKey == false {
		customErrorMessage := fmt.Sprintf("Could not find \"ActionChains\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}
	configActionChainsJSON := new(ActionChainsJSON)
	if err := json.Unmarshal(*rawActionChainsJSON, configActionChainsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"ActionChains\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// TODO: Remove these prints!
	//fmt.Print(reflect.TypeOf(configActionChainsJSON))
	//fmt.Println(" ", configActionChainsJSON)

	fraudionJSONConfig.ActionChains = *configActionChainsJSON

	// Contacts Section
	rawContactsJSON, hasKey := RawJSON["Contacts"]
	if !hasKey {
		customErrorMessage := fmt.Sprintf("Could not find \"Contacts\" section in config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}
	configContactsJSON := new(ContactsJSON)
	if err := json.Unmarshal(*rawContactsJSON, configContactsJSON); err != nil {
		customErrorMessage := fmt.Sprintf("Could not Unmarshal \"Contacts\" section inf config JSON")
		return utils.DebugLogAndGetError(customErrorMessage, true)
	}

	// TODO: Remove these prints!
	//fmt.Print(reflect.TypeOf(configContactsJSON))
	//fmt.Println(" ", configContactsJSON)

	fraudionJSONConfig.Contacts = *configContactsJSON

	// TODO: Remove this print!
	//fmt.Println("** JSON file parsing end.")

	return nil

}*/
