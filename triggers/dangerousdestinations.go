package triggers

import (
	"fmt"
	//"log"
	"regexp"
	"strconv"
	"time"

	"database/sql"
	//"net/smtp"
	//"os/exec"

	//"github.com/SlyMarbo/gmail"
	"github.com/fraudion/config"
	"github.com/fraudion/utils"
)

// DangerousDestinationsRun ...
func DangerousDestinationsRun(configs *config.FraudionConfig, db *sql.DB) {

	fmt.Println("Starting Trigger, \"DangerousDestinations\"")

	triggerConfigs := configs.Triggers.DangerousDestinations

	ticker := time.NewTicker(triggerConfigs.CheckPeriod)
	for t := range ticker.C {

		fmt.Println("DangerousDestinations executed at", t)

		err := db.Ping() // Open does not "open" a connection. This is the way to see if the server is available.
		if err != nil {
			utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to \"Open\" a connection to the CDRs database (%s)", err.Error()), false)
		} else {

			// TODO: Consider separating hits that start with the prefix and that have the prefix inside (there could be matches inside that do not correspond with a call to the prefix)
			// TODO: Also consider having configurable trunk selection prefixes to which this will add the prefixes in the list and consider that as the dialled number (this is to detect trunk selection tries that could match)
			hits := make(map[string]uint32)
			for _, prefix := range triggerConfigs.PrefixList {
				hits[prefix] = uint32(0)
			}

			// TODO: The interval should come from the configs? That's what the commented ConsiderCDRsFromLast field would be for
			var ConsiderCDRsFromLast time.Duration
			var ConsiderCDRsFromLastString = "30"

			numberOfDays, err2 := strconv.ParseInt(ConsiderCDRsFromLastString, 0, 32)
			if err2 != nil {
				// Assume it's a Parseable Duration!
				ConsiderCDRsFromLast, _ = time.ParseDuration(ConsiderCDRsFromLastString)
			} else {
				// Assume it's a Number of Days
				ConsiderCDRsFromLast, _ = time.ParseDuration(fmt.Sprintf("%vh", numberOfDays*24))
			}
			fmt.Println(uint32(ConsiderCDRsFromLast.Hours()))

			rows, err := db.Query(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(ConsiderCDRsFromLast.Hours())))
			if err != nil {
				utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to Query the CDRs database (%s)", err.Error()), false)
			} else {

				for rows.Next() {

					var calldate string
					var clid string
					var src string
					var dst string
					var dcontext string
					var channel string
					var dstchannel string
					var lastapp string
					var lastdata string
					var duration uint32
					var billsec uint32
					var disposition string
					var amaflags uint32
					var accountcode string
					var uniqueid string
					var userfield string

					err := rows.Scan(&calldate,
						&clid,
						&src,
						&dst,
						&dcontext,
						&channel,
						&dstchannel,
						&lastapp,
						&lastdata,
						&duration,
						&billsec,
						&disposition,
						&amaflags,
						&accountcode,
						&uniqueid,
						&userfield)

					/*
						fmt.Println(calldate,
							clid,
							src,
							dst,
							//dcontext,
							//channel,
							//dstchannel,
							//lastapp,
							//lastdata,
							duration,
							billsec,
							disposition,
							//amaflags,
							//accountcode,
							//uniqueid,
							//userfield

						)
					*/

					if err != nil {
						utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to get the CDR data (%s)", err.Error()), false)
					} else {

						if uint32(len(dst)) >= triggerConfigs.MinimumNumberLength {

							for _, prefix := range triggerConfigs.PrefixList {

								found, err := regexp.MatchString(fmt.Sprintf("00%s", prefix), dst)
								if err != nil {
									utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to match a Prefix with regexp (%s)", err.Error()), false)
								}

								if found == true {
									hits[prefix] = hits[prefix] + 1
								}

							}

						}

						for _, hits := range hits {

							if hits >= triggerConfigs.HitThreshold {

								addition := triggerConfigs.LastActionChainRunTime.Add(configs.General.DefaultActionChainSleepPeriod)

								if addition.Before(time.Now()) {

									fmt.Println("WOW! This would run action chain:", triggerConfigs.ActionChainName)
									triggerConfigs.LastActionChainRunTime = time.Now()

									// TODO Run actionChain!

									// E-mail tests:
									/*email := gmail.Compose("Email subject", "Email body")
									email.From = ""     // The same as username! Don't commit this!
									email.Password = "" // Don't commit this!

									// Defaults to "text/plain; charset=utf-8" if unset.
									email.ContentType = "text/html; charset=utf-8"

									// Normally you'll only need one of these, but I thought I'd show both.
									email.AddRecipient("am@voipit.pt")
									//email.AddRecipients("another@example.com", "more@example.com")

									err := email.Send()
									if err != nil {
										fmt.Println(err.Error())
									}*/

									// System commands!
									//err := exec.Command("touch", "buh.txt").Run()
									//fmt.Println(err)

								}

							}

						}

					}

				}

			}

			fmt.Println(hits)

			defer db.Close()

		}

	}

}
