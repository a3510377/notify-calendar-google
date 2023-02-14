package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// c_nbtiskrng1pkrcj168db62l4hg@group.calendar.google.com

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

	c := cron.New()
	// TODO add config cron rule
	c.AddFunc("0 0 12 * * ?", func() { checkAndNotification(CALENDAR_ID) })
	c.Start() // loop start
}

func notification(data CalenderV3ApiEventData) {
}

func checkAndNotification(CALENDAR_ID string) {
	nowTime := time.Now().AddDate(0, 0, 2)
	resp, err := http.Get(NewCalendarV3ApiRequest(nowTime, CALENDAR_ID).BaseURL().String())
	if err != nil {
		log.Println("Error getting calendar data: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Error getting calendar data: ", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading calendar data: ", err)
		return
	}
	data := CalenderV3ApiResponse{}
	json.Unmarshal(body, &data)

	log.Println("Calendar data: ", data)
	for _, item := range data.Items {
		// check if the item start time is before the current time and the status is confirmed
		if nowTime.Format("2006-01-02") == item.Start.Date && item.Status == "confirmed" {
			notification(item)
		}
	}
}
