package main

import (
	"flag"
	"log"
	"time"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db/postgres"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db/clickhouse"

	"github.com/matrosov-nikita/smart-generator/pkg/client/http"

	"github.com/matrosov-nikita/smart-generator/pkg/config"

	"github.com/matrosov-nikita/smart-generator/generator"
)

const (
	ClickhouseURL = "tcp://127.0.0.1:9000?debug=false"
	PostgresURL   = "postgresql://postgres@127.0.0.1:5432/generator?sslmode=disable"
	HttpURL       = "http://127.0.0.1:8780/api/v1/data/any/push"
)

func main() {
	cfg := &config.Config{}
	flag.IntVar(&cfg.ServersCount, "servers", 1, "Servers count")
	flag.IntVar(&cfg.TeamsCount, "teams", 2, "Teams count")
	flag.IntVar(&cfg.BatchSize, "batchSize", 100, "Batch size")
	flag.StringVar(&cfg.DetectorsConfigPath, "detectorsConfigPath", "./config.json", "Path to file stores amount events per detector")
	flag.StringVar(&cfg.StartDateStr, "startDate", "2020-01-01", "Start date for events generation")
	flag.StringVar(&cfg.EndDateStr, "endDate", "", "End date for events generation, default: now")
	flag.StringVar(&cfg.TimeGeneratorType, "generatorType", "normal", "Generate type for time seq for events: normal or random")
	flag.StringVar(&cfg.StorageType, "storageType", "all", "clickhouse, postgres, http, all")
	flag.Parse()

	if err := cfg.ParseFields(); err != nil {
		log.Fatal(err)
	}

	log.Printf("start with config: %+v", cfg)
	pgClient, err := postgres.NewClient(PostgresURL)
	if err != nil {
		log.Printf("failed to create postgres client: %+v", err)
	}

	chClient, err := clickhouse.NewClient(ClickhouseURL)
	if err != nil {
		log.Printf("failed to create clickhouse client: %+v", err)
	}

	httpClient := http.NewClient(HttpURL)
	defer func() {
		if err := pgClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	defer func() {
		if err := chClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	switch cfg.StorageType {
	case "clickhouse":
		runGeneratorWithClient(cfg.StorageType, chClient, cfg)
	case "postgres":
		runGeneratorWithClient(cfg.StorageType, pgClient, cfg)
	case "http":
		runGeneratorWithClient("http", httpClient, cfg)
	case "all":
		runGeneratorWithClient("clickhouse", chClient, cfg)
		runGeneratorWithClient("postgres", pgClient, cfg)
		runGeneratorWithClient("http", httpClient, cfg)
	}
}

func runGeneratorWithClient(clientType string, client generator.Client, cfg *config.Config) {
	gen := generator.New(client, cfg)
	startTime := time.Now()
	gen.Run()
	timeElapsed := time.Since(startTime)
	log.Printf("Loading to %s successfully completed, time elapsed %.3f sec.", clientType, timeElapsed.Seconds())
}
