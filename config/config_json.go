package config

import (
	"fmt"
	"os"
	"reflect"

	"encoding/json"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
)

const (
	DEFAULT_JSON_CONFIG_FILENAME = "fraudion.json"
)

// FraudionJSONConfig ...
type FraudionJSONConfig struct {
	General      JSONGeneral
	Triggers     JSONTriggers
	Actions      JSONActions
	ActionChains JSONActionChains
	Contacts     JSONContacts
}

// JSONGeneral ...
type JSONGeneral struct {
	Monitored_software                        string
	Cdrs_source                               string
	Default_trigger_check_period              string
	Default_action_chain_sleep_period         string
	Default_action_chain_run_count            uint32
	Default_minimum_destination_number_length uint32
}

// JSONTriggers ...
type JSONTriggers struct {
	Simultaneous_calls     map[string]interface{}
	Dangerous_destinations map[string]interface{}
	Expected_destinations  map[string]interface{}
	Small_duration_calls   map[string]interface{}
}

// JSONActions ...
type JSONActions struct {
	Email          jsonActionEmail
	Call           jsonActionCall
	HTTP           jsonActionHTTP
	Local_commands map[string]string
}

type jsonActionEmail struct {
	Default_message string
}

type jsonActionCall struct {
	Default_message string
}

type jsonActionHTTP struct {
	Default_url        string
	Default_method     string
	Default_parameters map[string]string
}

// JSONActionChains ...
type JSONActionChains map[string][]jsonActionChainAction

type jsonActionChainAction struct {
	Action   string
	Contacts []string
	Command  string
}

// JSONContacts ...
type JSONContacts map[string]map[string]interface{}

// LoadConfigFromJSONFile ...
func (fraudionJSONConfig *FraudionJSONConfig) LoadConfigFromJSONFile(configDir string) error {

	configFileName := DEFAULT_JSON_CONFIG_FILENAME
	configFileFullPath := filepath.Join(configDir, DEFAULT_JSON_CONFIG_FILENAME)

	fmt.Printf("Trying to open the JSON config file \"%s\" at \"%s\" (fullpath: \"%s\").\n", configFileName, configDir, configFileFullPath)
	configFileJSON, err := os.Open(filepath.Join(configDir, configFileName))
	if err != nil {
		// TODO Log this to Syslog
		fmt.Printf("ERROR: There was an error (\"%s\") opening the JSON config file \"%s\" at \"%s\" (fullpath: \"%s\"). :(\n", err.Error(), configFileName, configDir, configFileFullPath)
		os.Exit(-1)
	}
	defer configFileJSON.Close()

	/*scanner := bufio.NewScanner(configFileJSON)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}*/

	var RawJSON map[string]*json.RawMessage // Better than using the Lib example's Empty interface... https://tour.golang.org/methods/14
	JSONConfigReader := JsonConfigReader.New(configFileJSON)

	err = json.NewDecoder(JSONConfigReader).Decode(&RawJSON)
	if err != nil {
		// TODO Log this to Syslog
		fmt.Printf("ERROR: There was an error (\"%s\") parsing the JSON config file. :(\n", err.Error())
		return err
	}

	//fmt.Println(RawJSON)
	//fmt.Println()

	// General Section
	rawCfg, hasKey := RawJSON["general"]
	if hasKey == false {
		// TODO Log this to Syslog
		fmt.Println("ERROR: Could not find \"general\" section in config JSON. :(")
		return fmt.Errorf("could not find \"general\" section in config json")
	}
	cfgJSONGeneral := new(JSONGeneral)
	if err := json.Unmarshal(*rawCfg, cfgJSONGeneral); err != nil {
		// TODO Log this to Syslog
		fmt.Println("ERROR: Could not Unmarshal \"general\" JSON. :(")
		return fmt.Errorf("could not Unmarshal \"general\" JSON")
	}

	fmt.Println(reflect.TypeOf(cfgJSONGeneral))
	fmt.Println(cfgJSONGeneral)

	fraudionJSONConfig.General = *cfgJSONGeneral

	// Triggers Section
	rawCfg, hasKey = RawJSON["triggers"]
	if hasKey == false {
		// TODO Log this to Syslog
		fmt.Println("ERROR: Could not find \"triggers\" section in config JSON. :(")
		return fmt.Errorf("could not find \"triggers\" section in config JSON")
	}
	cfgJSONTriggers := new(JSONTriggers)
	if err := json.Unmarshal(*rawCfg, cfgJSONTriggers); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"triggers\" JSON. :(")
		return fmt.Errorf("could not Unmarshal \"triggers\" JSON")
	}

	fmt.Println(reflect.TypeOf(cfgJSONTriggers))
	fmt.Println(cfgJSONTriggers)

	fraudionJSONConfig.Triggers = *cfgJSONTriggers

	// Actions Section
	rawCfg, hasKey = RawJSON["actions"]
	if hasKey == false {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"actions\" section in config JSON. :(")
		return fmt.Errorf("could not find \"actions\" section in config JSON")
	}
	cfgJSONActions := new(JSONActions)
	if err := json.Unmarshal(*rawCfg, cfgJSONActions); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"actions\" JSON. :(")
		return fmt.Errorf("could not Unmarshal \"actions\" JSON")
	}

	fmt.Println(reflect.TypeOf(cfgJSONActions))
	fmt.Println(cfgJSONActions)

	fraudionJSONConfig.Actions = *cfgJSONActions

	// Action Chains Section
	rawCfg, hasKey = RawJSON["action_chains"]
	if hasKey == false {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"action_chains\" section in config JSON. :(")
		return fmt.Errorf("could not find \"action_chains\" section in config JSON")
	}
	cfgJSONActionChains := new(JSONActionChains)
	if err := json.Unmarshal(*rawCfg, cfgJSONActionChains); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"action_chains\" JSON. :(")
		return fmt.Errorf("could not Unmarshal \"action_chains\" JSON")
	}

	fmt.Println(reflect.TypeOf(cfgJSONActionChains))
	fmt.Println(cfgJSONActionChains)

	fraudionJSONConfig.ActionChains = *cfgJSONActionChains

	// Contacts Section
	rawCfg, hasKey = RawJSON["contacts"]
	if !hasKey {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"contacts\" section in config JSON. :(")
		return fmt.Errorf("could not find \"contacts\" section in config JSON")
	}
	cfgJSONContacts := new(JSONContacts)
	if err := json.Unmarshal(*rawCfg, cfgJSONContacts); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"contacts\" JSON. :(")
		return fmt.Errorf("could not Unmarshal \"contacts\" JSON")
	}

	fmt.Println(reflect.TypeOf(cfgJSONContacts))
	fmt.Println(cfgJSONContacts)

	fraudionJSONConfig.Contacts = *cfgJSONContacts

	return nil

}
