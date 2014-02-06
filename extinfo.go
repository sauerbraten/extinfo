// Package extinfo provides easy access to the state information of a Sauerbraten game server (called 'extinfo' in the Sauerbraten source code).
package extinfo

import (
	"errors"
	"net"
	"time"
)

// Constants describing the type of information to query for
const (
	EXTENDED_INFO = 0
	BASIC_INFO    = 1
)

// Constants describing the type of extended information to query for
const (
	EXTENDED_INFO_UPTIME       = 0
	EXTENDED_INFO_PLAYER_STATS = 1
	EXTENDED_INFO_TEAMS_SCORES = 2
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
	s.addr.Port++ // extinfo port is at game port + 1
	return
}

// GetTeamsScores queries a Sauerbraten server at addr on port for the teams' names and scores and returns the parsed response and/or an error in case something went wrong or the server is not running a team mode. Parsed response means that the int value sent as game mode is translated into the human readable name, e.g. '12' -> "insta ctf".
func (s *Server) GetTeamsScores() (TeamsScores, error) {
	teamsScoresRaw, err := s.GetTeamsScoresRaw()
	teamsScores := TeamsScores{getGameModeName(teamsScoresRaw.GameMode), teamsScoresRaw.SecsLeft, teamsScoresRaw.Scores}
	return teamsScores, err
}

// GetTeamsScoresRaw queries a Sauerbraten server at addr on port for the teams' names and scores and returns the raw response and/or an error in case something went wrong or the server is not running a team mode.
func (s *Server) GetTeamsScoresRaw() (TeamsScoresRaw, error) {
	teamsScoresRaw := TeamsScoresRaw{}

	request := buildRequest(EXTENDED_INFO, EXTENDED_INFO_TEAMS_SCORES, 0)
	response, err := s.queryServer(request)
	if err != nil {
		return teamsScoresRaw, err
	}

	positionInResponse = 0

	// first int is EXTENDED_INFO = 0
	_ = dumpInt(response)

	// next int is EXTENDED_INFO_TEAMS_SCORES = 2
	_ = dumpInt(response)

	// next int is EXT_ACK = -1
	_ = dumpInt(response)

	// next int is EXT_VERSION
	_ = dumpInt(response)

	// next int describes wether the server runs a team mode or not
	isTeamMode := true
	if dumpInt(response) != 0 {
		isTeamMode = false
	}

	teamsScoresRaw.GameMode = dumpInt(response)
	teamsScoresRaw.SecsLeft = dumpInt(response)

	if !isTeamMode {
		// no team scores following
		return teamsScoresRaw, errors.New("extinfo: server is not running a team mode\n")
	}

	name := ""
	score := 0
	numBases := 0

	for response[positionInResponse] != 0x00 {
		name = dumpString(response)
		score = dumpInt(response)
		numBases = dumpInt(response)

		bases := make([]int, 0)

		for i := 0; i < numBases; i++ {
			bases = append(bases, dumpInt(response))
		}

		teamsScoresRaw.Scores = append(teamsScoresRaw.Scores, TeamScore{name, score, bases})
	}

	return teamsScoresRaw, nil
}

// GetBasicInfo queries a Sauerbraten server at addr on port and returns the parsed response or an error in case something went wrong. Parsed response means that the int values sent as game mode and master mode are translated into the human readable name, e.g. '12' -> "insta ctf".
func (s *Server) GetBasicInfo() (info BasicInfo, err error) {
	var response []byte
	response, err = s.queryServer(buildRequest(BASIC_INFO, 0, 0))
	if err != nil {
		return info, err
	}

	positionInResponse = 0

	// first int is BASIC_INFO = 1
	_ = dumpInt(response)

	info.NumberOfClients = dumpInt(response)
	// next int is always 5, the number of additional attributes after the playercount and the strings for map and description
	//numberOfAttributes := dumpInt(response)
	_ = dumpInt(response)
	info.ProtocolVersion = dumpInt(response)
	info.GameMode = getGameModeName(dumpInt(response))
	info.SecsLeft = dumpInt(response)
	info.MaxNumberOfClients = dumpInt(response)
	info.MasterMode = getMasterModeName(dumpInt(response))
	info.Map = dumpString(response)
	info.Description = dumpString(response)

	return
}

// GetBasicInfoRaw queries a Sauerbraten server at addr on port and returns the raw response or an error in case something went wrong. Raw response means that the int values sent as game mode and master mode are NOT translated into the human readable name.
func (s *Server) GetBasicInfoRaw() (BasicInfoRaw, error) {
	BASIC_INFORaw := BasicInfoRaw{}

	response, err := s.queryServer(buildRequest(BASIC_INFO, 0, 0))
	if err != nil {
		return BASIC_INFORaw, err
	}

	positionInResponse = 0

	// first int is BASIC_INFO = 1
	_ = dumpInt(response)
	BASIC_INFORaw.NumberOfClients = dumpInt(response)
	// next int is always 5, the number of additional attributes after the playercount and the strings for map and description
	//numberOfAttributes := dumpInt(response)
	_ = dumpInt(response)
	BASIC_INFORaw.ProtocolVersion = dumpInt(response)
	BASIC_INFORaw.GameMode = dumpInt(response)
	BASIC_INFORaw.SecsLeft = dumpInt(response)
	BASIC_INFORaw.MaxNumberOfClients = dumpInt(response)
	BASIC_INFORaw.MasterMode = dumpInt(response)
	BASIC_INFORaw.Map = dumpString(response)
	BASIC_INFORaw.Description = dumpString(response)

	return BASIC_INFORaw, nil
}

