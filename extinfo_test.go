package extinfo

import (
	"net"
	"testing"
	"time"
)

var testAddress *net.UDPAddr

func init() {
	var err error
	testAddress, err = net.ResolveUDPAddr("udp", "sauerleague.org:10000")
	if err != nil {
		panic(err)
	}
}

func TestGetBasicInfo(t *testing.T) {
	psl1 := NewServer(testAddress, 5*time.Second)
	_, err := psl1.GetBasicInfo()
	if err != nil {
		t.Fail()
	}
}

func TestGetUptime(t *testing.T) {
	psl1 := NewServer(testAddress, 5*time.Second)
	_, err := psl1.GetUptime()
	if err != nil {
		t.Fail()
	}
}

func TestGetPlayerInfo(t *testing.T) {
	psl1 := NewServer(testAddress, 5*time.Second)
	_, err := psl1.GetPlayerInfo(2)
	if err != nil {
		t.Fail()
	}
}

func TestGetAllPlayerInfo(t *testing.T) {
	psl1 := NewServer(testAddress, 5*time.Second)
	_, err := psl1.GetAllPlayerInfo()
	if err != nil {
		t.Fail()
	}
}

func TestGetTeamsScores(t *testing.T) {
	psl1 := NewServer(testAddress, 5*time.Second)
	_, err := psl1.GetTeamsScores()
	if err != nil {
		t.Fail()
	}
}

func testGetMasterServerList(t *testing.T) {
    _, err := GetMasterServerList(5*time.Second)
    if err != nil {
        t.Fail()
    }
}
