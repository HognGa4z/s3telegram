package util

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"s3telegram/config"
	"s3telegram/s3"

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

func DownloadPhoto(url string, file_id string, path string) error {

	full_path := path + file_id
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
	}
	log.Printf("Download and save file size %d", size)

	return nil
}

// download file to path
func DownloadFile(url string, full_path string) error {
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

func ProcessDocument(doc *tgbotapi.Document, bot *tgbotapi.BotAPI, conf *config.BotConfig) (string, error) {
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

	err2 := DownloadFile(url, full_path)
	if err2 != nil {
		log.Fatal("Download file fail")
		return "", err2
	}
	log.Printf("Downlaod file to [%s]", full_path)

	file_desc := &config.FileDesc{
		FileID:   doc.FileID,
		FilePath: full_path,
	}

	s3_location, err := s3.UploadToS3(conf, file_desc)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	log.Printf(s3_location)
	return s3_location, nil
}
