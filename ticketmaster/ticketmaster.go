package ticketmaster

// https://developer.ticketmaster.com/products-and-docs/apis/discovery-api/v2/

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"
	"time"

	"sweden_events/collector"
)

const (
	defaultCookie = `eps_sid=906579a6a0882c8cf68ce9b2c32945bfeef2131f; pxcts=8cda8268-ed08-11ee-bf5b-0d0b54e43783; _pxvid=8cda7544-ed08-11ee-bf5b-b24ccae1eea0; reese84=3:EUnbQlwP/tWBsBIiKSoneg==:DcyFV2gZvkR0U/lC/Hs0sQ4pw84EqfIyc76CocIRzoXSiw0c4CYnnuivofGMguj1+NO6mDLV0lwzWOl6HAJsR0OeDdr9Z3igTZ+N4LQEPQV2MN1QyLHMypT/Snf9GRC664QYT7s3keF6hPSqgyxG8Q+pu49fBe5nnXoA71szR29VS59DMAHDMLl/PgJjukLm9/i2obuELJyJRyEkyjtCckCtOsfMXYCJccEZMswMxKGbCAJ9S7m9/5TnoXsw9pv5zh8908rNcyAw+8xko3gUaE7u34CSKW4aPBYPWBxkC2LfBgAko+xc6gcg5Q078eEJPoE7OqBXVnvAyu7jMKoOE86d78Vqv09rTzztz0oAJcv/CzpgIpjBzq0rMIU1lOrV/3O4lUqazQjvuZMW3EKxqn3+rE5gna/UqbVFuRSC0mzGR38GBWXv9lWRPE4azrOD/HtAz8jKuhQAt35x+jNZeQ==:i8ri7FVp8fPKf/sW8FizCpsCnP3bOR8r7YA6QSsOlnQ=; language=sv-se; NDMA=612; TMUO=west_RCDnKloRf69yiL837njTi1AOaGcQ1zrMoshNILMpn+8=; OptanonConsent=isGpcEnabled=0&datestamp=Thu+Mar+28+2024+14%3A42%3A25+GMT%2B0100+(Central+European+Standard+Time)&version=202402.1.0&browserGpcFlag=0&isIABGlobal=false&hosts=&consentId=acde5e54-b248-4c60-83c4-13e2dc0c5e74&interactionCount=1&isAnonUser=1&landingPath=NotLandingPage&groups=C0001%3A1%2CC0002%3A1%2CC0003%3A1%2CC0004%3A1%2CV2STACK42%3A1&geolocation=SE%3BM&AwaitingReconsent=false; LANGUAGE=sv-se; SID=GtsyGvkDddPqKhk8-Pmy17W_GYKcCwa8TmSg9CM86tuCLk2_1E7hg_VVreZqis0XSpyw4_hd3YAZaWGQ_BRH; BID=ezSWd_gY94addFjb4aOC227BuuZHWWgYO4ipgPX9T7EelfHhyRQmR5shZYQ8CFasCSWYyQj2apb7qKA; OptanonAlertBoxClosed=2024-03-28T13:39:11.869Z; eupubconsent-v2=CP8LnFgP8LnFgAcABBENAtEsAP_gAAAAACiQJ9JD7D7FbUFCwHpzaLsQMAkHRMCIQoQAAASBAGABQAKQIAQCkkAQFASgBAACAAAAICRBIQAECAAACUAAQAAAIAAEAAAAAAAIAAAAgAAAAAAIAAACAAAAEAAIgAAAEAAAmAgAAAIAGEAAhAAAIAAAAAAAAAAAAgAAAAAAAAIAAAAAACAAAQAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAACAAD8AAABASAoAAsACoAHAAQQAyADQAHgARAAmABVAD0AIQARAAwwB7AD9AMUAcQBSIC8wGTjoDYACwAKgAcABBADIANAAeABEACYAFUAMQAZgA9ACIAFGAMMAZQA9gB-gEWAMUAcQA6gCLwF5gMnAZYA1UcAHAAeABcALoBCACIgHSAX0QgEgALACYAFUAMQAegDFAHUAZOA1UlAOAAWABwAHgARAAmABVADEAIgAUYBigDqAIvAXmAyckABAAuUgKgALAAqABwAEEAMgA0AB4AEQAJgAUgAqgBiADMAIgAUYAygB-gEWAMUAi8BeYDJygAMAC4BdQF9A.f_wAAAAAAAAA; _gcl_au=1.1.1956880293.1711633152; session=1; referrer=https://www.ticketmaster.se/?utm_source=TM-google&gad_source=1&gclid=EAIaIQobChMIxN3bqYWXhQMVfVCRBR1Q5Qq6EAAYASAAEgJzJ_D_BwE&gclsrc=aw.ds; _ga_S7J60HHDPE=GS1.1.1711633152.1.1.1711633345.0.0.0; _ga=GA1.1.1498886711.1711633152`
)

type Venue struct {
	City    string `json:"city"`
	Name    string `json:"name"`
	Country string `json:"country"`
	State   string `json:"state"`
}

type Event struct {
	Title string `json:"title"`
	Venue Venue  `json:"venue"`
}

type TicketmasterData struct {
	Total              int     `json:"total"`
	TotalLocal         int     `json:"totalLocal"`
	TotalInternational int     `json:"totalInternational"`
	Events             []Event `json:"events"`
}

func fetchEvents() (TicketmasterData, error) {
	url := "https://www.ticketmaster.se/api/search/events/category/10001"

	data := TicketmasterData{}

	client := http.Client{
		Timeout: time.Second * 100, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("Cookie", defaultCookie)

	q := req.URL.Query()
	q.Add("region", "612")
	q.Add("page", "0")
	req.URL.RawQuery = q.Encode()

	res, getErr := client.Do(req)
	if getErr != nil {
		return data, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return data, readErr
	}

	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		return data, jsonErr
	}

	return data, nil
}

func CollectEvents() {
	data, err := fetchEvents()
	if err != nil {
		log.Fatal(err)
	}

	events := make([]collector.Event, 0)
	for _, event := range data.Events {
		/*year, _ := strconv.Atoi(event.InformationProgram.Startdatum[0:4])
		month, _ := strconv.Atoi(event.InformationProgram.Startdatum[4:6])
		day, _ := strconv.Atoi(event.InformationProgram.Startdatum[6:8])
		tt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)*/
		tt := time.Date(0, time.Month(0), 0, 0, 0, 0, 0, time.Local)
		events = append(events, collector.CreateEvent(event.Title, tt, event.Venue.Name, "Music", ""))
	}

	collector.AddEvents(events)
}
