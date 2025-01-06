package dexcomshare

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type DexcomSession struct {
	username  string
	password  string
	sessionid string
}

func Login(username string, password string) *DexcomSession {
	body, err := json.Marshal(map[string]string{
		"accountName":   username,
		"password":      password,
		"applicationId": ApplicationId,
	})

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", BaseUrlUS+LoginPath, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == 500 {
		log.Fatal("AuthError: Invalid Dexcom username and password!\n")
	}

	body, err = io.ReadAll(res.Body)
	return &DexcomSession{
		username:  username,
		password:  password,
		sessionid: string(body),
	}
}
