package extinfo

import (
	"errors"
	"net"
)

// BasicInfoRaw contains the raw information sent back from the server, i.e. state and privilege are ints.
type PlayerInfoRaw struct {
	ClientNum int    // player's client number or cn
	Ping      int    // player's ping
	Name      string // player's displayed name
	Team      string // the team the player is on, e.g. "good"
	Frags     int    // amount of frags/kills
	Flags     int    // amount of flags the player scored
	Deaths    int    // amount of deaths
	Teamkills int    // amount of teamkills
	Damage    int    // damage ?!?
	Health    int    // remaining HP (health points)
	Armour    int    // remaining armour
	Weapon    int    // weapon the player currently has selected
	Privilege int    // 0 ("none"), 1 ("master") or 2 ("admin")
	State     int    // player state, e.g. 1 ("alive") or 5 ("spectator"), see names.go for int -> string mapping
	IP        net.IP // player IP (only the first 3 bytes)
}

// BasicInfo contains the parsed information sent back from the server, i.e. weapon, state and privilege are translated into human readable strings.
type PlayerInfo struct {
	PlayerInfoRaw
	Weapon    string // weapon the player currently has selected
	Privilege string // "none", "master" or "admin"
	State     string // player state, e.g. "dead" or "spectator"
}

// GetPlayerInfoRaw returns the raw information about the player with the given clientNum.
func (s *Server) GetPlayerInfoRaw(clientNum int) (playerInfoRaw PlayerInfoRaw, err error) {
	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, clientNum))
	if err != nil {
		return
	}

	playerInfoRaw, err = parsePlayerInfoResponse(response)
	return
}

// GetPlayerInfo returns the parsed information about the player with the given clientNum.
func (s *Server) GetPlayerInfo(clientNum int) (playerInfo PlayerInfo, err error) {
	playerInfoRaw, err := s.GetPlayerInfoRaw(clientNum)
	if err != nil {
		return playerInfo, err
	}

	playerInfo.PlayerInfoRaw = playerInfoRaw
	playerInfo.Weapon = getWeaponName(playerInfo.PlayerInfoRaw.Weapon)
	playerInfo.Privilege = getPrivilegeName(playerInfo.PlayerInfoRaw.Privilege)
	playerInfo.State = getStateName(playerInfo.PlayerInfoRaw.State)

	return playerInfo, nil
}

// GetAllPlayerInfo returns the Information of all Players (including spectators) as a []PlayerInfo
func (s *Server) GetAllPlayerInfo() (allPlayerInfo map[int]PlayerInfo, err error) {
	allPlayerInfo = map[int]PlayerInfo{}

	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, -1))
	if err != nil {
		return allPlayerInfo, err
	}

	// response is multiple 64-byte responses, one for each player
	// parse each 64 byte packet (without the first 7 bytes) on its own and append to allPlayerInfo
	playerInfoRaw := PlayerInfoRaw{}
	for i := 0; i < len(response); i += 64 {
		playerInfoRaw, err = parsePlayerInfoResponse(response[i : i+64])
		if err != nil {
			return
		}

		allPlayerInfo[playerInfoRaw.ClientNum] = PlayerInfo{playerInfoRaw, getWeaponName(playerInfoRaw.Weapon), getStateName(playerInfoRaw.State), getPrivilegeName(playerInfoRaw.Privilege)}
	}

	return
}

// own function, because it is used in GetPlayerInfo() + GetAllPlayerInfo()
func parsePlayerInfoResponse(response []byte) (playerInfoRaw PlayerInfoRaw, err error) {
	// throw away 4 first bytes (EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, cn, EXTENDED_INFO_ACK)
	response = response[4:]

	positionInResponse = 0

	// next three bytes are EXTENDED_INFO_VERSION, EXTENDED_INFO_NO_ERROR, EXTENDED_INFO_PLAYER_STATS_RESP_STATS
	// check for correct extinfo protocol version
	if dumpByte(response) != EXTENDED_INFO_VERSION {
		err = errors.New("extinfo: wrong extinfo protocol version")
		return
	}

	if dumpByte(response) != EXTENDED_INFO_NO_ERROR {
		err = errors.New("extinfo: invalid client number")
		return
	}

	if dumpByte(response) != EXTENDED_INFO_PLAYER_STATS_RESPONSE_STATS {
		err = errors.New("extinfo: illegal response type")
		return
	}

	// set fields in raw player info

	playerInfoRaw.ClientNum, err = dumpInt(response)
	if err != nil {
		return
	}

	playerInfoRaw.Ping, err = dumpInt(response)
	if err != nil {
		return
	}

	playerInfoRaw.Name, err = dumpString(response)
	if err != nil {
		return
	}

	playerInfoRaw.Team, err = dumpString(response)
	if err != nil {
		return
	}

	playerInfoRaw.Frags, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Flags, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Deaths, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Teamkills, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Damage, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Health, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Armour, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Weapon, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.Privilege, err = dumpInt(response)
	if err != nil {
		return
	}
	playerInfoRaw.State, err = dumpInt(response)
	if err != nil {
		return
	}

	// IP from next 4 bytes
	ipBytes := response[positionInResponse : positionInResponse+4]
	playerInfoRaw.IP = net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])

	return
}
