package triggers

import (
	"fmt"
	"regexp"
	"time"

	"database/sql"

	"github.com/fraudion/config"
)

// DangerousDestinationsRun ...
func DangerousDestinationsRun(configs *config.FraudionConfig, db *sql.DB) {

	fmt.Println("Starting Trigger, \"DangerousDestinations\"")

	triggerConfigs := configs.Triggers.DangerousDestinations

	ticker := time.NewTicker(triggerConfigs.CheckPeriod)

	err := db.Ping() // Open doest not open a connection. This is the way to see if the server is available.
	fmt.Println(err)

	hits := make(map[string]uint32)
	for _, prefix := range triggerConfigs.PrefixList {
		hits[prefix] = uint32(0)
	}

	rows, err := db.Query("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL 2 MONTH) ORDER BY calldate DESC;")
	if err != nil {
		fmt.Println(err)
	}

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

		if err != nil {
			fmt.Println(err)
		}

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

		if uint32(len(dst)) >= triggerConfigs.MinimumNumberLength {

			for _, prefix := range triggerConfigs.PrefixList {

				found, err := regexp.MatchString(fmt.Sprintf("00%s", prefix), dst)
				if err != nil {
					fmt.Println(err)
				}

				if found == true {
					hits[prefix] = hits[prefix] + 1
				}

			}

		}

		/*fmt.Println(calldate,
		clid,
		src,
		dst,
		dcontext,
		channel,
		dstchannel,
		lastapp,
		lastdata,
		duration,
		billsec,
		disposition,
		amaflags,
		accountcode,
		uniqueid,
		userfield)*/

		/*

			CDR Format for Elastix 2.3 (maybe Asterisk 1.8?)

						calldate    | datetime     | NO   |     | 0000-00-00 00:00:00 |       |
						| clid        | varchar(80)  | NO   |     |                     |       |
						| src         | varchar(80)  | NO   |     |                     |       |
						| dst         | varchar(80)  | NO   |     |                     |       |
						| dcontext    | varchar(80)  | NO   |     |                     |       |
						| channel     | varchar(80)  | NO   |     |                     |       |
						| dstchannel  | varchar(80)  | NO   |     |                     |       |
						| lastapp     | varchar(80)  | NO   |     |                     |       |
						| lastdata    | varchar(80)  | NO   |     |                     |       |
						| duration    | int(11)      | NO   |     | 0                   |       |
						| billsec     | int(11)      | NO   |     | 0                   |       |
						| disposition | varchar(45)  | NO   |     |                     |       |
						| amaflags    | int(11)      | NO   |     | 0                   |       |
						| accountcode | varchar(20)  | NO   |     |                     |       |
						| uniqueid    | varchar(32)  | NO   | MUL |                     |       |
						| userfield

		*/

	}

	fmt.Println(hits)

	for _, hits := range hits {
		if hits >= triggerConfigs.HitThreshold {
			// TODO Run actionChain!
		}
	}

	for t := range ticker.C {

		fmt.Println("dangerousDestinations tick at", t)

		//queryResult, err := stmtOut.Query()
		//fmt.Println(err)

		//test, _ := queryResult.Columns()
		//fmt.Println(test)

	}

}
