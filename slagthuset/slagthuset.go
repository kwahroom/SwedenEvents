package slagthuset

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"sweden_events/collector"
)

type Event struct {
	ID		int `json:"id"`
	Title   struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Acf		struct {
		StartDatum string `json:"startdatum"`
	} `json:"acf"`
}
/*
{
    "id": 5156,
    "slug": "fairytale-of-new-york",
    "title": {
      "rendered": "Fairytale of New York"
    },
    "featured_media": 5167,
    "acf": {
      "startdatum": "20251125",
      "slutdatum": "",
      "oppnar": "18:00",
      "inslapp_till_salong": "19:00",
      "borjar": "19:30",
      "slutar": "",
  }*/


func fetchEvents() ([]Event, error) {
	url := "https://www.slagthuset.se/api/events"

	eventList := make([]Event, 0)

	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return eventList, err
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		return eventList, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return eventList, readErr
	}

	jsonErr := json.Unmarshal(body, &eventList)
	if jsonErr != nil {
		return eventList, jsonErr
	}

	return eventList, nil
}

func CollectEvents() {

	eventList, err := fetchEvents()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()

	events := make([]collector.Event, 0)
	for _, event := range eventList {
		year, _ := strconv.Atoi(event.Acf.StartDatum[0:4])
		month, _ := strconv.Atoi(event.Acf.StartDatum[4:6])
		day, _ := strconv.Atoi(event.Acf.StartDatum[6:8])
		tt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

		if tt.Before(now) {
			continue
		}

		events = append(events, collector.CreateEvent(event.Title.Rendered, tt, "Slagthuset malmö", ""))
	}

	collector.AddEvents(events)
}
