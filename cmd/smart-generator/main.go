package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db"

	"github.com/matrosov-nikita/smart-generator/pkg/config"

	"github.com/matrosov-nikita/smart-generator/generator"
)

func main() {
	var startDateStr, endDateStr string
	var serversCount, teamsCount int
	var configPath string
	flag.IntVar(&serversCount, "servers", 4, "Servers count")
	flag.IntVar(&teamsCount, "teams", 10, "Teams count")
	flag.StringVar(&configPath, "configPath", "./config.json", "path to config file")
	flag.StringVar(&startDateStr, "startDate", "2020-01-01", "Start date for events generation")
	flag.StringVar(&endDateStr, "endDate", "", "End date for events generation")
	flag.Parse()

	detectorConfig, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if serversCount <= 0 {
		log.Fatal("servers count must be positive")
	}

	if teamsCount <= 0 {
		log.Fatal("teams count must be positive")
	}

	client, err := db.NewClient("postgresql://postgres:postgres@127.0.0.1:5432/generator?sslmode=disable")
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

	fmt.Printf("startDate=%s\nendDate=%s\nserversCount=%d\ndetectorConfig=%+v\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), serversCount, detectorConfig)
	gen := generator.New(client, startDate, endDate, serversCount, teamsCount, detectorConfig)
	startTime := time.Now()
	gen.Run()
	timeElapsed := time.Since(startTime)
	log.Printf("Loading successfully completed, time elapsed %.3f sec.", timeElapsed.Seconds())
}
