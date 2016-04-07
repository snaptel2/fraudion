package triggers

import (
	"bufio"
	"fmt"
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
		command := exec.Command("git", "status")
		stdout, err := command.StdoutPipe()
		if err != nil {
			fmt.Println(err.Error())
		}

		if err := command.Start(); err != nil {
			fmt.Println(err.Error())
		}

		// read command's stdout line by line
		in := bufio.NewScanner(stdout)

		for in.Scan() {
			fmt.Println(in.Text()) // write each line to your log, or anything you need
		}
		if err := in.Err(); err != nil {
		}

		/*err := command.Run()
		if err != nil {
			fmt.Println(err.Error())
		}*/

	}

}
