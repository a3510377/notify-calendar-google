package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/araddon/dateparse"
)

const (
	GoogleCalendarAPIBaseURL = "https://clients6.google.com/calendar/v3/calendars/"
	// is google api default key
	GoogleCalendarBaseKey = "AIzaSyBNlYH01_9Hc5S1J9vuFmu2nUqBZJNAXxs"
)

// https://developers.google.com/calendar/api/v3/reference/events/list
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

func (c *CalendarV3ApiRequest) Get() (result map[string]any) {
	data, _ := json.Marshal(c)

	result = map[string]any{}
	json.Unmarshal(data, &result)
	return
}

func (c *CalendarV3ApiRequest) BaseURL() *url.URL {
	baseUrl, _ := url.Parse(GoogleCalendarAPIBaseURL + c.CalendarID + "/events")

	query := baseUrl.Query()
	for key, value := range c.Get() {
		query.Add(key, fmt.Sprintf("%v", value))
	}

	baseUrl.RawQuery = query.Encode()
	return baseUrl
}

func timeFormat(time time.Time) string { return time.Format("2006-01-02T00:00:00Z07:00") }

func NewCalendarV3ApiRequest(date time.Time, calendarID string) *CalendarV3ApiRequest {
	nowTime := time.Unix(int64(float64(date.Unix()/1e4)*1e4), 0)

	return &CalendarV3ApiRequest{
		CalendarID:   calendarID,
		TimeMin:      timeFormat(nowTime),
		TimeMax:      timeFormat(nowTime.AddDate(0, 0, 1)),
		MaxResults:   250,
		SingleEvents: true,
		MaxAttendees: 1,
		SanitizeHtml: true,
		Key:          GoogleCalendarBaseKey,
	}
}

// https://developers.google.com/calendar/api/v3/reference/events/list#response
type CalenderV3ApiResponse struct {
	Summary     string                   `json:"summary"`
	Description string                   `json:"description"`
	Updated     string                   `json:"updated"`
	Items       []CalenderV3ApiEventData `json:"items"`
}

// https://developers.google.com/calendar/api/v3/reference/events#resource
type CalenderV3ApiEventData struct {
	HtmlLink    string `json:"htmlLink"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Status      string `json:"status"` // confirmed, tentative, cancelled :: ??????, ??????, ??????
	Start       struct {
		Date     string `json:"date"`
		DateTime string `json:"dateTime"`
	} `json:"start"`
	End struct {
		Date     string `json:"date"`
		DateTime string `json:"dateTime"`
	} `json:"end"`
}

func (c *CalenderV3ApiEventData) StartTime() time.Time {
	dataDateString := c.Start.Date
	if dataDateString == "" {
		dataDateString = c.Start.DateTime
	}
	dataDate, _ := dateparse.ParseAny(dataDateString)
	return dataDate
}

func (c *CalenderV3ApiEventData) EndTime() time.Time {
	dataDateString := c.End.Date
	if dataDateString == "" {
		dataDateString = c.End.DateTime
	}
	dataDate, _ := dateparse.ParseAny(dataDateString)
	return dataDate
}

func (c *CalenderV3ApiEventData) timeFormat(dataDate time.Time) string {
	if c.IsAllDay() {
		return dataDate.Format("2006-01-02")
	}
	return dataDate.Format("2006-01-02 15:04:05")
}

func (c *CalenderV3ApiEventData) StartTimeString() string { return c.timeFormat(c.StartTime()) }
func (c *CalenderV3ApiEventData) EndTimeString() string   { return c.timeFormat(c.EndTime()) }

func (c *CalenderV3ApiEventData) IsSameStartDay(date time.Time) bool {
	return c.StartTime().Format("2006-01-02") == date.Format("2006-01-02")
}

func (c *CalenderV3ApiEventData) IsSameEndDay(date time.Time) bool {
	return c.EndTime().Format("2006-01-02") == date.Format("2006-01-02")
}

func (c *CalenderV3ApiEventData) IsAllDay() bool {
	return c.Start.Date != "" && c.End.Date != ""
}
