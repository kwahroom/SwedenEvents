package collector

import (
	"fmt"
	"sort"
	"time"
)

type Event struct {
	Title    string
	Date     time.Time
	Location string
	Category string
	Info string
}

var AllEvents []Event = make([]Event, 0)

func CreateEvent(title string, date time.Time, location string, category string, info string) Event {
	return Event{
		Title:    title,
		Date:     date,
		Location: location,
		Category: category,
		Info: info,
	}
}

func AddEvent(title string, date time.Time, location string, category string, info string) {
	AllEvents = append(AllEvents, CreateEvent(title, date, location, category, info))
}

func AddEvents(events []Event) {
	AllEvents = append(AllEvents, events...)

	sort.SliceStable(AllEvents, func(i, j int) bool {
		return AllEvents[i].Date.Before(AllEvents[j].Date)
	})
}

func PrintEvents() {

	year := 0
	var month time.Month

	fmt.Println(`Auto-generated page with scraper.

Objective is to collect as many events in Sweden ( or close enough ) as possible
from Various sources. Currently it only covers the following:

* Slagthuset in Malmö
* Royal Arena in Copenhagen`)

	for _, event := range AllEvents {

		if event.Date.Year() != year {
			year = event.Date.Year()
			fmt.Println("\n# ", year)
		}

		if event.Date.Month() != month {
			month = event.Date.Month()
			fmt.Println("\n## ", month)
		}

		fmt.Println("*", event.Date.Day(), ".", event.Title, "[", event.Location, "] [", event.Category, "]", "[", event.Info, "]")
	}
}
