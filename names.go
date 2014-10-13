package extinfo

// A slice containing the possible master modes
// In the sauerbraten protocol, 'auth' is -1, 'open' is 0, and so forth. Therefore getModeName() returns MasterModeNames[modeInt+1]. MasterModeNames should not be used directly
var masterModeNames = []string{"auth", "open", "veto", "locked", "private", "password"}

// A slice containing the possible game modes
// The index of a mode is equal to the int received in a response for that game mode, thus this slice maps the game mode ints to game mode strings
var gameModeNames = []string{"ffa", "coop edit", "teamplay", "instagib", "instagib team", "efficiency", "efficieny team", "tactics", "tactics team", "capture", "regen capture", "ctf", "insta ctf", "protect", "insta protect", "hold", "insta hold", "efficiency ctf", "efficiency protect", "efficiency hold", "collect", "insta collect", "efficiency collect"}

// A slice containing the weapon names
// The index of a weapon is equal to the int received in a response for player info, thus this slice maps the weapon ints to weapon strings
var weaponNames = []string{"chain saw", "shotgun", "chain gun", "rocket launcher", "rifle", "grenade launcher", "pistol", "fire ball", "ice ball", "slime ball", "bite", "barrel"}

// A slice containing the privilege names
// Maps the privilege ints to privilege strings
var privilegeNames = []string{"none", "master", "auth", "admin"}

// A slice containing the state names
// Maps the state ints to state strings
var stateNames = []string{"alive", "dead", "spawning", "lagged", "edited", "spectator"}

// wrapper function around masterModeNames
// returns the human readable name of the master mode as a string
func getMasterModeName(masterMode int) string {
	if masterMode < -1 || masterMode >= len(masterModeNames)-1 {
		return "unknown"
	}

	return masterModeNames[masterMode+1]
}

// wrapper function around gameModeNames
// returns the human readable name of the game mode as a string
func getGameModeName(gameMode int) string {
	if gameMode < 0 || gameMode >= len(gameModeNames) {
		return "unknown"
	}

	return gameModeNames[gameMode]
}

// wrapper function around weaponNames
// returns the human readable name of the weapon as a string
func getWeaponName(weapon int) string {
	if weapon < 0 || weapon >= len(weaponNames) {
		return "unknown"
	}

	return weaponNames[weapon]
}

// wrapper function around privilegeNames
// returns the human readable name of the privilege as a string
func getPrivilegeName(privilege int) string {
	if privilege < 0 || privilege >= len(privilegeNames) {
		return "unknown"
	}

	return privilegeNames[privilege]
}

// wrapper function around stateNames
// returns the human readable name of the state as a string
func getStateName(state int) string {
	if state < 0 || state >= len(stateNames) {
		return "unknown"
	}

	return stateNames[state]
}
