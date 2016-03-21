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
	supported_softwares = []string{"*ast_el_2.3", "*ast_alone_1.8"}
	configDir           = flag.String("config_dir", CONFIG_DIR, "<help message for 'config_dir'>")
)

// Config Data Holders
type configJsonGeneral struct {
	Monitored_software                string
	Cdrs_source                       string
	Default_trigger_check_period      string
	Default_action_chain_sleep_period string
	Default_action_chain_run_count    string
}

type ConfigGeneral struct {
	Monitored_software                string
	Cdrs_source                       string
	Default_trigger_check_period      time.Duration
	Default_action_chain_sleep_period time.Duration
	Default_action_chain_run_count    uint32
}

type configJsonTriggers struct {
	Simultaneous_calls     map[string]interface{}
	Dangerous_destinations map[string]interface{}
	Expected_destinations  map[string]interface{}
	Small_duration_calls   map[string]interface{}
}

type triggerSimultaneousCalls struct {
	Check_period       time.Duration
	Max_call_threshold uint32
}

type triggerDangerousDestinations struct {
	Check_period  time.Duration
	Prefix_list   []string
	Hit_threshold uint32
}

type triggerExpectedDestinations struct {
	Check_period  time.Duration
	Prefix_list   []string
	Hit_threshold uint32
}

type triggerSmallCallDurations struct {
	Check_period       time.Duration
	Duration_threshold time.Duration
	Hit_threshold      uint32
}

type configTriggers struct {
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

type configJsonActions struct {
	Email          configJsonActionEmail
	Call           configJsonActionCall
	Http           configJsonActionHttp
	Local_commands map[string]string
}

type configActionEmail struct {
	Default_message string
}

type configActionCall struct {
	Default_message string
}

type configActionHttp struct {
	Default_url        string
	Default_method     string
	Default_parameters map[string]string
}

type configActionLocalCommands struct {
	List map[string]string
}

type configActions struct {
	Email          configActionEmail
	Call           configActionCall
	Http           configActionHttp
	Local_commands map[string]string
}

type configActionChainAction struct {
	Action   string
	Contacts []string
	Command  string
}

type configActionChain struct {
	Chains map[string][]configActionChainAction
}

type configJsonActionChainAction struct {
	Action   string
	Contacts []string
	Command  string
}

type configJsonActionChains map[string][]configJsonActionChainAction

type configJsonContacts map[string]map[string]interface{}

type configContact struct {
	For_actions     []string
	Phone_number    string
	Email           string
	Message         string
	Http_url        string
	Http_method     string
	Http_parameters map[string]string
}

type configContacts struct {
	List map[string]configContact
}

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
	defer configFile.Close()

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
	// General
	/*
	   "monitored_software": "*ast_el_2.3.0", // ast_alone_<version>: Asterisk from <version> source, ast_el_<version>: asterisk from Elastix <version>, fs_<version>: Freeswitch <version>
	   "cdrs_source": "*db_mysql", // db_mysql: Database MySQL. TODO (v2) json: Json, csv: CSVs
	   "default_trigger_check_period": "10s", // Can be overriden in Trigger configuration
	   "default_action_chain_sleep_period": "30m", // Waits this time before checking if Chain run account is still > 0 the Action Chain for a match of some Trigger
	   "default_action_chain_run_count": "3" // Count for the Reprocessing of the Action Chains
	*/
	configGeneral := new(ConfigGeneral)

	fmt.Println("sadjaçidjfsdf", cfgGeneral.Monitored_software)
	fmt.Println("sadjaçidjfsdf", supported_softwares)
	found := stringInSlice(cfgGeneral.Monitored_software, supported_softwares)
	if found == false {
		// TODO Log this to Syslog
		fmt.Println("Could not find configured \"monitored_software\" in supported software list. It could be the syntax or you're system may not be supported. :(")
		os.Exit(-1)
	}
	configGeneral.Monitored_software = cfgGeneral.Monitored_software

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
		// Main "thread" has to Sleep or else 100% CPU, as expected!
		time.Sleep(100000 * time.Hour)
	}

	//os.Exit(1)

}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
