package main

import (
	"fmt"
	"time"
)

func RelativelyTime(date time.Time, showDate ...bool) (result string) {
	nowTime := time.Now()
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

func RelativelyTimeSlice(start time.Time, end time.Time, showDate ...bool) string {
	if start.Equal(end) {
		return RelativelyTime(start, showDate...)
	}
	return fmt.Sprintf("%s ~ %s", RelativelyTime(start, showDate...), RelativelyTime(end, showDate...))
}

var longDayNames = []string{"日", "一", "二", "三", "四", "五", "六"}

func TimeWeekdayString(d time.Weekday) string {
	if time.Sunday <= d && d <= time.Saturday {
		return longDayNames[d]
	}
	return d.String()
}
