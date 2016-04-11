package monitors

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"os/exec"

	"github.com/SlyMarbo/gmail"

	"github.com/andmar/fraudion/fraudion"
	"github.com/andmar/fraudion/logger"
)

// SimultaneousCallsRun ...
func SimultaneousCallsRun() {

	fraudion := fraudion.Global
	configs := fraudion.Configs
	state := fraudion.State.Triggers

	configsTrigger := configs.Triggers.SimultaneousCalls
	stateTrigger := state.StateDangerousDestinations

	logger.Log.Write(logger.ConstLoggerLevelInfo, "Starting Trigger, \"SimultaneousCalls\"...", false)

	ticker := time.NewTicker(configsTrigger.ExecuteInterval)

	for executionTime := range ticker.C {

		logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("SimultaneousCalls Trigger executed at %s", executionTime), false)

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
				logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Active Calls: %s", activeCalls), false)

				activeCallsInt, err := strconv.Atoi(activeCalls)
				if err == nil {

					runActionChain := false

					if uint32(activeCallsInt) > configsTrigger.HitThreshold {

						logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Active Calls greater than threshold (%v)\n", configsTrigger.HitThreshold), false)

						runActionChain = true

					}

					actionChainGuardTime := stateTrigger.LastActionChainRunTime.Add(configs.General.DefaultActionChainHoldoffPeriod)

					logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("stateTrigger.LastActionChainRunTime: %s", stateTrigger.LastActionChainRunTime), false)
					logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("configs.General.DefaultActionChainHoldoffPeriod: %s", configs.General.DefaultActionChainHoldoffPeriod), false)
					logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("actionChainGuardTime: %s", actionChainGuardTime), false)
					logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Now(): %s", time.Now()), false)
					logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("actionChainGuardTime < Now(): %v", actionChainGuardTime.Before(time.Now())), false)
					logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("stateTrigger.ActionChainRunCount > 0: %v", stateTrigger.ActionChainRunCount > 0), false)

					if runActionChain && actionChainGuardTime.Before(time.Now()) && stateTrigger.ActionChainRunCount > 0 {

						state.StateSimultaneousCalls.ActionChainRunCount = stateTrigger.ActionChainRunCount - 1

						actionChainName := configsTrigger.ActionChainName
						if actionChainName == "" {
							actionChainName = "*default"
						}

						logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Running action chain: %s", actionChainName), false)
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
