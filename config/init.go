package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/froggy-12/mooshroombase/utils"
)

func InitConfigs() Config {
	var configs Config
	if _, err := os.Stat("configs.json"); os.IsNotExist(err) {
		CreateDefaultConfig(&configs)
	} else {
		data, err := os.ReadFile("configs.json")
		if err != nil {
			log.Fatal("Error: ", err.Error())
		}
		err = json.Unmarshal(data, &configs)
		if err != nil {
			log.Fatal("Error: ", err.Error())
		}
	}
	return configs
}

func CreateDefaultConfig(configs *Config) {
	*configs = Config{
		PrimaryDB:                 "mongodb",
		RunningDatabaseContainers: []string{"mongodb", "redis", "mariadb"},
		MongoDBUsername:           "root",
		MongoDBPassword:           "password",
		MariaDBRootPassword:       "password",
		JWTSecret:                 "SuchASuperSecretMooshroom",
		JWTTokenExpiration:        time.Now().Add(time.Hour * 24 * 7),
		ChatFunctions:             true,
		Authentication:            true,
		BodySizeLimit:             100 * 1024 * 1024,
		SMTPServerAdress:          "smtp.gmail.com",
		SMTPServerPort:            "587",
		SMTPEmailFrom:             "mohi6644123@gmail.com",
		SMTPPassword:              "aahzizpbpbgiwtju",
		AppUrl:                    "http://localhost:6644",
		VerifyEmailRouteClient:    "/verify",
		EmailVerificationRoute:    "/verify-email",
		ClientAPICredintials:      map[string]string{},
		AllowedCorsOrigin:         []string{"*"},
		CorsHeadersMaxAge:         int(time.Hour * 24 * 7),
		JWTCookieAge:              int(time.Hour * 24 * 7),
		EmailVerificationAllowed:  false,
		GithubKey:                 "",
		GithubSecret:              "",
		DiscordKey:                "",
		DiscordSecret:             "",
		FacebookKey:               "",
		FaceBookSecret:            "",
		GoogleKey:                 "",
		GoogleSecret:              "",
		MicrosoftKey:              "",
		MicrosoftSecret:           "",
		LinkedInKey:               "",
		LinkedInSecret:            "",
		TwitterKey:                "",
		TwitterSecret:             "",
		AppleKey:                  "",
		AppleSecret:               "",
	}
	data, err := json.MarshalIndent(*configs, "", "  ")
	if err != nil {
		fmt.Println("Error creating default config:", err)
		os.Exit(1)
	}
	err = os.WriteFile("configs.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing default config:", err)
		os.Exit(1)
	}
	utils.DebugLogger("Important", "Default config file created. Please configure the settings in configs.json to avoid using default container names")
	utils.DebugLogger("Important", "Using default container names can lead to conflicts and security issues.")
	utils.DebugLogger("Important", "Please restart the application after configuring the settings.")
	utils.DebugLogger("Important", "Ohh yeah remember if the containers has been created u need to delete them manually")
	os.Exit(0)
}
