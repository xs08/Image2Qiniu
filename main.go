package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	flag "github.com/spf13/pflag"

	defineConfig "image2qiniu/config"
	defineErrors "image2qiniu/errors"
	"image2qiniu/utils"
)

// tmpImageStorePath temporary file store path
const tmpImageStorePath = "/tmp/image2qiniu/"

// default config filePath on user home dir
const defaultConfigFilePath = ".config/image4qiniu.yaml"

var (
	link       string // image URI
	image      string // local image path
	download   string // download image. store to desktop
	secretKey  string // upload sk
	accessKey  string // upload sk
	bucketName string // bucket for store image
	name       string // rename image
	keyPerfix  string // key perfix for store image
	nameSuffix string // Add suffix to name
	config     string // Config file path
)

func init() {
	flag.StringVarP(&link, "link", "l", "", "image link on net")
	flag.StringVarP(&image, "image", "i", "", "local image path")
	flag.StringVarP(&download, "download", "D", "", "download image")

	flag.StringVar(&secretKey, "sk", "", "Secret Key")
	flag.StringVar(&accessKey, "ak", "", "Access key")
	flag.StringVar(&bucketName, "bucketName", "", "image store bucket name")
	flag.StringVar(&name, "name", "", "rename image")
	flag.StringVar(&keyPerfix, "keyPrefix", "", "upload key prefix")
	flag.StringVar(&nameSuffix, "nameSuffix", "", "add suffix to name")
	flag.StringVar(&config, "config", "", "config file path")

	flag.Parse()
}

// parse config
func parseConfig() bool {
	// download image and upload, should check AssessKey and SecretKey

	if link == "" && image == "" && download == "" {
		log.Fatal(defineErrors.ErrNoImageSpecify)
		return false
	}

	return true
}

// downloa tasks
func startDownloakTasks(saveFilePath, link string, ok chan string) {
	file, err := os.Create(saveFilePath)
	if err != nil {
		log.Println("Create file error")
		log.Fatal(err)
	}
	defer file.Close()

	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	pix, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(file, bytes.NewReader(pix))
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	ok <- saveFilePath
}

// upload task
func startUpload(localFile, fileName, bucketName string, ok chan struct{}) {

	// 上传策略
	putPolicy := storage.PutPolicy{
		Scope: utils.JoinStrs(bucketName, ":", fileName),
	}

	// 认证消息
	mac := qbox.NewMac(accessKey, secretKey)
	// 策略与认证信息生成Token
	upToken := putPolicy.UploadToken(mac)

	// 存储配置
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// set proxy
	// proxyURL := "http://localhost:8888"
	// proxyURI, _ := url.Parse(proxyURL)

	//绑定网卡
	// nicIP := "192.168.0.110"
	dialer := &net.Dialer{
		// LocalAddr: &net.TCPAddr{
		// 	IP: net.ParseIP(nicIP),
		// },
	}

	//构建代理client对象
	client := http.Client{
		Transport: &http.Transport{
			// Proxy: http.ProxyURL(proxyURI),
			Dial: dialer.Dial,
		},
	}

	// 构建表单上传的对象
	formUploader := storage.NewFormUploaderEx(&cfg, &storage.Client{Client: &client})
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "wallpaper.png",
		},
	}
	//putExtra.NoCrc32Check = true
	err := formUploader.PutFile(context.Background(), &ret, upToken, fileName, localFile, &putExtra)

	if err != nil {
		log.Println("Upload error")
		log.Fatal(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
	fmt.Printf("%v\n", ret)
	ok <- struct{}{}
}

func main() {
	var appConfig *defineConfig.AppConfig
	var err error

	// config file path
	var configFilePath string

	// Read user home dir .config/image4qiniu.yaml
	if config == "" {
		// Get home dir and config file path
		user, _ := user.Current()
		configFilePath = filepath.Join(user.HomeDir, defaultConfigFilePath)
	} else {
		// useSpecify Path
		configFilePath = config
	}

	// Load config
	if config != "" {
		appConfig, err = defineConfig.LoadConfig(configFilePath)
		if err != nil {
			log.Print(err)
			// specify config file. but can't find it
			if err == defineErrors.ErrConfigFileNotExits {
				log.Fatal(err)
				return
			}
		}
	} else {
		// if config not exist, should initial appConfig
		appConfig = &defineConfig.AppConfig{}
	}
	// command line args first
	// checkout accessKey
	if accessKey != "" {
		appConfig.AppKey.AccessKey = accessKey
	} else {
		if appConfig.AppKey.AccessKey == "" {
			log.Fatal(defineErrors.ErrNoAccessKey)
			return
		}
	}
	// checkout secretKey
	if secretKey != "" {
		appConfig.AppKey.SecretKey = secretKey
	} else {
		if appConfig.AppKey.SecretKey == "" {
			log.Fatal(defineErrors.ErrNoSecretKey)
			return
		}
	}

	// checkout bucketName
	if bucketName != "" {
		appConfig.Bucket.Name = bucketName
	} else {
		if appConfig.Bucket.Name == "" {
			log.Fatal(defineErrors.ErrNoBucketName)
			return
		}
	}

	// fileName handler
	var fileName string
	if name != "" {
		fileName = name
	} else if link != "" {
		_, err = url.Parse(link)
		if err != nil {
			log.Fatal(err)
			log.Fatal(defineErrors.ErrLinkIsNotOk)
			return
		}
		// Download and upload
		// split link
		temPathArr := strings.SplitAfter(link, "/")
		// get last filename
		fileName = temPathArr[len(temPathArr)-1]
	}

	if fileName == "" {
		log.Fatal("file name not exists")
		return
	}
	// Download file to saveFilePath
	saveFilePath := filepath.Join(tmpImageStorePath, fileName)

	// key perfix
	if keyPerfix != "" {
		fileName = utils.JoinStrs(keyPerfix, fileName)
	} else if appConfig.Bucket.KeyPerfix != "" {
		fileName = utils.JoinStrs(appConfig.Bucket.KeyPerfix, fileName)
	}

	// name suffix
	if nameSuffix != "" {
		fileName = utils.JoinStrs(fileName, nameSuffix)
	}

	// template folder for store image
	if _, err = os.Stat(tmpImageStorePath); os.IsNotExist(err) {
		os.Mkdir(tmpImageStorePath, os.ModePerm)
	}

	downloadChan := make(chan string)
	uploadChan := make(chan struct{})

	// Download file
	go startDownloakTasks(saveFilePath, link, downloadChan)

	for {
		select {
		// download ok, start upload
		case storeFilePath := <-downloadChan:
			fmt.Println("Downloadok")
			go startUpload(storeFilePath, fileName, bucketName, uploadChan)
		case <-uploadChan:
			fmt.Println("upload success")
			os.Exit(0)
		}
	}
}
