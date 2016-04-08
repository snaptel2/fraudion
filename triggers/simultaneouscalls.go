package triggers

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"os/exec"

	"github.com/SlyMarbo/gmail"

	"github.com/andmar/fraudion/types"
)

// SimultaneousCallsRun ...
func SimultaneousCallsRun() {

	fraudion := types.Fraudion
	configs := fraudion.Configs
	state := fraudion.State

	configsTrigger := configs.Triggers.SimultaneousCalls
	stateTrigger := state.StateTriggers.StateDangerousDestinations

	fraudion.LogInfo.Println("Starting Trigger, \"SimultaneousCalls\"...")

	ticker := time.NewTicker(configsTrigger.ExecuteInterval)

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
				fraudion.LogInfo.Println("Active Calls:", activeCalls)

				activeCallsInt, err := strconv.Atoi(activeCalls)
				if err == nil {

					runActionChain := false

					if uint32(activeCallsInt) > configsTrigger.HitThreshold {

						fraudion.LogInfo.Printf("Active Calls greater than threshold (%v)\n", configsTrigger.HitThreshold)

						runActionChain = true

					}

					actionChainGuardTime := stateTrigger.LastActionChainRunTime.Add(configs.General.DefaultActionChainHoldoffPeriod)

					fraudion.LogInfo.Println("stateTrigger.LastActionChainRunTime:", stateTrigger.LastActionChainRunTime)
					fraudion.LogInfo.Println("configs.General.DefaultActionChainHoldoffPeriod:", configs.General.DefaultActionChainHoldoffPeriod)
					fraudion.LogInfo.Println("actionChainGuardTime:", actionChainGuardTime)
					fraudion.LogInfo.Println("Now():", time.Now())
					fraudion.LogInfo.Println("actionChainGuardTime < Now():", actionChainGuardTime.Before(time.Now()))
					fraudion.LogInfo.Println("stateTrigger.ActionChainRunCount > 0:", stateTrigger.ActionChainRunCount > 0)

					if runActionChain && actionChainGuardTime.Before(time.Now()) && stateTrigger.ActionChainRunCount > 0 {

						state.StateTriggers.StateSimultaneousCalls.ActionChainRunCount = stateTrigger.ActionChainRunCount - 1

						actionChainName := configsTrigger.ActionChainName
						if actionChainName == "" {
							actionChainName = "*default"
						}

						types.Fraudion.LogInfo.Println("Running action chain: ", actionChainName)
						stateTrigger.LastActionChainRunTime = time.Now()

						actionChain := configs.ActionChains.List[actionChainName]
						dataGroups := configs.DataGroups.List

						for _, v := range actionChain {

							if v.ActionName == "*email" {

								// TODO: Should we assert here that Email Action is enabled here or on config validation?

								body := fmt.Sprintf("Active Calls:\n\n%v", activeCallsInt)

								email := gmail.Compose("Fraudion ALERT: Simultaneous Calls!", fmt.Sprintf("\n\n%s", body))
								email.From = configs.Actions.Email.Username
								email.Password = configs.Actions.Email.Password
								fmt.Println(configs.Actions.Email.Username, configs.Actions.Email.Password)
								email.ContentType = "text/html; charset=utf-8"
								for _, dataGroupName := range v.DataGroupNames {
									fmt.Println(dataGroups[dataGroupName].EmailAddress)
									email.AddRecipient(dataGroups[dataGroupName].EmailAddress)
								}

								err := email.Send()
								if err != nil {
									fmt.Println(err.Error())
								}

							} else if v.ActionName == "*localcommand" {

								// TODO: Should we assert here that the run user of the process has "root" permissions?

								for _, dataGroupName := range v.DataGroupNames {

									command := exec.Command(dataGroups[dataGroupName].CommandName, dataGroups[dataGroupName].CommandArguments)

									err := command.Run()
									if err != nil {
										fmt.Println(err.Error())
									}

								}

							} else {

								fmt.Println("Unsupported Action!")

							}

						}

					}

				} else {
					fmt.Println(err.Error())
				}

			}

		}
		if err := in.Err(); err != nil {
			fmt.Println(err.Error())
		}

	}

}
