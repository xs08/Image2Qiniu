package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	defineErrors "image2qiniu/errors"
)

// AppKey qiniu appkey
type AppKey struct {
	AccessKey string `yaml:"AccessKey"`
	SecretKey string `yaml:"SecretKey"`
}

// Bucket 需要上传的 Bucket 名字
type Bucket struct {
	Name      string `yaml:"Name"`
	Domin     string `yaml:"Domin"`
	KeyPerfix string `yaml:"KeyPerfix"`
}

// AppConfig qiniu app
type AppConfig struct {
	AppKey AppKey `yaml:"AppKey"`
	Bucket Bucket `yaml:"Bucket"`
}

// LoadConfig load config from configPath
func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		return nil, nil
	}

	_, err := os.Stat(configPath)
	if err != nil {
		return nil, defineErrors.ErrConfigFileNotExits
	}

	var appConfig AppConfig

	// Read config file
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, defineErrors.ErrOpenConfig
	}

	err = yaml.Unmarshal(configFile, &appConfig)
	if err != nil {
		return nil, defineErrors.ErrLoadConfig
	}

	return &appConfig, nil
}
