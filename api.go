package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

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
	HtmlLink string `json:"htmlLink"`
	Summary  string `json:"summary"`
	Status   string `json:"status"` // confirmed, tentative, cancelled :: 確認, 暫定, 取消
	Start    struct {
		Date string `json:"date"`
	} `json:"start"`
	End struct {
		Date string `json:"date"`
	} `json:"end"`
}
