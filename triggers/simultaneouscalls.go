package triggers

import (
	"fmt"
	"time"

	"github.com/fraudion/config"
	"github.com/fraudion/interfaces/softswitches"
)

// SimultaneousCallsRun ...
func SimultaneousCallsRun(configs *config.FraudionConfig, softswitch softswitches.Softswitch) {

	fmt.Println("Starting Trigger, \"SimultaneousCalls\"")

	triggerConfigs := configs.Triggers.SimultaneousCalls

	ticker := time.NewTicker(triggerConfigs.CheckPeriod)

	for t := range ticker.C {

		fmt.Println("simultaneousCalls tick at", t)

	}

}
