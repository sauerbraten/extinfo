// Package extinfo provides easy access to the state information of a Sauerbraten game server (called 'extinfo' in the Sauerbraten source code).
package extinfo

import (
	"net"
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
	EXTENDED_INFO_PLAYER_STATS = 0x01 //1
	EXTENDED_INFO_TEAMS_SCORES = 0x02 //2

	EXTENDED_INFO_PLAYER_STATS_RESPONSE_IDS   = 0xF6 //-10
	EXTENDED_INFO_PLAYER_STATS_RESPONSE_STATS = 0xF5 //-11
)

// A server to query extinfo from.
type Server struct {
	addr    *net.UDPAddr
	timeOut time.Duration
}

func NewServer(addr *net.UDPAddr, timeOut time.Duration) (s *Server) {
	// copy the address to not touch the original port
	addrCopy := *addr
	s = &Server{
		addr:    &addrCopy,
		timeOut: timeOut,
	}
	s.addr.Port++ // info port is at game server port + 1
	return
}
