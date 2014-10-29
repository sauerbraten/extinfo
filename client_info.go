package extinfo

import (
	"errors"
	"log"
	"net"
)

// ClientInfoRaw contains the raw information sent back from the server, i.e. state and privilege are ints.
type ClientInfoRaw struct {
	ClientNum int    // client number or cn
	Ping      int    // client's ping to server
	Name      string //
	Team      string // name of the team the client is on, e.g. "good"
	Frags     int    // kills
	Flags     int    // number of flags the player scored
	Deaths    int    //
	Teamkills int    //
	Damage    int    // damage dealt by the client
	Health    int    // remaining HP (health points)
	Armour    int    // remaining armour
	Weapon    int    // weapon the client currently has selected
	Privilege int    // 0 ("none"), 1 ("master") or 2 ("admin")
	State     int    // client state, e.g. 1 ("alive") or 5 ("spectator"), see names.go for int -> string mapping
	IP        net.IP // client IP (only the first 3 bytes)
}

// ClientInfo contains the parsed information sent back from the server, i.e. weapon, state and privilege are translated into human readable strings.
type ClientInfo struct {
	ClientInfoRaw
	Weapon    string // weapon the client currently has selected
	Privilege string // "none", "master" or "admin"
	State     string // client state, e.g. "dead" or "spectator"
}

// GetClientInfoRaw returns the raw information about the client with the given clientNum.
func (s *Server) GetClientInfoRaw(clientNum int) (clientInfoRaw ClientInfoRaw, err error) {
	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_CLIENT_INFO, clientNum))
	if err != nil {
		return
	}

	clientInfoRaw, err = parseClientInfoResponse(response)
	return
}

// GetClientInfo returns the parsed information about the client with the given clientNum.
func (s *Server) GetClientInfo(clientNum int) (clientInfo ClientInfo, err error) {
	clientInfoRaw, err := s.GetClientInfoRaw(clientNum)
	if err != nil {
		return clientInfo, err
	}

	clientInfo.ClientInfoRaw = clientInfoRaw
	clientInfo.Weapon = getWeaponName(clientInfo.ClientInfoRaw.Weapon)
	clientInfo.Privilege = getPrivilegeName(clientInfo.ClientInfoRaw.Privilege)
	clientInfo.State = getStateName(clientInfo.ClientInfoRaw.State)

	return clientInfo, nil
}

// GetAllClientInfo returns the ClientInfo of all Players (including spectators) as a []ClientInfo
func (s *Server) GetAllClientInfo() (allClientInfo map[int]ClientInfo, err error) {
	allClientInfo = map[int]ClientInfo{}

	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_CLIENT_INFO, -1))
	if err != nil {
		return allClientInfo, err
	}

	// response is multiple 64-byte responses, one for each client
	// parse each 64 byte packet on its own and append to allClientInfo
	clientInfoRaw := ClientInfoRaw{}
	for i := 0; i < len(response); i += 64 {
		clientInfoRaw, err = parseClientInfoResponse(response[i : i+64])
		if err != nil {
			return
		}

		allClientInfo[clientInfoRaw.ClientNum] = ClientInfo{
			ClientInfoRaw: clientInfoRaw,
			Weapon:        getWeaponName(clientInfoRaw.Weapon),
			Privilege:     getPrivilegeName(clientInfoRaw.Privilege),
			State:         getStateName(clientInfoRaw.State),
		}
	}

	return
}

// own function, because it is used in GetClientInfo() + GetAllClientInfo()
func parseClientInfoResponse(response []byte) (clientInfoRaw ClientInfoRaw, err error) {
	log.Println(response)
	// throw away 4 first bytes (EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, cn, EXTENDED_INFO_ACK)
	response = response[4:]

	positionInResponse = 0

	// next three bytes are EXTENDED_INFO_VERSION, EXTENDED_INFO_NO_ERROR, EXTENDED_INFO_CLIENT_INFO_RESPONSE_INFO

	// check for correct extinfo protocol version
	if dumpByte(response) != EXTENDED_INFO_VERSION {
		err = errors.New("extinfo: wrong extinfo protocol version")
		return
	}

	if dumpByte(response) != EXTENDED_INFO_NO_ERROR {
		err = errors.New("extinfo: invalid client number")
		return
	}

	if dumpByte(response) != EXTENDED_INFO_CLIENT_INFO_RESPONSE_INFO {
		err = errors.New("extinfo: illegal response type")
		return
	}

	// set fields in raw client info

	clientInfoRaw.ClientNum, err = dumpInt(response)
	if err != nil {
		return
	}

	clientInfoRaw.Ping, err = dumpInt(response)
	if err != nil {
		return
	}

	clientInfoRaw.Name, err = dumpString(response)
	if err != nil {
		return
	}

	clientInfoRaw.Team, err = dumpString(response)
	if err != nil {
		return
	}

	clientInfoRaw.Frags, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Flags, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Deaths, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Teamkills, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Damage, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Health, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Armour, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Weapon, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.Privilege, err = dumpInt(response)
	if err != nil {
		return
	}
	clientInfoRaw.State, err = dumpInt(response)
	if err != nil {
		return
	}

	// IP from next 4 bytes
	ipBytes := response[positionInResponse : positionInResponse+4]
	clientInfoRaw.IP = net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])

	return
}
