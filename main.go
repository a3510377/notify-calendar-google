package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

const (
	UA                   = "notify-calendar-google(https://github.com/a3510377/notify-calendar-google,1.0.0)"
	DiscordMessageAPIUrl = "https://discord.com/api/channels/%d/messages"
)

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
	main := func() {
		retryCount := 0
		for retryCount < 3 {
			if err := checkAndNotification(CALENDAR_ID); err == nil {
				break
			}
			retryCount++
			time.Sleep(time.Second * 5) // retry after 5 seconds
		}
	}

	main() // run once
	time.Sleep(time.Second * 10)
	main()

	// TODO add config cron rule
	c.AddFunc("0 0 12 * * ?", main)

	c.Run() // loop start
}

func checkAndNotification(CALENDAR_ID string) error {
	nowTime := time.Now().AddDate(0, 0, 2)
	if nowTime.Format("2006-01-02") == GetTmpDate() {
		log.Println("Today already send notification, skip check")
		return nil
	}
	WriteTmpDate(nowTime)
	resp, err := http.Get(NewCalendarV3ApiRequest(nowTime, CALENDAR_ID).BaseURL().String())
	if err != nil {
		log.Println("Error getting calendar data: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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
	for _, item := range data.Items {
		// check if the item start time is before the current time and the status is confirmed
		if nowTime.Format("2006-01-02") == item.Start.Date && item.Status == "confirmed" {
			notification(item)
		}
	}
	return nil
}

func notification(data CalenderV3ApiEventData) {
	log.Println("Send: ", data.Summary)

	// discord
	if ConfigData.Discord.Enable {
		NotifyDiscord(data)
	}

	// line notify
	// TODO send line use line notify
}

func NotifyDiscord(data CalenderV3ApiEventData) {
	discordConfig := ConfigData.Discord
	TOKEN := discordConfig.TOKEN

	if TOKEN == "" {
		log.Println("Discord token is empty")
	} else {
		for _, id := range discordConfig.ChannelIDs {

			req, _ := http.NewRequest("POST", fmt.Sprintf(DiscordMessageAPIUrl, id), nil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bot "+TOKEN)
			req.Header.Set("User-Agent", UA)

			// multiple concurrent requests
			go func() {
				resp, err := (&http.Client{}).Do(req)
				if err != nil {
					log.Println("Error send discord: ", err)
					return
				}

				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					log.Println("Error send discord: ", resp.Status)
				}
			}()
		}
	}

	for _, url := range discordConfig.Webhooks {
		req, _ := http.NewRequest("POST", url, nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", UA)

		// multiple concurrent requests
		go func() {
			resp, err := (&http.Client{}).Do(req)
			if err != nil {
				log.Println("Error send discord webhook: ", err)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Println("Error send discord webhook: ", resp.Status)
			}
		}()
	}
}
