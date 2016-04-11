package monitors

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/andmar/fraudion/fraudion"
	"github.com/andmar/fraudion/logger"
)

// SmallDurationCallsRun ...
func SmallDurationCallsRun(db *sql.DB) {

	fraudion := fraudion.Global
	configs := fraudion.Configs
	//state := fraudion.State.Triggers

	configsTrigger := configs.Triggers.DangerousDestinations
	//stateTrigger := state.StateDangerousDestinations

	logger.Log.Write(logger.ConstLoggerLevelInfo, "Starting Trigger, \"SmallDurationCalls\"...", false)

	ticker := time.NewTicker(configsTrigger.ExecuteInterval)

	for executionTime := range ticker.C {

		logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("DangerousDestinations Trigger executed at %s", executionTime), false)

	}

}
