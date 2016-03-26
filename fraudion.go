package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fraudion/config"
)

// Defines constants
const (
	constDefaultConfigDir = "examples/config" // TODO: This will probably need to be changed to something like "/usr/share/fraudion" when we aproach a more usable version!
	/*MYSQL              = "mysql"
	FREESWITCH         = "fs"
	ASTERISK           = "ast_alone"
	ELASTIX            = "ast_el"*/
)

// Defines expected CLI arguments/flags
var (
	cliConfigDir = flag.String("configdir", constDefaultConfigDir, "<help message for 'config_dir'>")
)

// Starts here!
func main() {

	fmt.Println()
	fmt.Println("Starting...")

	fmt.Println("Parsing CLI parameters...")

	ConfigsJSON := new(config.FraudionConfigJSON)
	err := ConfigsJSON.LoadConfigFromJSONFile(*cliConfigDir)
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("** Parsed JSON:", ConfigsJSON)

	configs := new(config.FraudionConfig)
	err = configs.CheckJSONSanityAndLoadConfigs(ConfigsJSON)
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("** Configs:", configs)

	//fmt.Println()
	//fmt.Println(config)

	/*fmt.Println("sadjaçidjfsdf", cfgGeneral.Monitored_software)
	fmt.Println("sadjaçidjfsdf", supported_software)
	found := stringInSlice(cfgGeneral.Monitored_software, supported_software)
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
	}*/

	//os.Exit(1)

}
