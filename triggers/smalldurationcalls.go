package triggers

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/andmar/fraudion/types"
)

// SmallDurationCallsRun ...
func SmallDurationCallsRun(configs *types.FraudionConfig, db *sql.DB) {

	fmt.Println("Starting Trigger, \"SmallDurationCalls\"")

	triggerConfigs := configs.Triggers.SmallDurationCalls

	ticker := time.NewTicker(triggerConfigs.ExecuteInterval)

	for t := range ticker.C {

		fmt.Println("smallDurationCalls tick at", t)

	}

}
