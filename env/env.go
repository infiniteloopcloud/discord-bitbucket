package env

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	configuration *Static
	ConfigFlag    string = "BITBUCKET_BOT_CONFIG"
)

type Static struct {
	SkipRepoPushMessages string `json:"skip_reop_push_message"`
	BotToken             string `json:"bot_token"`
	BotGuild             string `json:"bot_guild"`
	Address              string `json:"address"`
}

func Configuration() *Static {
	if configuration == nil {
		var path string
		if path = os.Getenv(ConfigFlag); path == "" {
			path = "./config.json"
		}
		// read from path
		var s Static
		// json unmarshal into s
		configuration = &s

		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("[ERROR] %s", err.Error())
		}

		json.Unmarshal(file, &configuration)

	}
	return configuration
}