package royalarena

import (
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"

	"sweden_events/collector"
)

// https://www.zenrows.com/blog/web-scraping-golang

func getType(e *colly.HTMLElement) string {
	class := e.Attr("class")
	switch {

	case strings.Contains(class, "cid-1"):
		return "Music"
	case strings.Contains(class, "cid-2"):
		return "Sport"
	case strings.Contains(class, "cid-3"):
		return "Theater & Entertainment"
	case strings.Contains(class, "cid-4"):
		return "Family"
	case strings.Contains(class, "cid-5"):
		return "Comedy"
	default:
		return "unknown"
	}
}

func CollectEvents() {
	c := colly.NewCollector()

	events := make([]collector.Event, 0)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML(".event-row", func(e *colly.HTMLElement) {

		date := e.ChildText(".date-small")
		title := e.ChildText(".title")

		if title == "" {
			return
		}

		// format 28.03.2024 - ....
		year, _ := strconv.Atoi(date[6:10])
		month, _ := strconv.Atoi(date[3:5])
		day, _ := strconv.Atoi(date[0:2])

		tt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

		events = append(events, collector.CreateEvent(title, tt, "Royal Arena Copenhagen", getType(e), ""))
	})

	c.Visit("https://www.royalarena.dk/en/events")

	c.Wait()

	collector.AddEvents(events)
}
