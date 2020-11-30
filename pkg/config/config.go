package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	ServersCount        int
	TeamsCount          int
	BatchSize           int
	DetectorsConfigPath string
	StartDateStr        string
	StartDate           time.Time
	EndDateStr          string
	EndDate             time.Time
	TimeGeneratorType   string
	StorageType         string
	Detectors           map[string]int
}

func (c *Config) ParseFields() error {
	startDate, err := time.Parse("2006-01-02", c.StartDateStr)
	if err != nil {
		return errors.New("start date must be at YYYY-MM-DD format")
	}
	c.StartDate = startDate

	c.EndDate = time.Now()
	if c.EndDateStr != "" {
		endDate, err := time.Parse("2006-01-02", c.EndDateStr)
		if err != nil {
			return errors.New("end date must be at YYYY-MM-DD format")
		}
		c.EndDate = endDate
	}

	detectors, err := readDetectorsConfigFile(c.DetectorsConfigPath)
	if err != nil {
		return fmt.Errorf("failed to parse detectors config file: %+v", err)
	}

	c.Detectors = detectors
	return nil
}

func readDetectorsConfigFile(detectorsConfigPath string) (map[string]int, error) {
	detectorsConfigFile, err := os.Open(detectorsConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config.json file: %+v", err)
	}

	defer func() {
		if err := detectorsConfigFile.Close(); err != nil {
			log.Printf("failed to close config file: %v", err)
		}
	}()

	var detectorData map[string]int
	if err := json.NewDecoder(detectorsConfigFile).Decode(&detectorData); err != nil {
		return nil, fmt.Errorf("failed to decode config data: %+v", err)
	}

	return detectorData, nil
}
