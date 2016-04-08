package triggers

import (
	"bufio"
	"fmt"
	"regexp"
	"time"

	"os/exec"

	"github.com/andmar/fraudion/types"
)

// SimultaneousCallsRun ...
func SimultaneousCallsRun() {

	fraudion := types.Fraudion
	configs := fraudion.Configs
	//state := fraudion.State.StateTriggers

	triggerConfigs := configs.Triggers.SimultaneousCalls
	//stateTrigger := state.StateDangerousDestinations

	fraudion.LogInfo.Println("Starting Trigger, \"SimultaneousCalls\"...")

	ticker := time.NewTicker(triggerConfigs.ExecuteInterval)

	for executionTime := range ticker.C {

		fraudion.LogInfo.Println("SimultaneousCalls Trigger executed at", executionTime)

		//command := exec.Command("asterisk", "-rx 'core show channels'")
		command := exec.Command("asterisk", "-rx", "core show channels")
		stdout, err := command.StdoutPipe()
		if err != nil {
			fmt.Println(err.Error())
		}

		if err := command.Start(); err != nil {
			fmt.Println(err.Error())
		}

		in := bufio.NewScanner(stdout)

		for in.Scan() {

			searchActiveCallsNumber := regexp.MustCompile("^([0-9]+) active calls?$") // NOTE: Supported dial string format
			isActiveCallsLine := searchActiveCallsNumber.MatchString(in.Text())

			if isActiveCallsLine {

				activeCalls := searchActiveCallsNumber.FindStringSubmatch(in.Text())[1]
				fmt.Println("Active Calls:", activeCalls)

			}

		}
		if err := in.Err(); err != nil {
			fmt.Println(err.Error())
		}

	}

}
