package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db/clickhouse"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db/postgres"

	"github.com/matrosov-nikita/smart-generator/pkg/config"

	"github.com/matrosov-nikita/smart-generator/generator"
)

func main() {
	var startDateStr, endDateStr, generatorType string
	var serversCount, teamsCount int
	var configPath string
	// TODO: move this to config
	flag.IntVar(&serversCount, "servers", 4, "Servers count")
	flag.IntVar(&teamsCount, "teams", 10, "Teams count")
	flag.StringVar(&configPath, "configPath", "./config.json", "path to config file")
	flag.StringVar(&startDateStr, "startDate", "2020-01-01", "Start date for events generation")
	flag.StringVar(&endDateStr, "endDate", "", "End date for events generation")
	flag.StringVar(&generatorType, "generatorType", "normal", "Specifies how to generate time seq for events")
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

	client, err := postgres.NewClient("postgresql://postgres:postgres@127.0.0.1:5432/generator?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to create postgres client: %+v", err)
	}

	chClient, err := clickhouse.NewClient("tcp://127.0.0.1:9000?debug=false")
	if err != nil {
		log.Fatalf("failed to create clickhouse client: %+v", err)
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
	gen := generator.New(client, startDate, endDate, serversCount, teamsCount, generatorType, detectorConfig)
	startTime := time.Now()
	gen.Run()
	timeElapsed := time.Since(startTime)
	log.Printf("Loading to POSTGRES successfully completed, time elapsed %.3f sec.", timeElapsed.Seconds())

	gen = generator.New(chClient, startDate, endDate, serversCount, teamsCount, generatorType, detectorConfig)
	startTime = time.Now()
	gen.Run()
	timeElapsed = time.Since(startTime)
	log.Printf("Loading to CLICKHOUSE successfully completed, time elapsed %.3f sec.", timeElapsed.Seconds())
}
