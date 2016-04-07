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
	argCliLogFile   = flag.String("logto", constDefaultLogFile, "<help message for 'log'>")
	argCliConfigDir = flag.String("configin", constDefaultConfigDir, "<help message for 'configdir'>")
	argCliDBPass    = flag.String("dbpass", "", "<help message for 'dbpass'>")
)

func init() {
	types.Fraudion = new(types.FraudionGlobal)
	types.Fraudion.SetupConfigs()
	types.Fraudion.SetupState()
	types.Fraudion.Debug = true
}

// Starts here!
func main() {

	fraudion := types.Fraudion
	configs := fraudion.Configs

	fraudion.StartUpTime = time.Now()
	os.Stdout.WriteString(fmt.Sprintf("Starting at %s\n", fraudion.StartUpTime))
	os.Stdout.WriteString("Parsing CLI flags...\n")
	flag.Parse()

	/*
		var logFile *os.File
		if _, err := os.Stat(*argCliLogFile); os.IsNotExist(err) {
			logFile, err = os.Create(*argCliLogFile)
			if err != nil {
				os.Stdout.WriteString(fmt.Sprintf("Can't start, there was a problem (%s) creating the Log file. :(\n", err))
				os.Exit(1)
			}
		} else {
			logFile, err = os.Open(*argCliLogFile)
			if err != nil {
				os.Stdout.WriteString(fmt.Sprintf("Can't start, there was a problem ()%s) opening the Log file. :(\n", err))
				os.Exit(1)
			}
		}*/
	//fraudion.SetupLogging(logFile, logFile, logFile, logFile)
	types.Fraudion.SetupLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	fraudion.LogInfo.Printf("Starting at %s\n", fraudion.StartUpTime)

	configsJSON := new(types.FraudionConfigJSON)
	err := config.ParseConfigsFromJSON(configsJSON, *argCliConfigDir)
	if err != nil {
		fraudion.LogError.Fatalf("There was an error (%s) parsing the Fraudion JSON configuration file\n", err)
	}

	err = config.ValidateAndLoadConfigs(configsJSON, false)
	if err != nil {
		fraudion.LogError.Fatalf("There was an error (%s) validating/loading the Fraudion configuration\n", err)
	}

	var db *sql.DB
	if configs.Triggers.DangerousDestinations.Enabled == true || configs.Triggers.ExpectedDestinations.Enabled == true || configs.Triggers.SmallDurationCalls.Enabled == true {
		fraudion.LogInfo.Println("Connecting to the CDRs Database...")
		var dbstring string
		if *argCliDBPass == "" {
			dbstring = fmt.Sprintf("root:@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1")
		} else {
			dbstring = fmt.Sprintf("root:%s@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1", *argCliDBPass)
		}
		db, err = sql.Open("mysql", dbstring)
		if err != nil {
			fraudion.LogError.Fatalf("There was an error (%s) trying to open a connection to the database\n", err)
		}
	}

	// Launch Triggers!
	fraudion.LogInfo.Println("Launching enabled triggers...")
	/*if configs.Triggers.SimultaneousCalls.Enabled == true {
		go triggers.SimultaneousCallsRun(configs, new(softswitches.Asterisk1_8))
	}*/
	if configs.Triggers.DangerousDestinations.Enabled == true {
		go triggers.DangerousDestinationsRun(db)
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
