package extinfo

import "net"

// BasicInfoRaw contains the information sent back from the server in their raw form, i.e. no translation from ints to strings, even if possible.
type BasicInfoRaw struct {
	NumberOfClients    int    // the number of clients currently connected to the server (players and spectators)
	ProtocolVersion    int    // version number of the protocol in use by the server
	GameMode           int    // current game mode
	SecsLeft           int    // the time left until intermission in seconds
	MaxNumberOfClients int    // the maximum number of clients the server allows
	MasterMode         int    // the current master mode of the server
	Map                string // current map
	Description        string // server description
}

// BasicInfo contains the parsed information sent back from the server, i.e. game mode and master mode are translated into human readable strings.
type BasicInfo struct {
	NumberOfClients    int    // the number of clients currently connected to the server (players and spectators)
	ProtocolVersion    int    // version number of the protocol in use by the server
	GameMode           string // current game mode
	SecsLeft           int    // the time left until intermission in seconds
	MaxNumberOfClients int    // the maximum number of clients the server allows
	MasterMode         string // the current master mode of the server
	Map                string // current map
	Description        string // server description
}

// BasicInfo contains the parsed information sent back from the server, i.e. state and privilege are translated into human readable strings.
type PlayerInfo struct {
	ClientNum int    // player's client number or cn
	Ping      int    // player's ping
	Name      string // player's displayed name
	Team      string // the team the player is on e.g. "good"
	Frags     int    // amount of frags/kills
	Flags     int    // amount of flags the player scored
	Deaths    int    // amount of deaths
	Teamkills int    // amount of teamkills
	Damage    int    // damage ?!?
	Health    int    // remaining HP (health points)
	Armour    int    // remaining armour
	Weapon    string // current weapon
	Privilege string // "none", "master" or "admin"
	State     string // player state, e.g. "dead" or "spectator"
	IP        net.IP // player IP (only the first 3 bytes)
}

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
	Weapon    int    // remaining armour
	Privilege int    // 0 ("none"), 1 ("master") or 2 ("admin")
	State     int    // player state, e.g. 1 ("alive") or 5 ("spectator"), see names.go for int -> string mapping
	IP        net.IP // player IP (only the first 3 bytes)
}

// TeamsScores (teams's scores) contains the game mode as human readable string, the seconds left in the game, and a slice of TeamScores
type TeamsScores struct {
	GameMode string      // current game mode
	SecsLeft int         // the time left until intermission in seconds
	Scores   []TeamScore // a team score for each team
}

// TeamsScoresRaw (teams's scores) contains the game mode as raw int, the seconds left in the game, and a slice of TeamScores
type TeamsScoresRaw struct {
	GameMode int         // current game mode
	SecsLeft int         // the time left until intermission in seconds
	Scores   []TeamScore // a team score for each team
}

// TeamScore (team score) contains the name of the team and the score, i.e. flags scored in flag modes / points gained for holding bases in capture modes / frags achieved in DM modes / skulls collected
type TeamScore struct {
	Name  string // name of the team, e.g. "good"
	Score int    // amount of points (e.g. flags in ctf modes, frags in deathmatch modes, base points in capture, and so on)
	Bases []int  // the numbers/IDs of the bases the team possesses
}
