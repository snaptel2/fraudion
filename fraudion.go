package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"

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
type configJsonGeneral struct {
	Monitored_software                string
	Cdrs_source                       string
	Default_trigger_check_period      string
	Default_action_chain_sleep_period string
	Default_action_chain_run_count    string
}

type configJsonTriggers struct {
	Simultaneous_calls     map[string]interface{}
	Dangerous_destinations map[string]interface{}
	Expected_destinations  map[string]interface{}
	Small_duration_calls   map[string]interface{}
}

type configJsonActionEmail struct {
	Default_message string
}

type configJsonActionCall struct {
	Default_message string
}

type configJsonActionHttp struct {
	Default_url        string
	Default_method     string
	Default_parameters map[string]string
}

type configActionLocalCommands struct {
	list map[string]string
}

type configJsonActions struct {
	Email          configJsonActionEmail
	Call           configJsonActionCall
	Http           configJsonActionHttp
	Local_commands map[string]string
}

type configActionChainAction struct {
	Action   string
	Contacts []string
	Command  string
}

type configJsonActionChains map[string][]configActionChainAction

type configJsonContacts map[string]map[string]interface{}

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
	simultaneous_calls_check_period, found := cfgTriggers.Simultaneous_calls["check_period"].(string)
	if found == false {
		// TODO Confirm that cfgGeneral.Default_trigger_check_period has some value, because this var needs some.
		simultaneous_calls_check_period = cfgGeneral.Default_trigger_check_period
	}

	dangerous_destinations_check_period, found := cfgTriggers.Dangerous_destinations["check_period"].(string)
	if found == false {
		// TODO Confirm that cfgGeneral.Default_trigger_check_period has some value, because this var needs some.
		dangerous_destinations_check_period = cfgGeneral.Default_trigger_check_period
	}

	fmt.Println("Debug:", simultaneous_calls_check_period)
	fmt.Println("Debug:", dangerous_destinations_check_period)

	// Start Running!
	simultaneous_calls_check_period_ticker_duration, err := time.ParseDuration(simultaneous_calls_check_period)
	if err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not calculate duration for \"simultaneous_calls_check_period_ticker\". :(")
		os.Exit(-1)
	}
	simultaneous_calls_check_period_ticker := time.NewTicker(simultaneous_calls_check_period_ticker_duration)
	go func() { // TODO Future simultaneous_calls_checker()
		for t := range simultaneous_calls_check_period_ticker.C {
			fmt.Println("simultaneous_calls_check_period_ticker tick at", t)
		}
	}()

	dangerous_destinations_check_period_ticker_duration, err := time.ParseDuration(dangerous_destinations_check_period)
	if err != nil {
		// TODO Log this to Syslog
		fmt.Println("Could not calculate duration for \"dangerous_destinations_check_period_ticker\". :(")
		os.Exit(-1)
	}
	dangerous_destinations_check_period_ticker := time.NewTicker(dangerous_destinations_check_period_ticker_duration)
	go func() { // TODO Future dangerous_destinations_checker()
		for t := range dangerous_destinations_check_period_ticker.C {
			fmt.Println("dangerous_destinations_check_period_ticker tick at", t)
		}
	}()

	for {
	}

	//os.Exit(1)

}
