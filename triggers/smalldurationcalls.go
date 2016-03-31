package triggers

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/andmar/fraudion/config"
)

// SmallDurationCallsRun ...
func SmallDurationCallsRun(configs *config.FraudionConfig, db *sql.DB) {

	fmt.Println("Starting Trigger, \"SmallDurationCalls\"")

	triggerConfigs := configs.Triggers.SmallDurationCalls

	ticker := time.NewTicker(triggerConfigs.CheckPeriod)

	for t := range ticker.C {

		fmt.Println("smallDurationCalls tick at", t)

	}

}
