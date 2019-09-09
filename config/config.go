package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	defineError "image2qiniu/errors"
)

// AppKey qiniu appkey
type AppKey struct {
	AccessKey string `yaml:"AccessKey"`
	SecretKey string `yaml:"SecretKey"`
}

// Bucket 需要上传的 Bucket 名字
type Bucket struct {
	Name  string `yaml:"Name"`
	Domin string `yaml:"Domin"`
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

	var appConfig AppConfig

	// Read config file
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, defineError.ErrOpenConfig
	}

	err = yaml.Unmarshal(configFile, &appConfig)
	if err != nil {
		return nil, defineError.ErrLoadConfig
	}

	return &appConfig, nil
}
