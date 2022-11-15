package s3

import (
	"log"
	"os"
	"s3telegram/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func getExt(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}

	return ""
}

func UploadToS3(c *config.BotConfig, file_desc *config.FileDesc) (string, error) {
	bucket := aws.String(c.S3Bucket)
	key := aws.String(file_desc.FileID + getExt(file_desc.FilePath))
	access_key := c.S3AccessKey
	secret_key := c.S3SecretKey
	log.Printf("accesskey [%s] secretkey [%s] bucket [%s]", c.S3AccessKey, c.S3SecretKey, c.S3Bucket)
	myContentType := aws.String("image/png")
	myACL := aws.String("private")
	metadata_key := "udf-metadata"
	metadata_value := file_desc.FileID
	myMetadata := map[string]*string{
		metadata_key: &metadata_value,
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(access_key, secret_key, ""),
		Region:           aws.String("ap-southeast-1"),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(false),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	cparams := &s3.HeadBucketInput{
		Bucket: bucket,
	}
	_, err := s3Client.HeadBucket(cparams)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	uploader := s3manager.NewUploader(newSession)
	filename := file_desc.FilePath
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      bucket,
		Key:         key,
		Body:        f,
		ContentType: myContentType,
		ACL:         myACL,
		Metadata:    myMetadata,
	}, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
		u.LeavePartsOnError = true
		u.Concurrency = 3
	})

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Printf("file uploaded to, %s", result.Location)
	return result.Location, nil
}
