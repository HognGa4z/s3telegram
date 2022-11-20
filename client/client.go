package client

import (
	"context"
	"log"
	"s3telegram/config"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-redis/redis/v8"
)

var mutex = sync.Mutex{}

var redisClient *redis.Client = nil

func GetRedisClient() (*redis.Client, error) {
	var ctx context.Context
	if redisClient.Ping(ctx) == nil {
		return redisClient, nil
	}
	c := config.GetConfig()

	mutex.Lock()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	mutex.Unlock()
	return redisClient, nil
}

var awsSession *session.Session = nil

func GetAwsSession() *session.Session {
	if awsSession != nil {
		return awsSession
	}

	c := config.GetConfig()

	access_key := c.S3AccessKey
	secret_key := c.S3SecretKey

	awsConfig := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(access_key, secret_key, ""),
		Region:           aws.String("ap-southeast-1"),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(false),
	}

	mutex.Lock()
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		log.Printf("create aws session fail, %s", err.Error())
		return nil
	}
	mutex.Unlock()
	return awsSession
}

var s3Client *s3.S3 = nil

func getS3Client(awsSession *session.Session) *s3.S3 {
	if s3Client != nil {
		return s3Client
	}

	c := config.GetConfig()
	bucket := aws.String(c.S3Bucket)

	params := &s3.HeadBucketInput{
		Bucket: bucket,
	}

	mutex.Lock()
	s3Client = s3.New(awsSession)
	_, err := s3Client.HeadBucket(params)
	if err != nil {
		log.Printf("create s3 client fail, %s", err.Error())
		return nil
	}
	mutex.Unlock()
	return s3Client
}

var s3Uploader *s3manager.Uploader = nil

func GetS3Uploader() *s3manager.Uploader {
	if s3Uploader != nil {
		return s3Uploader
	}

	awsSession := GetAwsSession()
	if awsSession == nil {
		panic("create aws session fail")
	}
	s3Client := getS3Client(awsSession)
	if s3Client == nil {
		panic("create aws s3 client fail")
	}

	mutex.Lock()
	s3Uploader = s3manager.NewUploader(awsSession)
	mutex.Unlock()
	return s3Uploader
}
