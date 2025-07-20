package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DB_URL   string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	//Setup blank config struct
	config := Config{}

	filepath := getConfigFilepath()

	//Read File at filepath
	data, err := os.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	//Parse Json into Config Struct
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func SetUser(username string, configFile Config) error {
	//Set username
	configFile.UserName = username

	err := writeConfig(configFile)

	return err
}

func getConfigFilepath() string {
	filepath, _ := os.UserHomeDir()
	filepath += "/"
	filepath += configFileName
	return filepath
}

func writeConfig(configFile Config) error {
	//Convert Struct to Json
	jsonData, err := json.Marshal(configFile)
	if err != nil {
		return err
	}

	//fmt.Println(string(jsonData))

	//Write Json to file
	filepath := getConfigFilepath()
	os.WriteFile(filepath, jsonData, 0664)
	return nil
}
