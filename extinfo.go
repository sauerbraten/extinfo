// Package extinfo provides easy access to the state information of a Sauerbraten game server (called 'extinfo' in the Sauerbraten source code).
package extinfo

import (
	"net"
	"strconv"
	"time"
)

// Protocol constants
const (
	// Constants describing the type of information to query for
	EXTENDED_INFO = 0x00 //0
	BASIC_INFO    = 0x01 //1

	EXTENDED_INFO_ACK      = 0xFF // -1
	EXTENDED_INFO_VERSION  = 105
	EXTENDED_INFO_ERROR    = 0x01 //1
	EXTENDED_INFO_NO_ERROR = 0x00 //0

	// Constants describing the type of extended information to query for
	EXTENDED_INFO_UPTIME       = 0x00 //0
	EXTENDED_INFO_CLIENT_INFO  = 0x01 //1
	EXTENDED_INFO_TEAMS_SCORES = 0x02 //2

	EXTENDED_INFO_CLIENT_INFO_RESPONSE_CNS  = 0xF6 //-10
	EXTENDED_INFO_CLIENT_INFO_RESPONSE_INFO = 0xF5 //-11
)

const MAX_PLAYER_CN = 127

// A server to query extinfo from.
type Server struct {
	addr    *net.UDPAddr
	timeOut time.Duration
}

func NewServer(host string, port int, timeOut time.Duration) (s *Server, err error) {
	var addr *net.UDPAddr
	addr, err = net.ResolveUDPAddr("udp", host+":"+strconv.Itoa(port+1))
	if err != nil {
		return
	}

	s = &Server{
		addr:    addr,
		timeOut: timeOut,
	}

	return
}
