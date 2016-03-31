package triggers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"database/sql"
	//"os/exec"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/utils"
	//"github.com/SlyMarbo/gmail"
)

// DangerousDestinationsRun ...
func DangerousDestinationsRun(configs *config.FraudionConfig, db *sql.DB) {

	fmt.Println("Starting Trigger, \"DangerousDestinations\"...")

	triggerConfigs := configs.Triggers.DangerousDestinations

	ticker := time.NewTicker(triggerConfigs.CheckPeriod)
	for executionTime := range ticker.C {

		fmt.Println("DangerousDestinations Trigger executed at", executionTime)

		// TODO: db.Open() does not "open" a connection just returns a "pointer". db.Ping() is a way to see if the server is available. Should we try to find a way to close the connection between "ticks"?
		err := db.Ping()
		if err != nil {
			utils.DebugLogAndGetError(fmt.Sprintf("Something happened while checking the availability of the database connection to the CDRs (%s)", err.Error()), false)
		} else {

			// TODO: Consider separating hits that start with the prefix and that have the prefix inside (there could be matches inside that do not correspond with a call to the prefix)
			// TODO: Also consider having configurable trunk selection prefixes to which this will add the prefixes in the list and consider that as the dialled number (this is to detect trunk selection tries that could match)
			hits := make(map[string]uint32)
			//hitsIgnored := make(map[string]uint32)
			for _, prefix := range triggerConfigs.PrefixList {
				hits[prefix] = uint32(0)
			}
			hitValues := []string{}
			//hitedValuesIgnored := []string{}

			// TODO: The interval should come from the configs? That's what the commented ConsiderCDRsFromLast field would be for
			var ConsiderCDRsFromLast time.Duration
			var ConsiderCDRsFromLastString = "30" // NOTE: Number of days if duration rune is missing (with rune would be something like "5d")

			numberOfDays, err2 := strconv.ParseInt(ConsiderCDRsFromLastString, 0, 32)
			if err2 != nil { // Assume it's a Parseable Duration!
				ConsiderCDRsFromLast, _ = time.ParseDuration(ConsiderCDRsFromLastString)
			} else { // Assume it's a Number of Days
				ConsiderCDRsFromLast, _ = time.ParseDuration(fmt.Sprintf("%vh", numberOfDays*24))
			}

			// TODO: Remove this prints!
			//fmt.Println(uint32(ConsiderCDRsFromLast.Hours()))

			// TODO: Maybe we should consider that we only check from the program startup time or that time - a couple of hours/days, this is because if we have a detected attack we can restart the service without having it fire actions, or else having a way of saying from what time we should start considering and we can reset that time in case of an attack
			// TODO: From here on what is done is Elastix2.3 specific, where the tests were made, so we'll have to add some conditions to check configured softswitch
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

					// TODO: Remove this prints!
					/*fmt.Println(calldate,
						clid,
						src,
						dst,
						//dcontext,
						//channel,
						//dstchannel,
						lastapp,
						lastdata,
						duration,
						billsec,
						disposition,
						//amaflags,
						//accountcode,
						//uniqueid,
						//userfield

					)*/

					if err != nil {
						utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to get the CDR data (%s)", err.Error()), false)
					} else {

						if uint32(len(dst)) >= triggerConfigs.MinimumNumberLength {

							for _, prefix := range triggerConfigs.PrefixList {

								// TODO: Do we have to assert "DefaultMinimumDestinationNumberLength" of the number part in last data?

								// TODO: Maybe the "matchStringWithTag" should come from configs also? Also being able to add several?
								matchStringWithTag := "/([0-9]{0,8})?(0{2})?__prefix__[0-9]{5,}"
								matchString := strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)
								foundMatch, err := regexp.MatchString(matchString, lastdata)

								// TODO: Remove this prints!
								fmt.Println("Found:", matchString)

								if err != nil {
									utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to match (found) a Prefix with regexp (%s)", err.Error()), false)
								}

								// TODO: Maybe the "matchStringWithTag" should come from configs also? Also being able to add several?
								matchStringWithTag = "/[0-9]{9}" // NOTE: Ignores any number of 9 digits afther the dial string /
								matchString = strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)
								foundIgnore, err := regexp.MatchString(matchString, lastdata)

								// TODO: Remove this prints!
								fmt.Println("Ignore:", matchString)

								if err != nil {
									utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to match (ignore) a Prefix with regexp (%s)", err.Error()), false)
								}

								if foundMatch == true && foundIgnore == false && lastapp == "Dial" {
									hits[prefix] = hits[prefix] + 1
									hitValues = append(hitValues, dst)
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

			// TODO: Remove this prints!
			fmt.Println(hits)
			fmt.Println(hitValues)

			// TODO: Should we, and if so how, close the db connection between ticks? db.Close()

		}

	}

}
