package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type BotConfig struct {
	BotToken       string `validate:"required"`
	BotTimeout     int
	S3AccessKey    string `validate:"required"`
	S3SecretKey    string `validate:"required"`
	S3Bucket       string `validate:"required"`
	S3Host         string `validate:"required"`
	CloudFrontHost string `validate:"default=false"`
	TmpPath        string `validate:"default=false"`
}

var configFile string = ""

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

func getIniConfig() *BotConfig {
	if Config != nil && getLastModifyTime(configFile).Sub(lastModifyTime) <= 0 {
		log.Printf("already read config file")
		return Config
	}

	log.Println("read config file")
	file, err := ini.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config = &BotConfig{}

	config.BotToken = file.Section("Bot").Key("token").String()
	timeout, err := file.Section("Bot").Key("timeout").Int()
	if err == nil {
		config.BotTimeout = timeout
	}
	config.TmpPath = file.Section("Bot").Key("tmppath").String()
	config.S3AccessKey = file.Section("S3").Key("s3accesskey").String()
	config.S3SecretKey = file.Section("S3").Key("s3secretkey").String()
	config.S3Bucket = file.Section("S3").Key("s3bucket").String()
	config.S3Host = file.Section("S3").Key("s3host").String()
	config.CloudFrontHost = file.Section("S3").Key("cloudfronthost").String()

	lastModifyTime = getLastModifyTime(configFile)

	Config = config
	log.Printf("Config: %v\n", Config)
	return Config
}

func getEnvAsMap() map[string]string {
	rtn := make(map[string]string)
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			rtn[e[:i]] = e[i+1:]
		}
	}

	return rtn
}

func getEnvConfig(envMap map[string]string) *BotConfig {
	fmt.Println(envMap)
	config := &BotConfig{}
	config.BotToken = envMap["bottoken"]
	if timeoutstr, ok := envMap["bottimeout"]; ok {
		if timeout, err := strconv.Atoi(timeoutstr); err == nil {
			config.BotTimeout = timeout
		}
	}
	config.S3AccessKey = envMap["s3accesskey"]
	config.S3SecretKey = envMap["s3secretkey"]
	config.S3Bucket = envMap["s3bucket"]
	config.S3Host = envMap["s3host"]
	if cloudfronthost, ok := envMap["cloudfronthost"]; ok {
		config.CloudFrontHost = cloudfronthost
	}
	if tmppath, ok := envMap["tmppath"]; ok {
		config.TmpPath = tmppath
	} else {
		config.TmpPath = "."
	}

	Config = config
	log.Printf("Config: %v\n", Config)
	return Config
}

func GetConfig() *BotConfig {
	envMap := getEnvAsMap()
	if _, ok := envMap["config_file"]; ok {
		configFile = envMap["config_file"]
		return getIniConfig()
	}
	return getEnvConfig(envMap)
}

type FileDesc struct {
	FileID   string
	FilePath string
}
