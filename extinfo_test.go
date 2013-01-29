package extinfo

import (
	"testing"
)

func TestGetBasicInfo(t *testing.T) {
	psl1 := NewServer("sauerleague.org", 10000)
	_, err := psl1.GetBasicInfo()
	if err != nil {
		t.Fail()
	}
}

func TestGetUptime(t *testing.T) {
	psl1 := NewServer("sauerleague.org", 10000)
	_, err := psl1.GetUptime()
	if err != nil {
		t.Fail()
	}
}

func TestGetPlayerInfo(t *testing.T) {
	psl1 := NewServer("sauerleague.org", 10000)
	_, err := psl1.GetPlayerInfo(2)
	if err != nil {
		t.Fail()
	}
}

func TestGetTeamsScores(t *testing.T) {
	psl1 := NewServer("sauerleague.org", 10000)
	_, err := psl1.GetTeamsScores()
	if err != nil {
		t.Fail()
	}
}
