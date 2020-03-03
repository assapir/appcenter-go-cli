package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type config struct {
	OwnerName string `json:"owner"`
	AppName   string `json:"app"`
	APIKey    string `json:"apiKey"`
}

func getConfig() (*config, error) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var jsonConfig config
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return nil, err
	}

	return &jsonConfig, nil
}

func post(addr string, body interface{}, apiToken string, jsonBody bool) ([]byte, error) {
	var bodyValue []byte
	if jsonBody {
		bodyValue, _ = json.Marshal(body)
	}

	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(bodyValue))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Token", apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("Error while calling %s. Status: %s, body: %s", addr, resp.Status, string(bodyBytes))
	}

	log.Println("response status:", resp.Status)

	return bodyBytes, nil
}

func createUpload(config *config) (uploadID string, uploadURL string, err error) {
	const baseURL = "https://appcenter.ms/api/v0.1/apps/%s/%s/release_uploads"
	createUploadURL := fmt.Sprintf(baseURL, config.OwnerName, config.AppName)
	log.Println(createUploadURL)

	var body struct{}
	jsonBytes, err := post(createUploadURL, &body, config.APIKey, true)
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

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("error: no command line arguments supplied")
	}

	log.Println("got with this args:", args)

	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	_, _, err2 := createUpload(config)
	if err2 != nil {
		log.Fatal(err)
	}
}
