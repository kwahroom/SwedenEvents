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
	ID    int `json:"id"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Acf struct {
		StartDatum string `json:"startdatum"`
		Pris string `json:"pris"`
		TypAvEvenemang struct {
			Name string `json:"name"`
		} `json:"typ_av_evenemang"`
	} `json:"acf"`
}

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

		events = append(events, collector.CreateEvent(event.Title.Rendered, tt, "Slagthuset malmö", event.Acf.TypAvEvenemang.Name, event.Acf.Pris))
	}

	collector.AddEvents(events)
}
