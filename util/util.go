package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"s3telegram/config"
	"s3telegram/s3"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsRawPhoto(m *tgbotapi.Message) bool {
	return m.Document != nil
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func RandBase64() string {
	size := 10 // change the length of the generated random string here

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)
	return rs
}

// download file to path
func downloadFile(url string, full_path string) error {
	file, err := os.Create(full_path)
	if err != nil {
		log.Printf("Create file fail, [%s]", err.Error())
		return err
	}
	defer file.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Download and save file size %d", size)

	return nil
}

func ProcessUploadS3(doc *tgbotapi.Document, bot *tgbotapi.BotAPI) (string, error) {
	conf := config.GetConfig()
	file_conf := tgbotapi.FileConfig{
		FileID: doc.FileID,
	}
	file_conf.FileID = doc.FileID
	file, err := bot.GetFile(file_conf)
	if err != nil {
		return "", err
	}

	url := file.Link(bot.Token)
	log.Printf("url [%s]", url)
	full_path := filepath.Join(conf.TmpPath, doc.FileID, file.FilePath)
	dir_path := filepath.Dir(full_path)
	if !pathExist(dir_path) {
		os.MkdirAll(dir_path, os.ModePerm)
	}

	err2 := downloadFile(url, full_path)
	if err2 != nil {
		log.Fatal("Download file fail")
		return "", err2
	}
	log.Printf("Downlaod file to [%s]", full_path)

	file_desc := &config.FileDesc{
		FileID:   doc.FileID,
		FilePath: full_path,
	}
	s3_location, err := s3.UploadToS3(file_desc)
	if err != nil {
		log.Printf("upload to s3 fail, %s", err)
		return "", err
	}
	log.Println(s3_location)
	return s3_location, nil
}

func TmpFileClear() {
	for range time.Tick(time.Hour) {
		c := config.GetConfig()
		filepath.Walk(c.TmpPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && time.Since(info.ModTime()) > 24*time.Hour {
				os.RemoveAll(path)
				return nil
			}
			return nil
		})
	}
}
