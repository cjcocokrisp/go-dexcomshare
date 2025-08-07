package dexcomshare

import (
	"testing"
)

func TestGetLatestEGV(t *testing.T) {
	dexcom, err := Login("USERNAME", "PASSWORD", "REGION")
	if err != nil {
		t.Error(err)
	}

	_, err = dexcom.GetLatestEGV()
	if err != nil {
		t.Error(err)
	}
}

func TestGetEGV(t *testing.T) {
	dexcom, err := Login("USERNAME", "PASSWORD")
	if err != nil {
		t.Error(err)
	}

	egvs, err := dexcom.GetEGV(5, DefaultMinutes)
	if err != nil {
		t.Error(err)
	}

	if len(egvs) != 5 {
		t.Error(len(egvs))
	}
}
