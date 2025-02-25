package main

import (
	"sweden_events/collector"
	//"sweden_events/cruncho"
	"sweden_events/royalarena"
	"sweden_events/slagthuset"
	"sweden_events/ticketmaster"
)

func main() {
	slagthuset.CollectEvents()
	royalarena.CollectEvents()
	//cruncho.CollectEvents()
	ticketmaster.CollectEvents()

	collector.PrintEvents()
}
