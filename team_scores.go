package extinfo

// TeamScore contains the name of the team and the score, i.e. flags scored in flag modes / points gained for holding bases in capture modes / frags achieved in DM modes / skulls collected
type TeamScore struct {
	Name  string `json:"name"`  // name of the team, e.g. "good"
	Score int    `json:"score"` // flags in ctf modes, frags in deathmatch modes, points in capture, skulls in collect
	Bases []int  `json:"bases"` // the numbers/IDs of the bases the team possesses (only used in capture modes)
}

// TeamScoresRaw contains the game mode as raw int, the seconds left in the game, and a slice of TeamScores
type TeamScoresRaw struct {
	GameMode int                  `json:"gameMode"` // current game mode
	SecsLeft int                  `json:"secsLeft"` // the time left until intermission in seconds
	Scores   map[string]TeamScore `json:"scores"`   // a team score for each team, mapped to the team's name
}

// TeamScores contains the game mode as human readable string, the seconds left in the game, and a slice of TeamScores
type TeamScores struct {
	TeamScoresRaw
	GameMode string `json:"gameMode"` // current game mode
}

// GetTeamScoresRaw queries a Sauerbraten server at addr on port for the teams' names and scores and returns the raw response and/or an error in case something went wrong or the server is not running a team mode.
func (s *Server) GetTeamScoresRaw() (teamScoresRaw TeamScoresRaw, err error) {
	request := buildRequest(InfoTypeExtended, ExtInfoTypeTeamScores, 0)
	response, err := s.queryServer(request)
	if err != nil {
		return
	}

	teamScoresRaw.GameMode, err = response.ReadInt()
	if err != nil {
		return
	}

	teamScoresRaw.SecsLeft, err = response.ReadInt()
	if err != nil {
		return
	}

	teamScoresRaw.Scores = map[string]TeamScore{}

	for response.HasRemaining() {
		var name string
		name, err = response.ReadString()
		if err != nil {
			return
		}

		var score int
		score, err = response.ReadInt()
		if err != nil {
			return
		}

		var numBases int
		numBases, err = response.ReadInt()
		if err != nil {
			return
		}

		if numBases < 0 {
			numBases = 0
		}

		bases := make([]int, numBases)

		for i := 0; i < numBases; i++ {
			var base int
			base, err = response.ReadInt()
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
