package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var ConfigFile = "config.json"

type Config struct {
	CountryCode    string   `json:"country_code"`
	FilePath       string   `json:"file_path"`
	Interface      string   `json:"interface"`
	IgnoredSubnets []string `json:"ignored_subnets"`
	IgnoredIPs     []string `json:"ignored_ips"`
}

// Функция для загрузки конфигурационного файла
func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурационного файла: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора конфигурационного файла: %v", err)
	}

	return &config, nil
}
