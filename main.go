package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// c_nbtiskrng1pkrcj168db62l4hg@group.calendar.google.com
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
	// TODO add config cron rule
	c.AddFunc("0 0 12 * * ?", func() { checkAndNotification(CALENDAR_ID) })
	c.Start() // loop start
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

	for _, id := range discordConfig.ChannelIDs {
		if TOKEN == "" {
			log.Println("Discord token is empty")
			return
		}

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
