package extinfo

import (
	"testing"
)


func TestGetBasicInfo(t *testing.T) {
	_, err := GetBasicInfo("sauerleague.org", 10000)
	if err != nil {
		t.Fail()
	}
}

func TestGetUptime(t *testing.T) {
	_, err := GetUptime("sauerleague.org", 10000)
	if err != nil {
		t.Fail()
	}
}

func TestGetPlayerInfo(t *testing.T) {
	_, err := GetPlayerInfo("sauerleague.org", 10000, 14)
	if err != nil {
		t.Fail()
	}
}
