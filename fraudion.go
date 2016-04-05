package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	//"database/sql"

	"github.com/andmar/fraudion/config"
	//"github.com/andmar/fraudion/interfaces/softswitches"
	//"github.com/andmar/fraudion/triggers"

	_ "github.com/go-sql-driver/mysql"
)

// Defines constants
const (
	constDefaultConfigDir = "examples/config" // TODO: This will probably need to be changed to something like "/usr/share/fraudion" when we aproach a more usable version!
)

// Defines expected CLI arguments/flags
var (
	argCliConfigDir = flag.String("configdir", constDefaultConfigDir, "<help message for 'configdir'>")
	argCliTest      = flag.String("test", "", "<help message for 'test'>")
	argCliDBPass    = flag.String("dbpass", "", "<help message for 'dbpass'>")
	startUpTime     time.Time
)

// Starts here!
func main() {

	fmt.Println("Starting...")
	//startUpTime := time.Now()

	fmt.Println("Parsing CLI parameters...")
	flag.Parse()

	ConfigsJSON := new(config.FraudionConfigJSON2)
	err := ConfigsJSON.LoadConfigFromJSONFile2(*argCliConfigDir)
	if err != nil {
		fmt.Printf("There was an error (%s) parsing the Fraudion configuration file\n", err)
		os.Exit(-1)
	}

	fmt.Println("** Parsed JSON:", ConfigsJSON)
	fmt.Println()

	os.Exit(-1)

	/*

		configs := new(config.FraudionConfig)
		err = configs.CheckJSONSanityAndLoadConfigs(ConfigsJSON)
		if err != nil {
			fmt.Printf("There was an error (%s) loading the Fraudion configuration\n", err)
			os.Exit(-1)
		}

		fmt.Println("** Loaded Configs:", configs)
		fmt.Println()

		var dbstring string
		// TODO: This database connection is Elastix2.3 specific, where the tests were made, so later we'll have to add some conditions to check what is the configured softswitch
		if *argCliDBPass == "" {
			dbstring = fmt.Sprintf("root:@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1")
		} else {
			dbstring = fmt.Sprintf("root:%s@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1", *argCliDBPass)
		}

		db, err := sql.Open("mysql", dbstring)
		if err != nil {
			errorMessage := "There was an error (%s) trying to open a connection to the database\n"
			fmt.Printf(errorMessage)
			os.Exit(-1)
		}

		// Launch Triggers!
		if configs.Triggers.SimultaneousCalls.Enabled == true {
			go triggers.SimultaneousCallsRun(configs, new(softswitches.Asterisk1_8))
		}
		if configs.Triggers.DangerousDestinations.Enabled == true {
			go triggers.DangerousDestinationsRun(startUpTime, configs, db)
		}

		if configs.Triggers.ExpectedDestinations.Enabled == true {
			go triggers.ExpectedDestinationsRun(configs, db)
		}

		if configs.Triggers.SmallDurationCalls.Enabled == true {
			go triggers.SmallDurationCallsRun(configs, db)
		}

	*/

	// Sleep!
	for {

		// Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
