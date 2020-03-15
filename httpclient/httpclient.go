package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// PostJSON - POST HTTP request with JSON body
func PostJSON(addr string, body interface{}, apiToken string) ([]byte, error) {
	var bodyValue []byte
	bodyValue, _ = json.Marshal(body)

	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(bodyValue))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Token", apiToken)

	return doPost(req, addr)
}

// PostFile - Upload a file from path to addr
func PostFile(addr string, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("ipa", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", addr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	return doPost(req, addr)
}

func doPost(req *http.Request, addr string) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
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
