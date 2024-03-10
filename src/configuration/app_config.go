package configuration

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	App struct {
		Port int `json:"port"`
	}
	Cache struct {
		MaxSize        int64  `json:"max_size"`
		PruneSize      uint32 `json:"prune_size"`
		ExpiryTimeSecs int    `json:"expiry_time_secs"`
	} `json:"cache"`
	ApiDocs struct {
		IsEnabled bool `json:"is_enabled"`
	} `json:"api_docs"`
	Pprof struct {
		IsEnabled bool `json:"is_enabled"`
		Port      int  `json:"port"`
	} `json:"pprof"`
	Metrics struct {
		Port int `json:"port"`
	} `json:"metrics"`
	Worker struct {
		BufferSize int `json:"buffer_size"`
		PoolSize   int `json:"pool_size"`
	} `json:"worker"`
}

var AppConfig Config

func LoadConfig(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = json.Unmarshal(file, &AppConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
}
