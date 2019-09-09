package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"

	"context"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var (
	// accessKey = os.Getenv("QINIU_ACCESS_KEY")
	// secretKey = os.Getenv("QINIU_SECRET_KEY")
	// bucket    = os.Getenv("QINIU_TEST_BUCKET")
	accessKey = "-hk4KJ9Ph_yIccAPqMwI39iyP__j4KvBYUJ3j7im"
	secretKey = "cEUI_OFXlbj6TMCCujsqPVGPwKX6XLfdbBgh5kvR"
	bucket    = "imagine-space"
)

func joinStrs(strs ...string) string {
	var buf bytes.Buffer
	for _, str := range strs {
		buf.WriteString(str)
	}
	return buf.String()
}

func main() {

	// 本地需要上传的文件
	localFile := "wallpaper.png"
	// 文件保存的key
	key := "images/wallpaper2.png"

	// 上传策略
	putPolicy := storage.PutPolicy{
		Scope: joinStrs(bucket, ":", key),
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

	//设置代理
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
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
	fmt.Printf("%v", ret)
}
