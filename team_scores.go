package extinfo

import (
	"errors"
)

// TeamScore contains the name of the team and the score, i.e. flags scored in flag modes / points gained for holding bases in capture modes / frags achieved in DM modes / skulls collected
type TeamScore struct {
	Name  string // name of the team, e.g. "good"
	Score int    // flags in ctf modes, frags in deathmatch modes, points in capture, skulls in collect
	Bases []int  // the numbers/IDs of the bases the team possesses (only used in capture modes)
}

// TeamScoresRaw contains the game mode as raw int, the seconds left in the game, and a slice of TeamScores
type TeamScoresRaw struct {
	GameMode int                  // current game mode
	SecsLeft int                  // the time left until intermission in seconds
	Scores   map[string]TeamScore // a team score for each team, mapped to the team's name
}

// TeamScores contains the game mode as human readable string, the seconds left in the game, and a slice of TeamScores
type TeamScores struct {
	TeamScoresRaw
	GameMode string // current game mode
}

// GetTeamScoresRaw queries a Sauerbraten server at addr on port for the teams' names and scores and returns the raw response and/or an error in case something went wrong or the server is not running a team mode.
func (s *Server) GetTeamScoresRaw() (teamScoresRaw TeamScoresRaw, err error) {
	request := buildRequest(EXTENDED_INFO, EXTENDED_INFO_TEAMS_SCORES, 0)
	response, err := s.queryServer(request)
	if err != nil {
		return
	}

	// ignore first 3 bytes: EXTENDED_INFO, EXTENDED_INFO_TEAMS_SCORES, EXTENDED_INFO_ACK
	response = response[3:]

	positionInResponse = 0

	// check for correct extinfo protocol version
	version, err := dumpInt(response)
	if err != nil {
		return
	}
	if version != EXTENDED_INFO_VERSION {
		err = errors.New("extinfo: wrong extinfo protocol version")
		return
	}

	// next int describes wether the server runs a team mode or not
	isTeamMode := true
	teamModeErrorValue, err := dumpInt(response)
	if err != nil {
		return
	}
	if teamModeErrorValue != 0 {
		isTeamMode = false
	}

	teamScoresRaw.GameMode, err = dumpInt(response)
	if err != nil {
		return
	}
	teamScoresRaw.SecsLeft, err = dumpInt(response)
	if err != nil {
		return
	}

	if !isTeamMode {
		// no team scores following
		err = errors.New("extinfo: server is not running a team mode")
		return
	}

	for response[positionInResponse] != 0x0 {
		var name string
		name, err = dumpString(response)
		if err != nil {
			return
		}

		var score int
		score, err = dumpInt(response)
		if err != nil {
			return
		}

		var numBases int
		numBases, err = dumpInt(response)
		if err != nil {
			return
		}

		bases := make([]int, numBases)

		for i := 0; i < numBases; i++ {
			var base int
			base, err = dumpInt(response)
			if err != nil {
				return
			}
			bases = append(bases, base)
		}

		teamScoresRaw.Scores[name] = TeamScore{name, score, bases}
	}

	return
}

// GetTeamScores queries a Sauerbraten server at addr on port for the teams' names and scores and returns the parsed response and/or an error in case something went wrong or the server is not running a team mode. Parsed response means that the int value sent as game mode is translated into the human readable name, e.g. '12' -> "insta ctf".
func (s *Server) GetTeamScores() (TeamScores, error) {
	teamScores := TeamScores{}

	teamScoresRaw, err := s.GetTeamScoresRaw()
	if err != nil {
		return teamScores, err
	}

	teamScores.TeamScoresRaw = teamScoresRaw
	teamScores.GameMode = getGameModeName(teamScoresRaw.GameMode)

	return teamScores, nil
}