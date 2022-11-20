package s3

import (
	"log"
	"os"
	"s3telegram/client"
	"s3telegram/config"

	"github.com/aws/aws-sdk-go/aws"
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

func UploadToS3(file_desc *config.FileDesc) (string, error) {
	c := config.GetConfig()

	bucket := aws.String(c.S3Bucket)
	key := aws.String(file_desc.FileID + getExt(file_desc.FilePath))
	log.Printf("accesskey [%s] secretkey [%s] bucket [%s]", c.S3AccessKey, c.S3SecretKey, c.S3Bucket)

	content_type := aws.String("image/png")
	acl := aws.String("private")
	metadata_key := "udf-metadata"
	metadata_value := file_desc.FileID
	metadata := map[string]*string{
		metadata_key: &metadata_value,
	}

	uploader := client.GetS3Uploader()
	filename := file_desc.FilePath
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("open file fail, %s", err.Error())
		return "", err
	}
	defer f.Close()

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      bucket,
		Key:         key,
		Body:        f,
		ContentType: content_type,
		ACL:         acl,
		Metadata:    metadata,
	}, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
		u.LeavePartsOnError = true
		u.Concurrency = 3
	})

	if err != nil {
		log.Printf("upload to s3 fail, %s", err.Error())
		return "", err
	}

	log.Printf("file uploaded to, %s", result.Location)
	return result.Location, nil
}
