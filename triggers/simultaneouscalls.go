package triggers

import (
	"fmt"
	"time"

	"github.com/andmar/fraudion/interfaces/softswitches"
	"github.com/andmar/fraudion/types"
)

// SimultaneousCallsRun ...
func SimultaneousCallsRun(configs *types.FraudionConfig2, softswitch softswitches.Softswitch) {

	fmt.Println("Starting Trigger, \"SimultaneousCalls\"")

	triggerConfigs := configs.Triggers.SimultaneousCalls

	ticker := time.NewTicker(triggerConfigs.ExecuteInterval)

	for t := range ticker.C {

		fmt.Println("simultaneousCalls tick at", t)

	}

}
