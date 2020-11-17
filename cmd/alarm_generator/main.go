package main

import (
	"flag"
	"log"
	"time"

	"github.com/matrosov-nikita/alarm_generator/generator"

	"github.com/matrosov-nikita/alarm_generator/postgres"
)

func main() {
	var startDateStr, endDateStr string
	count := 100000
	flag.IntVar(&count, "count", 100000, "Events count to generate")
	flag.StringVar(&startDateStr, "startDate", "2020-01-01", "Start date for events generation")
	flag.StringVar(&endDateStr, "endDate", "", "End date for events generation")
	flag.Parse()
	if count <= 0 {
		log.Fatal("events count must be positive")
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

	gen := generator.New(client, startDate, endDate, count)
	startTime := time.Now()
	if err := gen.LoadEvents(); err != nil {
		log.Fatalf("failed to load events to db: %+v", err)
	}
	timeElapsed := time.Since(startTime)
	log.Printf("Loading successfully completed, time elapsed %.3f sec.", timeElapsed.Seconds())
}
