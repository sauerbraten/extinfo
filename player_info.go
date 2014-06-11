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
func (s *Server) GetPlayerInfoRaw(clientNum int) (PlayerInfoRaw, error) {
	playerInfoRaw := PlayerInfoRaw{}

	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, clientNum))
	if err != nil {
		return playerInfoRaw, err
	}

	if response[5] != EXTENDED_INFO_NO_ERROR {
		// server says the cn was invalid
		return playerInfoRaw, errors.New("extinfo: invalid cn\n")
	}

	return parsePlayerInfoResponse(response), nil
}

// GetPlayerInfo returns the parsed information about the player with the given clientNum.
func (s *Server) GetPlayerInfo(clientNum int) (PlayerInfo, error) {
	playerInfo := PlayerInfo{}

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
func (s *Server) GetAllPlayerInfo() (map[int]PlayerInfo, error) {
	allPlayerInfo := map[int]PlayerInfo{}

	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, -1))
	if err != nil {
		return allPlayerInfo, err
	}

	// response is multiple 64-byte responses, one for each player
	// parse each 64 byte packet (without the first 7 bytes) on its own and append to allPlayerInfo
	for i := 0; i < len(response); i += 64 {
		playerInfoRaw := parsePlayerInfoResponse(response[i : i+64])
		allPlayerInfo[playerInfoRaw.ClientNum] = PlayerInfo{playerInfoRaw, getWeaponName(playerInfoRaw.Weapon), getStateName(playerInfoRaw.State), getPrivilegeName(playerInfoRaw.Privilege)}
	}

	return allPlayerInfo, nil
}

// own function, because it is used in GetPlayerInfo() + GetAllPlayerInfo()
func parsePlayerInfoResponse(response []byte) PlayerInfoRaw {
	// throw away 7 first bytes (EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, cn, EXTENDED_INFO_ACK, EXTENDED_INFO_VERSION, EXTENDED_INFO_NO_ERROR, EXTENDED_INFO_PLAYER_STATS_RESP_STATS)
	response = response[7:]

	positionInResponse = 0

	playerInfoRaw := PlayerInfoRaw{
		ClientNum: dumpInt(response),
		Ping:      dumpInt(response),
		Name:      dumpString(response),
		Team:      dumpString(response),
		Frags:     dumpInt(response),
		Flags:     dumpInt(response),
		Deaths:    dumpInt(response),
		Teamkills: dumpInt(response),
		Damage:    dumpInt(response),
		Health:    dumpInt(response),
		Armour:    dumpInt(response),
		Weapon:    dumpInt(response),
		Privilege: dumpInt(response),
		State:     dumpInt(response),
	}

	// IP from next 4 bytes
	ipBytes := response[positionInResponse : positionInResponse+4]
	playerInfoRaw.IP = net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])

	return playerInfoRaw
}
