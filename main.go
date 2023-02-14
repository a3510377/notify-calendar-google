package main

import (
	"fmt"
	"time"
)

func main() {
	baseLoop := func() {
		fmt.Println("Foo")
	}
	baseLoop()
	for range time.Tick(time.Hour * 24) { // 24 hour clock
		baseLoop()
	}
	// data := CalendarV3ApiRequest{
	// 	CalendarID: os.Getenv("CALENDAR_ID"),
	// }
	// http.Get(data.BaseURL())
}

type CalendarV3ApiRequest struct {
	CalendarID   string `json:"calendarId"`
	SingleEvents bool   `json:"singleEvents"`
	TimeZone     string `json:"timeZone"`
	MaxAttendees int    `json:"maxAttendees"`
	MaxResults   int    `json:"maxResults"`
	SanitizeHtml bool   `json:"sanitizeHtml"`
	TimeMin      string `json:"timeMin"`
	TimeMax      string `json:"timeMax"`
	Key          string `json:"key"`
}

func (c *CalendarV3ApiRequest) JSON() CalendarV3ApiRequest {
	return CalendarV3ApiRequest{
		SingleEvents: true,
		TimeZone:     "Asia/Taipei",
		MaxAttendees: 1,
		SanitizeHtml: true,
		// TODO check is google api default key not auto generated
		Key: "AIzaSyBNlYH01_9Hc5S1J9vuFmu2nUqBZJNAXxs",
	}
}

func (c *CalendarV3ApiRequest) BaseURL() string {
	return fmt.Sprintf("https://clients6.google.com/calendar/v3/calendars/%s/events", c.CalendarID)
}

type CalenderV3ApiResult struct{}

// https://clients6.google.com/calendar/v3/calendars/
// c_nbtiskrng1pkrcj168db62l4hg@group.calendar.google.com/events?calendarId=c_nbtiskrng1pkrcj168db62l4hg%40group.calendar.google.com&singleEvents=true&timeZone=Asia%2FTaipei&maxAttendees=1&maxResults=250&sanitizeHtml=true&timeMin=2023-01-29T00%3A00%3A00%2B08%3A00&timeMax=2023-03-05T00%3A00%3A00%2B08%3A00&key=AIzaSyBNlYH01_9Hc5S1J9vuFmu2nUqBZJNAXxs
