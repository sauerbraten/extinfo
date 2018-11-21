// Package extinfo provides easy access to the state information of a Sauerbraten game server (called 'extinfo' in the Sauerbraten source code).
package extinfo

import (
	"net"
	"time"
)

// Protocol constants
const (
	// Constants describing the type of information to query for
	InfoTypeExtended byte = 0x00
	InfoTypeBasic    byte = 0x01

	// Constants used in responses to extended info queries
	ExtInfoACK     byte = 0xFF // -1
	ExtInfoVersion byte = 105
	ExtInfoError   byte = 0x01

	// Constants describing the type of extended information to query for
	ExtInfoTypeUptime     byte = 0x00
	ExtInfoTypeClientInfo byte = 0x01
	ExtInfoTypeTeamScores byte = 0x02

	// Constants used in responses to client info queries
	ClientInfoResponseTypeCNs  byte = 0xF6 // -10
	ClientInfoResponseTypeInfo byte = 0xF5 // -11
)

// Constants generally useful in this package
const (
	MaxPlayerCN     = 127 // Highest CN an actual player can have; bots' CNs start at 128
	MaxPacketLength = 512 // better to be safe
)

// Server represents a Sauerbraten game server.
type Server struct {
	addr    *net.UDPAddr
	timeOut time.Duration
}

// NewServer returns a Server to query information from.
func NewServer(addr net.UDPAddr, timeOut time.Duration) (*Server, error) {
	addr.Port++
	_addr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return nil, err
	}

	return &Server{
		addr:    _addr,
		timeOut: timeOut,
	}, nil
}
