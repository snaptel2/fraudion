package triggers

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/fraudion/config"
)

// ExpectedDestinationsRun ...
func ExpectedDestinationsRun(configs *config.FraudionConfig, db *sql.DB) {

	fmt.Println("Starting Trigger, \"ExpectedDestinations\"")

	triggerConfigs := configs.Triggers.ExpectedDestinations

	ticker := time.NewTicker(triggerConfigs.CheckPeriod)

	for t := range ticker.C {

		fmt.Println("expectedDestinations tick at", t)

	}

}
