package extinfo

import (
	"log"
	"net"
	"testing"
	"time"
)

var srv *Server

func init() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:28785")
	if err != nil {
		panic(err)
	}
	srv, err = NewServer(*addr, 5*time.Second)
	if err != nil {
		panic(err)
	}
}

func TestGetBasicInfo(t *testing.T) {
	_, err := srv.GetBasicInfo()
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestGetUptime(t *testing.T) {
	_, err := srv.GetUptime()
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestGetClientInfo(t *testing.T) {
	_, err := srv.GetClientInfo(10)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestGetAllClientInfo(t *testing.T) {
	_, err := srv.GetAllClientInfo()
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestGetTeamScores(t *testing.T) {
	_, err := srv.GetTeamScores()
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}
