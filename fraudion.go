package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	//"time"

	"encoding/json"
	//"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
)

// Defines constants
const (
	CONFIG_DIR = "examples/config" // "/usr/share/fraudion"
	MYSQL      = "mysql"
	FREESWITCH = "fs"
	ASTERISK   = "ast_alone"
	ELASTIX    = "ast_el"
)

// Defines expected CLI arguments/flags
var (
	configDir = flag.String("config_dir", CONFIG_DIR, "<help message for 'config_dir'>")
)

// Config Data Holders
type Test struct {
	TestArrBuh map[string]string
}

type configJsonGeneral struct {
	Pbx_software                      string
	Cdrs_source                       string
	Default_trigger_check_period      string
	Default_action_chain_sleep_period string
	Default_action_chain_run_count    string
}

type configJsonTriggers struct {
	Simultaneous_calls_threshold map[string]interface{}
	Dangerous_destinations       map[string]interface{}
	Expected_destinations        map[string]interface{}
	Small_duration_calls         map[string]interface{}
}

type configActionLocalEmail struct {
	Message string
}
type configActionHTTPSMS struct {
	Url     string
	Method  string
	Message string
}
type configActionCall struct {
	Message string
}

type configActionHTTPPost struct {
	Url        string
	Parameters map[string]string
}

type configActionLocalCommands struct {
	list map[string]string
}

type configJsonActions struct {
	Email          configActionLocalEmail
	Http_sms       configActionHTTPSMS
	Call           configActionCall
	Http_post      configActionHTTPPost
	Local_commands map[string]string
}

type configActionChainAction struct {
	Action_name   string
	Contact_names []string
}

type configJsonActionChains map[string][]configActionChainAction

type configJsonContacts map[string]map[string]interface{}

/*struct configJsonTriggers {
	var Test string

	"simultaneous_calls_threshold": {},
	"dangerous_destinations": {},
	"expected_destinations": {}
}*/

// Starts here!
func main() {

	fmt.Println("Starting...")
	fmt.Println()

	flag.Parse()

	//if *version {
	//	fmt.Println("CGRateS " + utils.VERSION)
	//	return
	//}

	// Don't forget "pid" file

	// Parse config file
	/*fi, err := os.Stat(cfgDir)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			return cfg, nil
		}
		return nil, err
	} else if !fi.IsDir() && cfgDir != utils.CONFIG_DIR { // If config dir defined, needs to exist, not checking for default
		return nil, fmt.Errorf("Path: %s not a directory.", cfgDir)
	}
	if fi.IsDir() {

	err = filepath.Walk(cfgDir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					return nil
				}
				cfgFiles, err := filepath.Glob(filepath.Join(path, "*.json"))

	*/

	configFile, err := os.Open("examples/config/fraudion.json")
	if err != nil {
		// TODO Log this to Syslog
		fmt.Printf("There was an error opening the config file (\"fraudion.json\") at %s. :(\n", *configDir)
		os.Exit(-1)
	}

	var jsonConfig map[string]*json.RawMessage // Better than using the Lib example's Empty interface... https://tour.golang.org/methods/14
	jsonConfigReader := JsonConfigReader.New(configFile)

	err = json.NewDecoder(jsonConfigReader).Decode(&jsonConfig)
	if err != nil {
		// TODO Log this to Syslog
		fmt.Printf("There was an error parsing the config file (\"fraudion.json\") at %s. :(\n", *configDir)
		os.Exit(-1)
	}

	fmt.Println(jsonConfig)
	fmt.Println()

	rawCfg, hasKey := jsonConfig["general"]
	if !hasKey {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"general\" key in config Json. :(")
		os.Exit(-1)
	}
	cfgGeneral := new(configJsonGeneral)
	if err := json.Unmarshal(*rawCfg, cfgGeneral); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"general\" Json. :(")
		os.Exit(-1)
	}

	fmt.Println(reflect.TypeOf(cfgGeneral))
	fmt.Println(cfgGeneral)
	fmt.Println()

	rawCfg, hasKey = jsonConfig["triggers"]
	if !hasKey {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"triggers\" key in config Json. :(")
		os.Exit(-1)
	}
	cfgTriggers := new(configJsonTriggers)
	if err := json.Unmarshal(*rawCfg, cfgTriggers); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"triggers\" Json. :(")
		os.Exit(-1)
	}

	fmt.Println(reflect.TypeOf(cfgTriggers))
	fmt.Println(cfgTriggers)
	fmt.Println()

	rawCfg, hasKey = jsonConfig["actions"]
	if !hasKey {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"actions\" key in config Json. :(")
		os.Exit(-1)
	}
	cfgActions := new(configJsonActions)
	if err := json.Unmarshal(*rawCfg, cfgActions); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"actions\" Json. :(")
		os.Exit(-1)
	}

	fmt.Println(reflect.TypeOf(cfgActions))
	fmt.Println(cfgActions)
	fmt.Println()

	rawCfg, hasKey = jsonConfig["action_chains"]
	if !hasKey {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"action_chains\" key in config Json. :(")
		os.Exit(-1)
	}
	cfgActionChains := new(configJsonActionChains)
	if err := json.Unmarshal(*rawCfg, cfgActionChains); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"action_chains\" Json. :(")
		os.Exit(-1)
	}

	fmt.Println(reflect.TypeOf(cfgActionChains))
	fmt.Println(cfgActionChains)
	fmt.Println()

	rawCfg, hasKey = jsonConfig["contacts"]
	if !hasKey {
		// TODO Log this to Syslog
		fmt.Println("Could not find \"contacts\" key in config Json. :(")
		os.Exit(-1)
	}
	cfgContacts := new(configJsonContacts)
	if err := json.Unmarshal(*rawCfg, cfgContacts); err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not Unmarshal \"contacts\" Json. :(")
		os.Exit(-1)
	}

	fmt.Println(reflect.TypeOf(cfgContacts))
	fmt.Println(cfgContacts)
	fmt.Println()

	// TODO Parse configs for acceptable values.

	// Start Running!
	// ...

	os.Exit(1)

}
