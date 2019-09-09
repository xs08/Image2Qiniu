package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
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

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(user.HomeDir)
	var appConfig AppConfig

	// Read config file
	configFile, err := ioutil.ReadFile("./image4qiniu.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(configFile, &appConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(appConfig)
}
