package main

import (
	"fmt"
	"strings"
	"time"
)

func RelativelyTime(nowTime time.Time, date time.Time, showDate ...bool) (result string) {
	subDuration := date.Sub(nowTime)
	ShowDate := len(showDate) > 0 && showDate[0]

	// // seconds
	// if subDuration.Seconds() < 60 {
	// 	return fmt.Sprintf("%d 秒左右", int(subDuration.Seconds()))
	// }

	// // minutes
	// if subDuration.Hours() < 1 {
	// 	return fmt.Sprintf("%d 分鐘左右", int(subDuration.Minutes()))
	// }

	// // hours
	// if subDuration.Hours() < 24 {
	// 	return fmt.Sprintf("%d 小時左右", int(subDuration.Hours()))
	// }

	// days
	dayHour := time.Hour * 24
	if dateDay := date.Day(); dateDay == nowTime.Day() {
		result = "今天"
	} else if nextDay := nowTime.Add(dayHour); dateDay == nextDay.Day() {
		result = "明天"
	} else if dateDay == nextDay.Add(dayHour).Day() {
		result = "後天"
	} else if subDuration.Hours() < 24*7 {
		result = fmt.Sprintf("%d 天後", int(subDuration.Hours()/24))
	}
	if result != "" {
		if ShowDate {
			result += fmt.Sprintf(" (%s)", date.Format("2006-01-02"))
		}
		return
	}

	// weeks
	if subDuration.Hours() < 24*7*2 {
		result = "下週" + TimeWeekdayString(date.Weekday())
	} else {
		result = fmt.Sprintf("%d 週後", int(subDuration.Hours()/(24*7)))
	}
	if ShowDate {
		result += fmt.Sprintf(" (%s)", date.Format("2006-01-02"))
	}
	return
}

func RelativelyTimeSlice(fromTime time.Time, start time.Time, end time.Time, showDate ...bool) string {
	arg1 := RelativelyTime(fromTime, start, showDate...)
	if start.Format("2006-01-02") == end.Format("2006-01-02") { // as same day
		return arg1
	}
	return fmt.Sprintf("%s ~ %s", arg1, RelativelyTime(fromTime, end, showDate...))
}

var longDayNames = []string{"日", "一", "二", "三", "四", "五", "六"}

func TimeWeekdayString(d time.Weekday) string {
	if time.Sunday <= d && d <= time.Saturday {
		return longDayNames[d]
	}
	return d.String()
}

func notification(fromTime time.Time, data ...CalenderV3ApiEventData) {
	content := ""

	for _, item := range data {
		description := ""
		if item.Description != "" {
			data := strings.Split(item.Description, "\n")
			description += " >>> \n"
			for _, item := range data {
				description += "   >> " + item + "\n"
			}
		}

		endTime := item.EndTime()
		if item.Start.Date != "" {
			endTime = endTime.Add(-time.Hour * 24)
		}
		content += fmt.Sprintf("%s是 %s 的日子 %s\n", RelativelyTimeSlice(
			fromTime, item.StartTime(),
			endTime,
		), item.Summary, description)
	}

	content = strings.TrimSuffix(content, "\n") // remove trailing newline

	// line notify
	if ConfigData.Line.Enable {
		NotifyLine(content)
	}

	// discord
	if ConfigData.Discord.Enable {
		NotifyDiscord(content)
	}
}
