package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type BotConfig struct {
	Token          string
	TmpPath        string
	BotTimeout     int
	S3AccessKey    string
	S3SecretKey    string
	S3Bucket       string
	S3Host         string
	CloudFrontHost string
	RedisAddr      string
}

const configFile = "bot.ini"

var Config *BotConfig = nil
var lastModifyTime = time.Now()

func getLastModifyTime(filename string) time.Time {
	file, err := os.Stat(filename)
	if err != nil {
		log.Printf("read file err %s", err.Error())
		return time.Now()
	}

	return file.ModTime()
}

func GetConfig() *BotConfig {
	if Config != nil && getLastModifyTime(configFile).Sub(lastModifyTime) <= 0 {
		log.Printf("already read config file")
		return Config
	}

	log.Println("read config file")
	file, err := ini.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config = new(BotConfig)

	config.Token = file.Section("Bot").Key("token").String()
	config.TmpPath = file.Section("Bot").Key("tmppath").String()
	timeout, err := file.Section("Bot").Key("timeout").Int()
	if err == nil {
		config.BotTimeout = timeout
	}
	config.S3AccessKey = file.Section("S3").Key("s3accesskey").String()
	config.S3SecretKey = file.Section("S3").Key("s3secretkey").String()
	config.S3Bucket = file.Section("S3").Key("s3bucket").String()
	config.S3Host = file.Section("S3").Key("s3host").String()
	config.CloudFrontHost = file.Section("S3").Key("cloudfronthost").String()
	config.RedisAddr = file.Section("Redis").Key("addr").String()

	lastModifyTime = getLastModifyTime(configFile)

	Config = config
	fmt.Printf("Config: %v\n", Config)
	return Config
}

type FileDesc struct {
	FileID   string
	FilePath string
}
