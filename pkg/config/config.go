package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func LoadConfig(configPath string) (map[string]int, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config.json file: %+v", err)
	}

	defer func() {
		if err := configFile.Close(); err != nil {
			log.Printf("failed to close config file: %v", err)
		}
	}()

	var detectorData map[string]int
	if err := json.NewDecoder(configFile).Decode(&detectorData); err != nil {
		return nil, fmt.Errorf("failed to decode config data: %+v", err)
	}

	return detectorData, nil
}
