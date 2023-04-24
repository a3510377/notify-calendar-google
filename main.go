package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

const (
	UA                   = "notifyGoogleCalendar (https://github.com/a3510377, 1.0.0) Golang/1.19.4"
	DiscordMessageAPIUrl = "https://discord.com/api/channels/%d/messages"
	LineMessageAPIUrl    = "https://notify-api.line.me/api/notify"
)

func getLocation() *time.Location {
	loc, err := time.LoadLocation(os.Getenv("LOC"))
	if err == nil {
		return loc
	}
	return time.Local
}

func getNowTime() time.Time { return time.Now().In(getLocation()) }

func main() {
	godotenv.Load()

	CALENDAR_ID := os.Getenv("CALENDAR_ID")
	if len(os.Args) > 1 {
		CALENDAR_ID = os.Args[1]
	} else if CALENDAR_ID == "" {
		CALENDAR_ID = ConfigData.CALENDAR_ID
	}

	if CALENDAR_ID == "" {
		panic("CALENDAR_ID is empty")
	}

	main := func(checkTimes ...time.Time) {
		checkTime := getNowTime().AddDate(0, 0, 1)
		if len(checkTimes) > 0 {
			checkTime = checkTimes[0]
		}

		log.Println("check", checkTime.Format("2006-01-02"))

		if checkTime.Format("2006-01-02") == GetTmpDate() {
			log.Println("Today already send notification, skip check")
			return
		}

		retryCount := 0
		for ; retryCount < 3; retryCount++ {
			if err := checkAndNotification(CALENDAR_ID, checkTime, nil); err != nil {
				time.Sleep(time.Second * 5) // retry after 5 seconds
				continue
			}
			WriteTmpDate(checkTime)
			break
		}
		if retryCount >= 3 {
			log.Println("Retry 3 times, skip check")
		}

		if ConfigData.Options.AdvanceReminder {
			if err := checkAndNotification(CALENDAR_ID, checkTime.AddDate(0, 0,
				ConfigData.Options.AdvanceReminderDays,
			), func(item CalenderV3ApiEventData) bool {
				// TODO
				return false
			}); err != nil {
				log.Printf("Advance reminder error %s\n", err)
			}
		}
	}

	/* for test */
	// for i := 4; i < 30*4; i++ {
	// 	main(time.Now().AddDate(0, 0, i))
	// 	time.Sleep(time.Second)
	// }
	// return

	const specTime = "50 14 * * *"

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := parser.Parse(specTime)
	if now := getNowTime(); s.Next(now).In(getLocation()).Day() != now.Day() {
		main()
	}

	c := cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(log.Default())),
		cron.WithLocation(getLocation()),
	)

	c.AddFunc(specTime, func() { main() })

	c.Run() // loop start
}

func checkAndNotification(CALENDAR_ID string, nowTime time.Time, callback func(item CalenderV3ApiEventData) bool) error {
	resp, err := http.Get(NewCalendarV3ApiRequest(nowTime, CALENDAR_ID).BaseURL().String())
	if err != nil {
		log.Println("Error getting calendar data: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Println("Error getting calendar data: ", resp.Status)
		return errors.New("Error getting calendar data: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading calendar data: ", err)
		return err
	}

	data := CalenderV3ApiResponse{}
	json.Unmarshal(body, &data)
	log.Println("Calendar data: ", data)

	notifications := map[string][]CalenderV3ApiEventData{}
	for _, item := range data.Items {
		if callback == nil {
			// check if the item start time is before the current time and the status is confirmed
			if !item.IsSameStartDay(nowTime) || item.Status != "confirmed" {
				continue
			}
		} else if !callback(item) {
			continue
		}

		key := item.StartTimeString() + "-" + item.EndTimeString()
		notifications[key] = append(notifications[key], item)
	}

	for _, item := range notifications {
		notification(nowTime.Add(-time.Hour*24), item...)
	}

	return nil
}
