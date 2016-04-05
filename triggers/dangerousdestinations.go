package triggers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"database/sql"
	//"os/exec"

	"github.com/SlyMarbo/gmail"
	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/utils"
)

// DangerousDestinationsRun ...
func DangerousDestinationsRun(startUpTime *time.Time, configs *config.FraudionConfig2, db *sql.DB) {

	fmt.Println("Starting Trigger, \"DangerousDestinations\"...")

	configsTrigger := configs.Triggers.DangerousDestinations

	//ticker := time.NewTicker(configsTrigger.CheckPeriod)
	//for executionTime := range ticker.C {

	//fmt.Println("DangerousDestinations Trigger executed at", executionTime)

	// TODO: db.Open() does not "open" a connection just returns a "pointer". db.Ping() is a way to see if the server is available. Should we try to find a way to close the connection between "ticks"?
	err := db.Ping()
	if err != nil {
		utils.DebugLogAndGetError(fmt.Sprintf("Something (%s) happened while checking the availability of the database connection to the CDRs", err.Error()), false)
		// TODO: Maybe also try to send an e-mail because this means the system is not being able to do it's job
	} else {

		// TODO: Consider separating hits that start with the prefix and that have the prefix inside (there could be matches inside that do not correspond with a call to the prefix)
		// TODO: Also consider having configurable trunk selection prefixes to which this will add the prefixes in the list and consider that as the dialled number (this is to detect trunk selection tries that could match)
		hits := make(map[string]uint32)
		//hitsIgnored := make(map[string]uint32)
		for _, prefix := range configsTrigger.PrefixList {
			hits[prefix] = uint32(0)
		}
		hitValues := []string{}
		//hitedValuesIgnored := []string{}

		// TODO: The interval should come from the configs? That's what the commented ConsiderCDRsFromLast field would be for
		var considerCDRsFromLast time.Duration
		var considerCDRsFromLastString = "200" // NOTE: Number of days if duration "rune" is missing (with rune would be something like "5d")

		numberOfDays, err := strconv.ParseInt(considerCDRsFromLastString, 0, 32)
		if err != nil { // Assume it's a Parseable Duration!
			considerCDRsFromLast, err = time.ParseDuration(considerCDRsFromLastString)
			if err != nil {
				utils.DebugLogAndGetError(fmt.Sprintf("Something (%s) happened while trying to parse \"considerCDRsFromLastString\"", err.Error()), false)
			}
		} else { // Assume it's a Number of Days
			considerCDRsFromLast, err = time.ParseDuration(fmt.Sprintf("%vh", numberOfDays*24))
			if err != nil {
				utils.DebugLogAndGetError(fmt.Sprintf("Something (%s) happened while trying to parse \"considerCDRsFromLastString\"", err.Error()), false)
			}
		}

		// NOTE: "guardDuration" and "guardTime" makes it so that when the service is restarted (maybe after an attack), the CDRs, from the start up time forward will only be considered from "startUpTime" - "guardTime" onwards, to try to prevent the system from reexecuting the Action Chain, that would execute for sure
		stringGuardDuration := "5000h" // TODO: This value should also come from the loaded configuration
		guardDuration, err := time.ParseDuration(stringGuardDuration)
		if err != nil {
			utils.DebugLogAndGetError(fmt.Sprintf("Something (%s) happened while trying to parse \"stringGuardTime\"", err.Error()), false)
		}
		guardTime := startUpTime.Add(-guardDuration)
		durationSinceGuardTime := time.Now().Sub(guardTime)

		// TODO: From here on what is done is Elastix2.3 specific, where the tests were made, so later we'll have to add some conditions to check what is the configured softswitch
		fmt.Println(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) AND calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(considerCDRsFromLast.Hours()), uint32(durationSinceGuardTime.Hours())))
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) AND calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(considerCDRsFromLast.Hours()), uint32(durationSinceGuardTime.Hours())))
		if err != nil {
			utils.DebugLogAndGetError(fmt.Sprintf("Something (%s) happened while trying to Query the CDRs database", err.Error()), false)
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

					// TODO: Should we match dials to more than one destination SIP/test/<number>&SIP/test/<number2>
					// TODO: Maybe the dial string match code should be from the interfaces because it's a softswitch specific thing
					matchesDialString := regexp.MustCompile("(?:SIP|DAHDI)/[^@&]+/([0-9]+)") // NOTE: Supported dial string format
					matchedString := matchesDialString.FindString(lastdata)
					if lastapp != "Dial" /*|| strings.Contains(lastapp, "Local") || !test */ || matchedString == "" { // NOTE: Ignore if "lastapp" is no Dial and "lastdata" does not contain an expected dial string
						continue
					}

					dialedNumber := matchesDialString.FindStringSubmatch(lastdata)[1]

					// TODO: Remove this print!
					//fmt.Println(dialedNumber)

					if uint32(len(dialedNumber)) >= configsTrigger.MinimumNumberLength {

						for _, prefix := range configsTrigger.PrefixList {

							// TODO: Maybe the "matchStringWithTag" should come from configs also? Also being able to add several?
							//matchStringWithTag := "([0-9]{0,8})?(0{2})?__prefix__[0-9]{5,}"
							matchStringWithTag := "([0-9]{0,8})?(0{2})?__prefix__[0-9]{5,}"
							matchString := strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)
							foundMatch, err := regexp.MatchString(matchString, lastdata)

							if err != nil {
								utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to match (found) a Prefix with regexp (%s)", err.Error()), false)
							}

							// TODO: Maybe the "matchStringWithTag" should come from configs also? Also being able to add several?
							matchStringWithTag = "[92][0-9]{8}" // NOTE: Ignores any number of 9 digits that starts with 9 or 2
							matchString = strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)
							foundIgnore, err := regexp.MatchString(matchString, lastdata)

							if err != nil {
								utils.DebugLogAndGetError(fmt.Sprintf("Something happened while trying to match (ignore) a Prefix with regexp (%s)", err.Error()), false)
							}

							if foundMatch == true && foundIgnore == false {
								hits[prefix] = hits[prefix] + 1
								hitValues = append(hitValues, dst)
							}

						}

					}

				}

			}

			runActionChain := false
			for _, hits := range hits {
				if hits >= configsTrigger.HitThreshold {
					runActionChain = true
				}
			}

			//actionChainGuardTime := configsTrigger.LastActionChainRunTime.Add(configs.General.DefaultActionChainSleepPeriod)

			if runActionChain /*&& actionChainGuardTime.Before(time.Now())*/ {

				actionChainName := configsTrigger.ActionChainName
				if actionChainName == "" {
					actionChainName = "Default"
				}

				fmt.Println("Running action chain:", actionChainName)
				configsTrigger.LastActionChainRunTime = time.Now()

				actionChain := configs.ActionChains.List[actionChainName]
				dataGroups := configs.DataGroups.List

				for _, v := range actionChain {

					// TODO: Remove this print!
					//fmt.Println(k, v)

					if v.ActionName == "*email" {

						// TODO: Should we assert here that Email Action is enabled or on config loading?

						email := gmail.Compose("Email subject", "Email body")
						email.From = configs.Actions.Email.Username
						email.Password = configs.Actions.Email.Password
						email.ContentType = "text/html; charset=utf-8"
						for _, dataGroupName := range v.DataGroupNames {
							email.AddRecipient(dataGroups[dataGroupName].EmailAddress)
						}

						err := email.Send()
						if err != nil {
							fmt.Println(err.Error())
						}

					} else if v.ActionName == "*localcommand" {

						// TODO: Should we assert here that the run user of the process has "root" permissions?
						/*for _, contact := range contacts {

							command := exec.Command(contact.CommandName, contact.CommandArguments)

							err := command.Run()
							if err != nil {
								fmt.Println(err.Error())
							}

						}*/

					}

				}

			}

		}

	}

	// TODO: Should we, and if so how, close the db connection between ticks? db.Close()

	//}

}
