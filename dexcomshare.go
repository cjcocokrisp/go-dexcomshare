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

// Enum for region types.
type Region int

type RegionInfo struct {
	BaseURL       string
	ApplicationID string
}

// Possible regions
const (
	RegionUS Region = iota
	RegionOUS
	RegionJP
)

// Map von Region → RegionInfo
var regionInfos = map[Region]RegionInfo{
	RegionUS: {
		BaseURL:       BaseUrlUS,
		ApplicationID: ApplicationIdUS,
	},
	RegionOUS: {
		BaseURL:       BaseUrlOUS,
		ApplicationID: ApplicationIdOUS,
	},
	RegionJP: {
		BaseURL:       BaseUrlJP,
		ApplicationID: ApplicationIdJP,
	},
}

func (r Region) BaseUrl() string {
	return regionInfos[r].BaseURL
}

func (r Region) AppID() string {
	return regionInfos[r].ApplicationID
}

// Function to make a POST request to the Dexcom API.
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

// Log into Dexcom with your username and password.
func Login(username string, password string, region Region) (*DexcomSession, error) {
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

	if needToAuth {
		body, err := json.Marshal(map[string]string{
			"accountName":   username,
			"password":      password,
			"applicationId": region.AppID(),
		})
		if err != nil {
			return nil, err
		}
		res, err := DexcomAPIRequest(region.BaseUrl()+AuthPath, body)
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
		"applicationId": region.AppID(),
		"password":      password,
	})

	if err != nil {
		return nil, err
	}
	res, err := DexcomAPIRequest(region.BaseUrl()+LoginPathId, body)
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
		BaseUrl:   region.BaseUrl(),
	}, nil
}

// Get EGV values from CGM.
// @arg - amount (amount of readings to get)
// @arg - minutes (not really sure what this is documentation was not very specific, it uses 1440 so just stick with that)
func (dexcom DexcomSession) GetEGV(amount int, minutes int) ([]EstimatedGlucoseValue, error) {
	if dexcom.sessionid == nil {
		return nil, errors.New("Invalid Session Token.")
	}

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
