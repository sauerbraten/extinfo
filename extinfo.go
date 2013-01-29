// Package extinfo provides easy access to the state information of a Sauerbraten game server (called 'extinfo' in the Sauerbraten source code).
package extinfo

import (
	"errors"
	"net"
)

// the current position in a response ([]byte)
// needed, since values are encoded in variable amount of bytes
// global to not have to pass around an int on every dump
var positionInResponse int

// Constants describing the type of information to query for
const (
	extendedInfo = 0
	basicInfo    = 1
)

// Constants describing the type of extended information to query for
const (
	uptimeInfo      = 0
	playerStatsInfo = 1
	teamScoreInfo   = 2
)

// A server to query extinfo from.
type Server struct {
	addr string
	port int
}

func NewServer(addr string, port int) *Server {
	return &Server{addr, port}
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

	request := buildRequest(extendedInfo, teamScoreInfo, 0)
	response, err := queryServer(s.addr, s.port, request)
	if err != nil {
		return teamsScoresRaw, err
	}

	positionInResponse = 0

	// first int is extendedInfo = 0
	_ = dumpInt(response)

	// next int is teamScoreInfo = 2
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
		return teamsScoresRaw, errors.New("server is not running a team mode")
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
	response, err := queryServer(s.addr, s.port, buildRequest(basicInfo, 0, 0))
	if err != nil {
		return info, err
	}

	positionInResponse = 0

	// first int is basicInfo = 1
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
	basicInfoRaw := BasicInfoRaw{}

	response, err := queryServer(s.addr, s.port, buildRequest(basicInfo, 0, 0))
	if err != nil {
		return basicInfoRaw, err
	}

	positionInResponse = 0

	// first int is basicInfo = 1
	_ = dumpInt(response)
	basicInfoRaw.NumberOfClients = dumpInt(response)
	// next int is always 5, the number of additional attributes after the playercount and the strings for map and description
	//numberOfAttributes := dumpInt(response)
	_ = dumpInt(response)
	basicInfoRaw.ProtocolVersion = dumpInt(response)
	basicInfoRaw.GameMode = dumpInt(response)
	basicInfoRaw.SecsLeft = dumpInt(response)
	basicInfoRaw.MaxNumberOfClients = dumpInt(response)
	basicInfoRaw.MasterMode = dumpInt(response)
	basicInfoRaw.Map = dumpString(response)
	basicInfoRaw.Description = dumpString(response)

	return basicInfoRaw, nil
}

// GetUptime returns the uptime of the server in seconds.
func (s *Server) GetUptime() (int, error) {
	response, err := queryServer(s.addr, s.port, buildRequest(extendedInfo, uptimeInfo, 0))
	if err != nil {
		return -1, err
	}

	positionInResponse = 0

	// first int is extendedInfo
	_ = dumpInt(response)

	// next int is EXT_uptimeInfo = 0
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

	response, err := queryServer(s.addr, s.port, buildRequest(extendedInfo, playerStatsInfo, clientNum))
	if err != nil {
		return playerInfo, err
	}

	if response[5] != 0x00 {
		// there was an error
		return playerInfo, errors.New("invalid cn")
	}

	// throw away 7 first ints (extendedInfo, playerStatsInfo, clientNum, server ACK byte, server VERSION byte, server NO_ERROR byte, server playerStatsInfo_RESP_STATS byte)
	response = response[7:]

	playerInfo = parsePlayerInfo(response)

	return playerInfo, nil
}

// GetPlayerInfoRaw returns the raw information about the player with the given clientNum.
func (s *Server) GetPlayerInfoRaw(clientNum int) (PlayerInfoRaw, error) {
	playerInfoRaw := PlayerInfoRaw{}

	response, err := queryServer(s.addr, s.port, buildRequest(extendedInfo, playerStatsInfo, clientNum))
	if err != nil {
		return playerInfoRaw, err
	}

	if response[5] != 0x00 {
		// there was an error
		return playerInfoRaw, errors.New("invalid cn")
	}

	// throw away 7 first ints (extendedInfo, playerStatsInfo, clientNum, server ACK byte, server VERSION byte, server NO_ERROR byte, server playerStatsInfo_RESP_STATS byte)
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
func (s *Server) GetAllPlayerInfo() ([]PlayerInfo, error) {
	allPlayerInfo := []PlayerInfo{}

	response, err := queryServer(s.addr, s.port, buildRequest(extendedInfo, playerStatsInfo, -1))
	if err != nil {
		return allPlayerInfo, err
	}

	// response is multiple 64-byte responses, one for each player
	playerCount := len(response) / 64

	// parse each 64 byte packet (without the first 7 bytes) on its own and append to allPlayerInfo
	for i := 0; i < playerCount; i++ {
		allPlayerInfo = append(allPlayerInfo, parsePlayerInfo(response[i*64+7:(i*64)+64]))
	}

	return allPlayerInfo, nil
}
