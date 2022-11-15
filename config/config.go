package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type BotConfig struct {
	Token          string
	TmpPath        string
	S3AccessKey    string
	S3SecretKey    string
	S3Bucket       string
	S3Host         string
	CloudFrontHost string
}

type FileDesc struct {
	FileID   string
	FilePath string
}

func Read(r *BotConfig) error {
	cfg, err := ini.Load("bot.ini")
	if err != nil {
		log.Printf("Read config fail")
		return err
	}

	r.Token = cfg.Section("Bot").Key("token").String()
	r.TmpPath = cfg.Section("Bot").Key("tmppath").String()
	r.S3AccessKey = cfg.Section("S3").Key("s3accesskey").String()
	r.S3SecretKey = cfg.Section("S3").Key("s3secretkey").String()
	r.S3Bucket = cfg.Section("S3").Key("s3bucket").String()
	r.S3Host = cfg.Section("S3").Key("s3host").String()
	r.CloudFrontHost = cfg.Section("S3").Key("cloudfronthost").String()
	return nil
}
