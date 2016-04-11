package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	//"database/sql"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/fraudion"
	"github.com/andmar/fraudion/logger"
	//"github.com/andmar/fraudion/monitors"

	_ "github.com/go-sql-driver/mysql"
)

// Defines Constants
const (
	constDefaultConfigDir = "/etc/fraudion"
	constDefaultLogFile   = "/var/log/fraudion.log" // TODO: The system now defaults to STDOUT so this will be removed soon
)

// Defines expected CLI flags
var (
	argCliLogFile   = flag.String("logto", "", "<help message for 'logto'>") // NOTE: The default is "" because we use this to detect if the user has specifiec any file, if not, the system defaults to using STDOUT automatically.
	argCliConfigDir = flag.String("configin", constDefaultConfigDir, "<help message for 'configin'>")
	argCliDBPass    = flag.String("dbpass", "", "<help message for 'dbpass'>")
)

// Starts here!
func main() {

	fraudion := fraudion.Global // NOTE: fraudion.Global (and it's pointers) is (are) initialized on fraudion's package init() function

	fraudion.StartUpTime = time.Now()
	os.Stdout.WriteString(fmt.Sprintf("Starting Fraudion at %s\n", fraudion.StartUpTime))
	os.Stdout.WriteString("Parsing CLI flags...\n")
	flag.Parse()

	if strings.ToLower(*argCliLogFile) != "" {

		var logFile *os.File
		logFile, err := os.OpenFile(*argCliLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			os.Stdout.WriteString(fmt.Sprintf("Can't start, there was a problem (%s) opening the Log file. :(\n", err))
			os.Exit(1)
		}

		os.Stdout.WriteString(fmt.Sprintf("Outputting Log to \"%s\"\n", *argCliLogFile))
		logger.Log.SetHandles(logFile, logFile, logFile, logFile) // NOTE: Overwrite the default handles on the Logger object.
		logFile.WriteString("\n")

	}

	logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Starting Fraudion Log at %s\n", fraudion.StartUpTime), false)

	configsJSON, err := config.Parse(*argCliConfigDir)
	if err != nil {
		logger.Log.Write(logger.ConstLoggerLevelError, fmt.Sprintf("There was an error (%s) parsing the Fraudion JSON configuration file\n", err), true)
	}

	configs, err := config.Load(configsJSON)
	if err != nil {
		logger.Log.Write(logger.ConstLoggerLevelError, fmt.Sprintf("There was an error (%s) validating/loading the Fraudion configuration\n", err), true)
	}

	// TODO: We'll call config.Validate() here in the future

	fmt.Println(configs)

	// TODO: This will maybe be done elsewhere!
	/*var db *sql.DB
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
	if configs.Triggers.SimultaneousCalls.Enabled == true {
		go monitors.SimultaneousCallsRun()
	}
	if configs.Triggers.DangerousDestinations.Enabled == true {
		go monitors.DangerousDestinationsRun(db)
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
