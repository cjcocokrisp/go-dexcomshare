package dexcomshare

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Struct that represents a Dexcom session.
type DexcomSession struct {
	username  string
	password  string
	accountId string
	sessionid *string
	BaseUrl   string
}

// Struct that represents the parts of an EGV reading.
type EstimatedGlucoseValue struct {
	WT         string `json:"WT"`
	ST         string `json:"ST"`
	DT         string `json:"DT"`
	Value      int    `json:"Value"`
	Trend      string `json:"Trend"`
	TrendArrow string
}

// Log into Dexcom with your username and password.
func Login(username string, password string, region string) (*DexcomSession, error) {
	var needToAuth bool
	var accountId string

	if IsEmail(username) {
		// conitnue auth
		needToAuth = true
	} else if IsUUID(username) {
		// continue login with id
		needToAuth = false
		accountId = username
	} else {
		return nil, errors.New("invalid username format. use email or uuid.")
	}

	var BaseUrl string
	switch region {
	case "us":
		BaseUrl = BaseUrlUS
	case "ous":
		BaseUrl = BaseUrlOUS
	default:
		return nil, errors.New("invalid region specified, use 'us' or 'ous'")
	}

	if needToAuth {
		body, err := json.Marshal(map[string]string{
			"accountName":   username,
			"password":      password,
			"applicationId": ApplicationId,
		})
		if err != nil {
			return nil, err
		}
		res, err := DexcomAPIRequest(BaseUrl+AuthPath, body)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		accountId = strings.ReplaceAll(string(resBody), "\"", "")

	}

	body, err := json.Marshal(map[string]string{
		"accountId":     accountId,
		"applicationId": ApplicationId,
		"password":      password,
	})

	if err != nil {
		return nil, err
	}

	res, err := DexcomAPIRequest(BaseUrl+LoginPathId, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	id := strings.ReplaceAll(string(resBody), "\"", "")
	if err != nil {
		return nil, err
	}

	return &DexcomSession{
		username:  username,
		password:  password,
		accountId: accountId,
		sessionid: &id,
		BaseUrl:   BaseUrl,
	}, nil
}

// Get EGV values from CGM.
// @arg - amount (amount of readings to get)
// @arg - minutes (not really sure what this is documentation was not very specific, it uses 1440 so just stick with that)
func (dexcom DexcomSession) GetEGV(amount int, minutes int) ([]EstimatedGlucoseValue, error) {
	if dexcom.sessionid == nil {
		return nil, errors.New("Invalid Session Token.")
	}

	//TODO: Remove Debugging prints
	//fmt.Printf("%s\n", dexcom.BaseUrl+CurrentEGVPath)
	//fmt.Printf("%s\n", *dexcom.sessionid)
	url, err := url.Parse(dexcom.BaseUrl + CurrentEGVPath)
	if err != nil {
		log.Fatal(err)
	}

	q := url.Query()
	q.Set("sessionId", *dexcom.sessionid)
	q.Set("maxCount", strconv.Itoa(amount))
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

// Get latest EGV from CGM.
func (dexcom DexcomSession) GetLatestEGV() (*EstimatedGlucoseValue, error) {
	egvs, err := dexcom.GetEGV(1, DefaultMinutes)
	if len(egvs) == 0 {
		return nil, errors.New("ReadingError: No readings were available.")
	}
	return &egvs[0], err
}