// GetUptime returns the uptime of the server in seconds.
func (s *Server) GetUptime() (int, error) {
	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_UPTIME, 0))
	if err != nil {
		return -1, err
	}

	positionInResponse = 0

	// first int is EXTENDED_INFO
	_ = dumpInt(response)

	// next int is EXT_EXTENDED_INFO_UPTIME = 0
	_ = dumpInt(response)

	// next int is EXT_ACK = -1
	_ = dumpInt(response)

	// next int is EXT_VERSION
	_ = dumpInt(response)

	// next int is the actual uptime
	uptime := dumpInt(response)

	return uptime, nil
}

// GetPlayerInfo returns the parsed information about the player with the given clientNum.
func (s *Server) GetPlayerInfo(clientNum int) (PlayerInfo, error) {
	playerInfo := PlayerInfo{}

	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, clientNum))
	if err != nil {
		return playerInfo, err
	}

	if response[5] != 0x00 {
		// server says the cn was invalid
		return playerInfo, errors.New("extinfo: invalid cn\n")
	}

	// throw away 7 first ints (EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, clientNum, server ACK byte, server VERSION byte, server NO_ERROR byte, server EXTENDED_INFO_PLAYER_STATS_RESP_STATS byte)
	response = response[7:]

	playerInfo = parsePlayerInfo(response)

	return playerInfo, nil
}

// GetPlayerInfoRaw returns the raw information about the player with the given clientNum.
func (s *Server) GetPlayerInfoRaw(clientNum int) (PlayerInfoRaw, error) {
	playerInfoRaw := PlayerInfoRaw{}

	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, clientNum))
	if err != nil {
		return playerInfoRaw, err
	}

	if response[5] != 0x00 {
		// server says the cn was invalid
		return playerInfoRaw, errors.New("extinfo: invalid cn\n")
	}

	// throw away 7 first ints (EXTENDED_INFO, EXTENDED_INFO_PLAYER_STATS, clientNum, server ACK byte, server VERSION byte, server NO_ERROR byte, server EXTENDED_INFO_PLAYER_STATS_RESP_STATS byte)
	response = response[7:]

	positionInResponse = 0

	playerInfoRaw.ClientNum = dumpInt(response)
	playerInfoRaw.Ping = dumpInt(response)
	playerInfoRaw.Name = dumpString(response)
	playerInfoRaw.Team = dumpString(response)
	playerInfoRaw.Frags = dumpInt(response)
	playerInfoRaw.Flags = dumpInt(response)
	playerInfoRaw.Deaths = dumpInt(response)
	playerInfoRaw.Teamkills = dumpInt(response)
	playerInfoRaw.Damage = dumpInt(response)
	playerInfoRaw.Health = dumpInt(response)
	playerInfoRaw.Armour = dumpInt(response)
	playerInfoRaw.Weapon = dumpInt(response)
	playerInfoRaw.Privilege = dumpInt(response)
	playerInfoRaw.State = dumpInt(response)
	// IP from next 4 bytes
	ip := response[positionInResponse : positionInResponse+4]
	playerInfoRaw.IP = net.IPv4(ip[0], ip[1], ip[2], ip[3])

	return playerInfoRaw, nil
}

// own function, because it is used in GetPlayerInfo() + GetAllPlayerInfo()
func parsePlayerInfo(response []byte) PlayerInfo {
	playerInfo := PlayerInfo{}

	positionInResponse = 0

	playerInfo.ClientNum = dumpInt(response)
	playerInfo.Ping = dumpInt(response)
	playerInfo.Name = dumpString(response)
	playerInfo.Team = dumpString(response)
	playerInfo.Frags = dumpInt(response)
	playerInfo.Flags = dumpInt(response)
	playerInfo.Deaths = dumpInt(response)
	playerInfo.Teamkills = dumpInt(response)
	playerInfo.Damage = dumpInt(response)
	playerInfo.Health = dumpInt(response)
	playerInfo.Armour = dumpInt(response)
	playerInfo.Weapon = getWeaponName(dumpInt(response))
	playerInfo.Privilege = getPrivilegeName(dumpInt(response))
	playerInfo.State = getStateName(dumpInt(response))
	// IP from next 4 bytes
	ipBytes := response[positionInResponse : positionInResponse+4]
	playerInfo.IP = net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])

	return playerInfo
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
		playerInfo := parsePlayerInfo(response[i+7 : i+64])
		allPlayerInfo[playerInfo.ClientNum] = playerInfo
	}

	return allPlayerInfo, nil
}
