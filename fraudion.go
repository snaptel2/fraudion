package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"database/sql"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/triggers"
	"github.com/andmar/fraudion/types"

	_ "github.com/go-sql-driver/mysql"
)

// Defines Constants
const (
	constDefaultConfigDir = "/etc/fraudion"
	constDefaultLogFile   = "/var/log/fraudion.log"
)

// Defines expected CLI flags
var (
	argCliLogFile   = flag.String("log", constDefaultLogFile, "<help message for 'log'>")
	argCliConfigDir = flag.String("configdir", constDefaultConfigDir, "<help message for 'configdir'>")
	argCliDBPass    = flag.String("dbpass", "", "<help message for 'dbpass'>")
)

func init() {
	types.Globals = new(types.Fraudion)
	types.Globals.Debug = true
	types.Globals.SetupLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

// Starts here!
func main() {

	types.Globals.StartUpTime = time.Now()
	types.Globals.LogInfo.Printf("Starting at %s\n", types.Globals.StartUpTime)

	types.Globals.LogInfo.Println("Parsing CLI flags...")
	flag.Parse()

	configsJSON := new(types.FraudionConfigJSON2)
	err := config.LoadConfigFromJSONFile(configsJSON, *argCliConfigDir)
	if err != nil {
		types.Globals.LogError.Fatalf("There was an error (%s) parsing the Fraudion JSON configuration file\n", err)
	}

	configs := new(types.FraudionConfig2)
	err = config.ValidateAndLoadConfigs(configs, configsJSON, false)
	if err != nil {
		types.Globals.LogError.Fatalf("There was an error (%s) validating/loading the Fraudion configuration\n", err)
	}

	os.Exit(-1)

	types.Globals.LogInfo.Println("Connecting to the CDRs Database...")
	// TODO: This is here only for testing purposes, maybe this will move to the Triggers code, but the information must be global, maybe come from config file
	// TODO: This database connection method is Elastix2.3 specific, where the tests were made, so later we'll have to add some conditions to check what is the configured softswitch
	var dbstring string
	if *argCliDBPass == "" {
		dbstring = fmt.Sprintf("root:@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1")
	} else {
		dbstring = fmt.Sprintf("root:%s@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1", *argCliDBPass)
	}
	db, err := sql.Open("mysql", dbstring)
	if err != nil {
		types.Globals.LogError.Fatalf("There was an error (%s) trying to open a connection to the database\n", err)
	}

	// Launch Triggers!
	types.Globals.LogInfo.Println("Launching enabled triggers...")
	/*if configs.Triggers.SimultaneousCalls.Enabled == true {
		go triggers.SimultaneousCallsRun(configs, new(softswitches.Asterisk1_8))
	}*/
	if configs.Triggers.DangerousDestinations.Enabled == true {
		go triggers.DangerousDestinationsRun(configs, db)
	}

	/*if configs.Triggers.ExpectedDestinations.Enabled == true {
		go triggers.ExpectedDestinationsRun(configs, db)
	}

	if configs.Triggers.SmallDurationCalls.Enabled == true {
		go triggers.SmallDurationCallsRun(configs, db)
	}*/

	// Sleep!
	for {

		// Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
