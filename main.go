package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

		for retryCount := 0; retryCount < 3; retryCount++ {
			if checkTime.Format("2006-01-02") == GetTmpDate() {
				log.Println("Today already send notification, skip check")
				break
			}
			if err := checkAndNotification(CALENDAR_ID, checkTime); err != nil {
				retryCount++

				if retryCount >= 3 {
					log.Println("Retry 3 times, skip check")
				}
				time.Sleep(time.Second * 5) // retry after 5 seconds
				continue
			}
			WriteTmpDate(checkTime)
			break
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

func checkAndNotification(CALENDAR_ID string, nowTime time.Time) error {
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
		// check if the item start time is before the current time and the status is confirmed
		if !item.IsSameStartDay(nowTime) || item.Status != "confirmed" {
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

func NotifyLine(text string) {
	TOKEN := ConfigData.Line.TOKEN

	if TOKEN == "" {
		log.Println("Line token is empty")
		return
	}

	data := url.Values{"message": {"\n" + text}}.Encode()
	req, _ := http.NewRequest("POST", LineMessageAPIUrl, strings.NewReader(data))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))
	req.Header.Set("Authorization", "Bearer "+TOKEN)
	req.Header.Set("User-Agent", UA)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error send Line notification: %s\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		data, _ := io.ReadAll(resp.Body)
		log.Printf("Error send Line notification: %s\nResponse: %s\nSend: ", err, data)
		data, _ = io.ReadAll(resp.Request.Body)
		log.Println(string(data))
	}
}

func NotifyDiscord(text string) {
	discordConfig := ConfigData.Discord
	TOKEN := discordConfig.TOKEN

	contentByte, _ := json.Marshal(map[string]string{"content": text})
	bodyReader := bytes.NewReader(contentByte)

	if TOKEN == "" {
		log.Println("Discord token is empty")
	} else {
		for _, id := range discordConfig.ChannelIDs {
			// multiple concurrent requests
			go func(data bytes.Reader, id int64) { // id is channel ID
				req, _ := http.NewRequest("POST", fmt.Sprintf(DiscordMessageAPIUrl, id), &data)

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bot "+TOKEN)
				req.Header.Set("User-Agent", UA)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("Error send discord: %s\nID: %d\n", err, id)
					return
				}

				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					data, _ := io.ReadAll(resp.Body)
					log.Printf("Error send discord: %s\nID: %d\nResponse: %s\nSend: ", err, id, data)
					data, _ = io.ReadAll(resp.Request.Body)
					log.Println(string(data))
				}
			}(*bodyReader, id)
		}
	}

	for _, url := range discordConfig.Webhooks {
		// multiple concurrent requests
		go func(data bytes.Reader, url string) {
			req, _ := http.NewRequest("POST", url, &data)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", UA)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Error send discord webhook: ", err, "\nURL:", url)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				data, _ := io.ReadAll(resp.Body)
				log.Printf("Error send discord webhook: %s\nURL: %s\nResponse: %s\n", resp.Status, url, data)
			}
		}(*bodyReader, url)
	}
}
