# extinfo

A  Go package to query information ('extinfo') from a Sauerbraten game server.

## Usage

Get the package:

	$ go get github.com/sauerbraten/extinfo

Import the package:

	import (
		"github.com/sauerbraten/extinfo"
	)

## Documentation

Detailed Documentation [here](http://godoc.org/github.com/sauerbraten/extinfo).

## Example

Here is code to get you the state of the PSL1 server:

	package main

	import (
		"fmt"
		"github.com/sauerbraten/extinfo"
		"net"
	)

	func main() {
		addr, err := net.ResolveUDPAddr("udp", "sauerleague.org:10000")
		if err != nil {
			fmt.Println(err)
			return
		}

		psl1 := extinfo.NewServer(addr)
		basicInfo, err := psl1.GetBasicInfo()
		if err != nil {
			fmt.Print("Error getting basic information: ", err)
			return
		}

		fmt.Print("Basic Server Information:\n")
		fmt.Printf("Description:                %v\n", basicInfo.Description)
		fmt.Printf("Master Mode:                %v\n", basicInfo.MasterMode)
		fmt.Printf("Game Mode:                  %v\n", basicInfo.GameMode)
		fmt.Printf("Map:                        %v\n", basicInfo.Map)
		fmt.Printf("Players:                    %v\n", basicInfo.NumberOfClients)
		fmt.Printf("Maximum Number of Players:  %v\n", basicInfo.MaxNumberOfClients)
		fmt.Printf("Time Left (seconds):        %v\n", basicInfo.SecsLeft)
		fmt.Printf("Protocol Version:           %v\n", basicInfo.ProtocolVersion)
	}

The output should be something like this:

	Basic Server Information:
	Description:                PSL.sauerleague.org
	Master Mode:                auth
	Game Mode:                  insta ctf
	Map:                        garden
	Players:                    14
	Maximum Number of Players:  23
	Time Left (seconds):        262
	Protocol Version:           259

`GetPlayerInfo()` and `GetPlayerInfoRaw()` work pretty much the same; here is an example to get the player information of the player with the cn 4 on PSL1:

	...

	func main() {
		psl1 := extinfo.NewServer("sauerleague.org", 10000)
		playerInfo, err := psl1.GetPlayerInfo(14)
		if err != nil {
			fmt.Print("Error getting player information: ", err)
			return
		}

		fmt.Print("Player Information:\n")
		fmt.Printf("Name:                       %v\n", playerInfo.Name)
		fmt.Printf("Client Number:              %v\n", playerInfo.ClientNum)
		fmt.Printf("Ping:                       %v\n", playerInfo.Ping)
		fmt.Printf("Team:                       %v\n", playerInfo.Team)
		fmt.Printf("Frags:                      %v\n", playerInfo.Frags)
		// here you could get more things like deaths, health, armour, the player state (dead/alive/spectator/...), and so on
		fmt.Printf("Privilege:                  %v\n", playerInfo.Privilege)
		fmt.Printf("IP:                         %v\n", playerInfo.IP)
	}

Output would look like this:

	Player Information:
	Name:                       [tBMC]Rsn
	Client Number:              14
	Ping:                       45
	Team:                       good
	Frags:                      8
	Privilege:                  none
	IP:                         85.8.108.0

There is also `GetTeamsScores()` which returns all teams' scores (a TeamsScores containing a TeamScore (team score) for every team in the current game). Example code:

	...

	func main() {
		psl1 := extinfo.NewServer("sauerleague.org", 10000)
		scores, err := psl1.GetTeamsScores()
		if err != nil {
			fmt.Print("Error getting teams' scores: ", err)
			return
		}

		fmt.Print("Teams' Scores:\n")

		fmt.Printf("Game Mode:                  %v\n", scores.GameMode)
		fmt.Printf("Time Left (seconds):        %v\n", scores.SecsLeft)
		fmt.Print("Scores:\n")

		for _, score := range scores.Scores {
			fmt.Printf("   Team:                    %v\n", score.Name)
			fmt.Printf("   Score:                   %v\n", score.Score)

			if scores.GameMode == "capture" || scores.GameMode == "regen capture" {
				fmt.Printf("   Bases:                   %v\n", score.Bases)
			}
		}
	}

Output:

	Teams' Scores:
	Game Mode:                  insta ctf
	Time Left (seconds):        114
	Scores:
	   Team:                    good
	   Score:                   6
	   Team:                    evil
	   Score:                   4

More methods:

- `GetUptime()`: returns the amount of seconds the sauerbraten server is running
- `GetAllPlayerInfo()`: returns a PlayerInfo for every client connected to the server
- `GetTeamsScoresRaw()`: returns a TeamsScoresRaw containing a TeamScore for every team in the current game
