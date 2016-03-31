package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"database/sql"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/interfaces/softswitches"
	"github.com/andmar/fraudion/triggers"

	_ "github.com/go-sql-driver/mysql"
)

// Defines constants
const (
	constDefaultConfigDir = "examples/config" // TODO: This will probably need to be changed to something like "/usr/share/fraudion" when we aproach a more usable version!
)

// Defines expected CLI arguments/flags
var (
	cliConfigDir = flag.String("configdir", constDefaultConfigDir, "<help message for 'configdir'>")
	cliTest      = flag.String("test", "", "<help message for 'test'>")
	cliDBPass    = flag.String("dbpass", "", "<help message for 'dbpass'>")
)

// Starts here!
func main() {

	fmt.Println("Starting...")

	fmt.Println("Parsing CLI parameters...")
	flag.Parse()

	ConfigsJSON := new(config.FraudionConfigJSON)
	err := ConfigsJSON.LoadConfigFromJSONFile(*cliConfigDir)
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("** Parsed JSON:", ConfigsJSON)
	fmt.Println()

	configs := new(config.FraudionConfig)
	err = configs.CheckJSONSanityAndLoadConfigs(ConfigsJSON)
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("** Loaded Configs:", configs)
	fmt.Println()

	var dbstring string
	if *cliDBPass == "" {
		dbstring = fmt.Sprintf("root:@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1")
	} else {
		dbstring = fmt.Sprintf("root:%s@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1", *cliDBPass)
	}

	//dbstring := fmt.Sprintf("root:%s@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1", dbpass)
	db, err := sql.Open("mysql", dbstring)
	//db, err := sql.Open("mysql", "root:@/test")
	fmt.Println(db)
	fmt.Println(err)

	err = db.Ping() // Open doest not open a connection. This is the way to see if the server is available.
	fmt.Println(err)

	// Launch Triggers!
	if configs.Triggers.SimultaneousCalls.Enabled == true {
		go triggers.SimultaneousCallsRun(configs, new(softswitches.Asterisk1_8))
	}
	if configs.Triggers.DangerousDestinations.Enabled == true {
		go triggers.DangerousDestinationsRun(configs, db)
	}

	if configs.Triggers.ExpectedDestinations.Enabled == true {
		go triggers.ExpectedDestinationsRun(configs, db)
	}

	if configs.Triggers.SmallDurationCalls.Enabled == true {
		go triggers.SmallDurationCallsRun(configs, db)
	}

	// Sleep!
	for {

		// Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
