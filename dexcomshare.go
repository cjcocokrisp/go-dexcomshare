package dexcomshare

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type DexcomSession struct {
	username  string
	password  string
	sessionid uuid.UUID
}

type EstimatedGlucoseValue struct {
	WT         string `json:"WT"`
	ST         string `json:"ST"`
	DT         string `json:"DT"`
	Value      int    `json:"Value"`
	Trend      string `json:"Trend"`
	TrendArrow string
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
	id := strings.ReplaceAll(string(body), "\"", "")
	sessionUuid, err := uuid.FromString(id)
	if err != nil {
		log.Fatal(err)
	}

	return &DexcomSession{
		username:  username,
		password:  password,
		sessionid: sessionUuid,
	}
}

func (dexcom DexcomSession) GetEGV(params ...int) []EstimatedGlucoseValue {
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

	url, err := url.Parse(BaseUrlUS + CurrentEGVPath)
	if err != nil {
		log.Fatal(err)
	}

	q := url.Query()
	q.Set("sessionId", dexcom.sessionid.String())
	q.Set("maxCount", strconv.Itoa(maxCount))
	q.Set("minutes", strconv.Itoa(minutes))
	url.RawQuery = q.Encode()

	res, err := http.Post(url.String(), "", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var egvs []EstimatedGlucoseValue
	json.Unmarshal(data, &egvs)
	for i, egv := range egvs {
		egvs[i].TrendArrow = TrendArrowMap[egv.Trend]
	}

	return egvs
}

func (dexcom DexcomSession) GetLatestEGV() EstimatedGlucoseValue {
	egvs := dexcom.GetEGV()
	return egvs[0]
}
