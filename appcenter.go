package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/assapir/appcenter-go-cli/config"
	"github.com/assapir/appcenter-go-cli/httpclient"
)

func createUpload(config *config.Config) (uploadID string, uploadURL string, err error) {
	const baseURL = "https://appcenter.ms/api/v0.1/apps/%s/%s/release_uploads"
	createUploadURL := fmt.Sprintf(baseURL, config.OwnerName, config.AppName)
	log.Println(createUploadURL)

	var body struct{}
	jsonBytes, err := httpclient.PostJSON(createUploadURL, &body, config.APIKey)
	if err != nil {
		return
	}

	var result map[string]string
	json.Unmarshal(jsonBytes, &result)

	uploadID = result["upload_id"]
	uploadURL = result["upload_url"]
	log.Printf("upload id: %s, upload url: %s", uploadID, uploadURL)
	return
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("error: no command line arguments supplied")
	}
	arg := os.Args[1]
	if len(arg) == 0 {
		log.Fatal("error: no command line arguments supplied")
	}

	log.Println("got with this arg:", arg)

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	_, uploadURL, err := createUpload(config)
	if err != nil {
		log.Fatal(err)
	}

	_, err2 := httpclient.PostFile(uploadURL, arg)
	if err2 != nil {
		log.Fatal(err2)
	}
}
