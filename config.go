package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	BotName  string `json:"name"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func loadConfig() Config {
	file, err := os.Open("init.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode config JSON: %v", err)
	}
	return config
}
