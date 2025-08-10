package dexcomshare

import (
	"bytes"
	"errors"
	"net/http"
	"regexp"
	"time"
)

func Post(url string, body []byte) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func IsEmail(inputStr string) bool {
	emailRegex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	return emailRegex.MatchString(inputStr)
}

func IsUUID(inputStr string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(inputStr)
}

func DexcomAPIRequest(url string, payload []byte) (*http.Response, error) {
	res, err := Post(url, payload)
	if err != nil {
		return nil, err
	}

	// Check if the response status is OK
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("API request failed with status: " + res.Status)
	}

	return res, nil
}
