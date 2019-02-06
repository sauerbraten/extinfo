package extinfo

import (
	"strconv"
)

// A slice containing the possible master modes
// In the sauerbraten protocol, 'auth' is -1, 'open' is 0, and so forth. Therefore getModeName() returns MasterModeNames[modeInt+1]. MasterModeNames should not be used directly
var masterModeNames = []string{"auth", "open", "veto", "locked", "private", "password"}

// wrapper function around masterModeNames
// returns the human readable name of the master mode as a string
func getMasterModeName(masterMode int) string {
	if masterMode < -1 || masterMode >= len(masterModeNames)-1 {
		return "unknown"
	}

	return masterModeNames[masterMode+1]
}

// A slice containing the possible game modes
// The index of a mode is equal to the int received in a response for that game mode, thus this slice maps the game mode ints to game mode strings
var gameModeNames = []string{"ffa", "coop edit", "teamplay", "instagib", "instagib team", "efficiency", "efficiency team", "tactics", "tactics team", "capture", "regen capture", "ctf", "insta ctf", "protect", "insta protect", "hold", "insta hold", "efficiency ctf", "efficiency protect", "efficiency hold", "collect", "insta collect", "efficiency collect"}

// wrapper function around gameModeNames
// returns the human readable name of the game mode as a string
func getGameModeName(gameMode int) string {
	if gameMode < 0 || gameMode >= len(gameModeNames) {
		return "unknown"
	}

	return gameModeNames[gameMode]
}

// IsTeamMode returns true when mode is a team mode, false otherwise.
func IsTeamMode(mode string) bool {
	switch mode {
	case "teamplay",
		"instagib team",
		"efficiency team",
		"tactics team",
		"capture",
		"regen capture",
		"ctf",
		"insta ctf",
		"protect",
		"insta protect",
		"hold",
		"insta hold",
		"efficiency ctf",
		"efficiency protect",
		"efficiency hold",
		"collect",
		"insta collect",
		"efficiency collect":
		return true
	default:
		return false
	}
}

// A slice containing the weapon names
// The index of a weapon is equal to the int received in a response for client info, thus this slice maps the weapon ints to weapon strings
var weaponNames = []string{"chain saw", "shotgun", "chain gun", "rocket launcher", "rifle", "grenade launcher", "pistol", "fire ball", "ice ball", "slime ball", "bite", "barrel"}

// wrapper function around weaponNames
// returns the human readable name of the weapon as a string
func getWeaponName(weapon int) string {
	if weapon < 0 || weapon >= len(weaponNames) {
		return "unknown"
	}

	return weaponNames[weapon]
}

// A slice containing the privilege names
// Maps the privilege ints to privilege strings
var privilegeNames = []string{"none", "master", "auth", "admin"}

// wrapper function around privilegeNames
// returns the human readable name of the privilege as a string
func getPrivilegeName(privilege int) string {
	if privilege < 0 || privilege >= len(privilegeNames) {
		return "unknown"
	}

	return privilegeNames[privilege]
}

// A slice containing the state names
// Maps the state ints to state strings
var stateNames = []string{"alive", "dead", "spawning", "lagged", "editing", "spectator"}

// wrapper function around stateNames
// returns the human readable name of the state as a string
func getStateName(state int) string {
	if state < 0 || state >= len(stateNames) {
		return "unknown"
	}

	return stateNames[state]
}

// getServerModName returns the human readable name of the server mod as a string
func getServerModName(mod int) string {
	switch mod {
	case -2:
		return "hopmod"
	case -3:
		return "oomod"
	case -4:
		return "spaghettimod"
	case -5:
		return "suckerserv"
	case -6:
		return "remod"
	case -7:
		return "noobmod"
	case -8:
		return "zeromod"
	case -9:
		return "waiter"
	default:
		return "unknown (" + strconv.Itoa(mod) + ")"
	}
}
