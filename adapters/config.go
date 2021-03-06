package adapters

import (
	"log"
	"os"

	"github.com/AlecSmith96/dnd-bot/entities"
	"gopkg.in/yaml.v2"
)

func GetConfig() entities.Config {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Unable to open config file: %v", err)
	}
	defer f.Close()

	var config entities.Config
	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Unable to decode config file: %v", err)
	}

	return config
}