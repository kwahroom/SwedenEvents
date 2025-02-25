package cruncho

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	"time"
	//	"sweden_events/collector"
)

const (
	defaultJson = `
{
	"pageContext":{
		"destinationSlug":"lund",
		"l1":"events",
		"previousL1":"",
		"clientTime":"17:21",
		"ip":"",
		"area":"All areas"
	},
	"endDate":"2044-03-27T12:28:59.999Z",
	"l2":["other-events"],
	"l3":[
		"arts-exhibition",
		"christmas-concerts",
		"comedy-quiz",
		"dancing",
		"diy-crafts",
		"food-drink",
		"guided-tours",
		"holiday-and-festival",
		"international-citizen-hub",
		"kamratkortet","kids",
		"lecture",
		"literature-writing",
		"litteraturdagarna",
		"lov-i-lund",
		"movies-film",
		"music",
		"nature",
		"online-event-and-esport",
		"online-other-events",
		"rysliga-veckan",
		"sports-fitness",
		"theater",
		"young-adults"
	],
	"format":"",
	"startDate":"2024-03-28T23:00:59.999Z",
	"timezone":"Europe/Stockholm",
	"handpicked":false
}`
)

type CrunchoData struct {
	Address               string   `json:"address"`
	ApiRatings            []string `json:"apiRatings"`
	AvailableApis         []string `json:"availableApis"`
	AvailableTranslations []string `json:"availableTranslations"`
	BookProvider          []string `json:"bookProvider"`
	BookUrl               string   `json:"bookUrl"`
	Categories            []string `json:"categories"`
	City                  string   `json:"city"`
	Description           string   `json:"description"`
	DestinationSlug       string   `json:"destinationSlug"`
	Email                 string   `json:"email"`
	EventEnd              []string `json:"eventEnd"`
	EventStart            []string `json:"eventStart"`
	EventTargetGroup      []string `json:"eventTargetGroup"`
	EventVenueName        string   `json:"eventVenueName"`
	Formats               []string `json:"formats"`
	Geometry              struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"geometry"`
	Hide          bool   `json:"hide"`
	Hybrid        bool   `json:"hybrid"`
	IsFree        bool   `json:"isFree"`
	IsSponsored   bool   `json:"isSponsored"`
	IsPinned      bool   `json:"isPinned"`
	Label         string `json:"label"`
	LastUpdated   string `json:"lastUpdated"`
	LongTermEvent bool   `json:"longTermEvent"`
	L1            string `json:"l1"`
	Name          string `json:"name"`
	Online        bool   `json:"online"`
	Organizer     string `json:"organizer"`
	Photos        []struct {
		Url    string `json:"url"`
		Alt    string `json:"alt"`
		Credit string `json:"credit"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Api    string `json:"api"`
	} `json:"photos"`
	PostDate    string   `json:"postDate"`
	Price       string   `json:"price"`
	Reviews     []string `json:"reviews"`
	ShouldMatch bool     `json:"shouldMatch"`
	tag         string   `json:"tag"`
	Videos      []string `json:"videos"`
	WebsiteUrl  string   `json:"websiteUrl"`
	Id          string   `json:"id"`
}

func fetchEvents() ([]CrunchoData, error) {
	url := "https://api-ts.cruncho.co/landing-page/recommendations"

	events := make([]CrunchoData, 0)

	client := http.Client{
		Timeout: time.Second * 100, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(defaultJson))
	if err != nil {
		return events, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	q := req.URL.Query()
	q.Add("destination", "lund")
	q.Add("size", "1")
	q.Add("sponsored", "false")
	req.URL.RawQuery = q.Encode()

	res, getErr := client.Do(req)
	if getErr != nil {
		return events, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return events, readErr
	}

	jsonErr := json.Unmarshal(body, &events)
	if jsonErr != nil {
		return events, jsonErr
	}

	return events, nil
}

func CollectEvents() {
	events, _ := fetchEvents()
	fmt.Println(events)
}
