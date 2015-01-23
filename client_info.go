package extinfo

import (
	"net"

	"github.com/sauerbraten/cubecode"
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
	Accuracy  int    // damage the client could have dealt * 100 / damage actually dealt by the client
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
	for i := 0; i < response.Len(); i += 64 {
		partialResponse := response.SubPacket(i, i+64)
		clientInfoRaw, err = parseClientInfoResponse(partialResponse)
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

// own function, because it is used in GetClientInfo() & GetAllClientInfo()
func parseClientInfoResponse(response *cubecode.Packet) (clientInfoRaw ClientInfoRaw, err error) {
	// omit 7 first bytes: EXTENDED_INFO, EXTENDED_INFO_CLIENT_INFO, CN, EXTENDED_INFO_ACK, EXTENDED_INFO_VERSION, EXTENDED_INFO_NO_ERROR, EXTENDED_INFO_CLIENT_INFO_RESPONSE_INFO
	for i := 0; i < 7; i++ {
		_, err = response.ReadInt()
		if err != nil {
			return
		}
	}

	clientInfoRaw.ClientNum, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Ping, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Name, err = response.ReadString()
	if err != nil {
		return
	}

	clientInfoRaw.Team, err = response.ReadString()
	if err != nil {
		return
	}

	clientInfoRaw.Frags, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Flags, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Deaths, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Teamkills, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Accuracy, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Health, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Armour, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Weapon, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.Privilege, err = response.ReadInt()
	if err != nil {
		return
	}

	clientInfoRaw.State, err = response.ReadInt()
	if err != nil {
		return
	}

	// IP from next 4 bytes
	var ipByte1, ipByte2, ipByte3, ipByte4 byte

	ipByte1, err = response.ReadByte()
	if err != nil {
		return
	}

	ipByte2, err = response.ReadByte()
	if err != nil {
		return
	}

	ipByte3, err = response.ReadByte()
	if err != nil {
		return
	}

	ipByte4 = 0

	clientInfoRaw.IP = net.IPv4(ipByte1, ipByte2, ipByte3, ipByte4)

	return
}
