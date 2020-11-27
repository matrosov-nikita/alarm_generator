package main

import (
	"flag"
	"log"
	"time"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db/clickhouse"

	"github.com/matrosov-nikita/smart-generator/pkg/client/db/postgres"

	"github.com/matrosov-nikita/smart-generator/pkg/config"

	"github.com/matrosov-nikita/smart-generator/generator"
)

const (
	ClickhouseURL = "tcp://127.0.0.1:9000?debug=false"
	PostgresURL   = "postgresql://postgres@127.0.0.1:5432/generator?sslmode=disable"
)

func main() {
	cfg := &config.Config{}
	flag.IntVar(&cfg.ServersCount, "servers", 1, "Servers count")
	flag.IntVar(&cfg.TeamsCount, "teams", 2, "Teams count")
	flag.StringVar(&cfg.DetectorsConfigPath, "detectorsConfigPath", "./config.json", "Path to file stores amount events per detector")
	flag.StringVar(&cfg.StartDateStr, "startDate", "2020-01-01", "Start date for events generation")
	flag.StringVar(&cfg.EndDateStr, "endDate", "", "End date for events generation, default: now")
	flag.StringVar(&cfg.TimeGeneratorType, "generatorType", "normal", "Generate type for time seq for events: normal or random")
	flag.StringVar(&cfg.StorageType, "storageType", "all", "clickhouse, postgres, all")
	flag.Parse()

	if err := cfg.ParseFields(); err != nil {
		log.Fatal(err)
	}

	log.Printf("start with config: %+v", cfg)
	pgClient, err := postgres.NewClient(PostgresURL)
	if err != nil {
		log.Fatalf("failed to create postgres client: %+v", err)
	}

	chClient, err := clickhouse.NewClient(ClickhouseURL)
	if err != nil {
		log.Fatalf("failed to create clickhouse client: %+v", err)
	}
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
	case "all":
		runGeneratorWithClient("clickhouse", chClient, cfg)
		runGeneratorWithClient("postgres", pgClient, cfg)
	}
}

func runGeneratorWithClient(storageType string, client generator.Client, cfg *config.Config) {
	gen := generator.New(client, cfg.StartDate, cfg.EndDate, cfg.ServersCount, cfg.TeamsCount,
		cfg.TimeGeneratorType, cfg.Detectors)
	startTime := time.Now()
	gen.Run()
	timeElapsed := time.Since(startTime)
	log.Printf("Loading to %s successfully completed, time elapsed %.3f sec.", storageType, timeElapsed.Seconds())
}
