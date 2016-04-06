package triggers

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/andmar/fraudion/types"
)

// ExpectedDestinationsRun ...
func ExpectedDestinationsRun(configs *types.FraudionConfig2, db *sql.DB) {

	fmt.Println("Starting Trigger, \"ExpectedDestinations\"")

	triggerConfigs := configs.Triggers.ExpectedDestinations

	ticker := time.NewTicker(triggerConfigs.ExecuteInterval)

	for t := range ticker.C {

		fmt.Println("expectedDestinations tick at", t)

	}

}
