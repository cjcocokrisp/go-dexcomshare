package dexcomshare

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DexcomSession struct {
	username  string
	password  string
	sessionid *string
}

type EstimatedGlucoseValue struct {
	WT         string `json:"WT"`
	ST         string `json:"ST"`
	DT         string `json:"DT"`
	Value      int    `json:"Value"`
	Trend      string `json:"Trend"`
	TrendArrow string
}

func Login(username string, password string) (*DexcomSession, error) {
	body, err := json.Marshal(map[string]string{
		"accountName":   username,
		"password":      password,
		"applicationId": ApplicationId,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", BaseUrlUS+LoginPath, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 500 {
		return nil, errors.New("AuthError: Invalid username or password!")
	}

	body, err = io.ReadAll(res.Body)
	id := strings.ReplaceAll(string(body), "\"", "")
	if err != nil {
		return nil, err
	}

	return &DexcomSession{
		username:  username,
		password:  password,
		sessionid: &id,
	}, nil
}

func (dexcom DexcomSession) GetEGV(params ...int) ([]EstimatedGlucoseValue, error) {
	var maxCount, minutes int
	if len(params) > 0 {
		maxCount = params[0]
		if len(params) > 1 {
			minutes = params[1]
		} else {
			minutes = 1440
		}
	} else {
		maxCount = 1
		minutes = 1440
	}

	if dexcom.sessionid == nil {
		return nil, errors.New("Invalid Session Token.")
	}

	url, err := url.Parse(BaseUrlUS + CurrentEGVPath)
	if err != nil {
		log.Fatal(err)
	}

	q := url.Query()
	q.Set("sessionId", dexcom.sessionid)
	q.Set("maxCount", strconv.Itoa(maxCount))
	q.Set("minutes", strconv.Itoa(minutes))
	url.RawQuery = q.Encode()

	res, err := http.Post(url.String(), "", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var egvs []EstimatedGlucoseValue
	json.Unmarshal(data, &egvs)
	for i, egv := range egvs {
		egvs[i].TrendArrow = TrendArrowMap[egv.Trend]
	}

	return egvs, nil
}

func (dexcom DexcomSession) GetLatestEGV() (EstimatedGlucoseValue, error) {
	egvs, err := dexcom.GetEGV()
	return egvs[0], err
}
