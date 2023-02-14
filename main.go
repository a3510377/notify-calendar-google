package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// c_nbtiskrng1pkrcj168db62l4hg@group.calendar.google.com
// https://developers.google.com/calendar/api/v3/reference/events/list

const (
	GoogleCalendarAPIBaseURL = "https://clients6.google.com/calendar/v3/calendars/"
	// is google api default key
	GoogleCalendarBaseKey = "AIzaSyBNlYH01_9Hc5S1J9vuFmu2nUqBZJNAXxs"
)

func main() {
	godotenv.Load()
	var CALENDAR_ID string
	if len(os.Args) > 1 {
		CALENDAR_ID = os.Args[1]
	} else {
		CALENDAR_ID = os.Getenv("CALENDAR_ID")
	}

	if CALENDAR_ID == "" {
		panic("CALENDAR_ID is empty")
	}

	baseLoop := func() {
		nowTime := time.Now().AddDate(0, 0, 2)
		resp, err := http.Get(NewCalendarV3ApiRequest(nowTime, CALENDAR_ID).BaseURL().String())
		if err != nil {
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		data := CalenderV3ApiResponse{}
		json.Unmarshal(body, &data)

		for _, item := range data.Items {
			// check if the item start time is before the current time and the status is confirmed
			if nowTime.Format("2006-01-02") == item.Start.Date && item.Status == "confirmed" {
				// TODO: send notification
			}
		}
	}
	baseLoop()
	for range time.Tick(time.Hour * 24) { // 24 hour clock
		baseLoop()
	}
}
