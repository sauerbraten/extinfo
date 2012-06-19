package extinfo

import (
	"testing"
)


func TestBasicInfo(t *testing.T) {
	_, err := GetBasicInfo("sauerleague.org", 10000)
	if err != nil {
		t.Fail()
	}
}

func TestUptime(t *testing.T) {
	_, err := GetUptime("sauerleague.org", 10000)
	if err != nil {
		t.Fail()
	}
}

func TestPlayerInfo(t *testing.T) {
	_, err := GetPlayerInfo("sauerleague.org", 10000, 14)
	if err != nil {
		t.Fail()
	}
}
