package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/matrosov-nikita/alarm_generator/generator"

	"github.com/matrosov-nikita/alarm_generator/postgres"
)

func main() {
	var startDateStr, endDateStr string
	var alertsCount, serversCount int
	flag.IntVar(&alertsCount, "alertsCount", 2000, "Alerts count to generate for one server")
	flag.IntVar(&serversCount, "servers", 4, "Servers count")
	flag.StringVar(&startDateStr, "startDate", "2020-01-01", "Start date for events generation")
	flag.StringVar(&endDateStr, "endDate", "", "End date for events generation")
	flag.Parse()
	if alertsCount <= 0 {
		log.Fatal("events count must be positive")
	}

	if serversCount <= 0 {
		log.Fatal("servers count must be positive")
	}

	client, err := postgres.NewClient("postgresql://postgres:postgres@127.0.0.1:5432/generator?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to create postgres client: %+v", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Fatal("dates must be at YYYY-MM-DD format")
	}

	endDate := time.Now()
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			log.Fatal("dates must be at YYYY-MM-DD format")
		}
	}

	fmt.Printf("alertsCount=%d\nstartDate=%s\nendDate=%s\nserversCount=%d\n",
		alertsCount, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), serversCount)

	gen := generator.New(client, startDate, endDate, alertsCount, serversCount)
	startTime := time.Now()
	if err := gen.LoadEvents(); err != nil {
		log.Fatalf("failed to load events to db: %+v", err)
	}
	timeElapsed := time.Since(startTime)
	log.Printf("Loading successfully completed, time elapsed %.3f sec.", timeElapsed.Seconds())
}
