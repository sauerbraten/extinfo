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

Detailed Documentation [here](http://go.pkgdoc.org/github.com/sauerbraten/extinfo).

## Example

Here is code to get you the state of the PSL1 server:

	package main
	
	import (
		"fmt"
		"github.com/sauerbraten/extinfo"
	)
	
	func main() {
		pslBasicInfo, err := extinfo.GetBasicInfo("sauerleague.org", 10000)
		if err != nil {
			fmt.Print("Error getting basic information: ", err)
			return
		}
	
		fmt.Print("Basic Server Information:\n")
		fmt.Printf("Description:\t\t\t%v\n", pslBasicInfo.Description)
		fmt.Printf("Master Mode:\t\t\t%v\n", pslBasicInfo.MasterMode)
		fmt.Printf("Game Mode:\t\t\t%v\n", pslBasicInfo.GameMode)
		fmt.Printf("Map:\t\t\t\t%v\n", pslBasicInfo.Map)
		fmt.Printf("Players:\t\t\t%v\n", pslBasicInfo.NumberOfPlayers)
		fmt.Printf("Maximum Number of Players:\t%v\n", 	pslBasicInfo.MaxNumberOfPlayers)
		fmt.Printf("Time Left (seconds):\t\t%v\n", pslBasicInfo.SecsLeft)
		fmt.Printf("Protocol Version:\t\t%v\n", pslBasicInfo.ProtocolVersion)
	}

The output should be something like this:

	Basic Server Information:
	Description:				PSL.sauerleague.org
	Master Mode:				auth
	Game Mode:					insta ctf
	Map:						damnation
	Players:					20
	Maximum Number of Players:	20
	Time Left (seconds):		17
	Protocol Version:			258

`GetPlayerInfo()` and `GetPlayerInfoRaw()` work pretty much the same; here is an example to get the player information of the player with the cn 4 on PSL1:

	package main
	
	import (
		"fmt"
		"github.com/sauerbraten/extinfo"
	)
	
	func main() {
		playerInfo, err := extinfo.GetPlayerInfo("sauerleague.org", 10000, 4)
		if err != nil {
			fmt.Print("Error getting player information: ", err)
			return
		}
	
		fmt.Print("Player Information:\n")
		fmt.Printf("Name:\t\t\t\t%v\n", playerInfo.Name)
		fmt.Printf("Client Number:\t\t\t%v\n", playerInfo.ClientNum)	
		fmt.Printf("Ping:\t\t\t\t%v\n", playerInfo.Ping)
		fmt.Printf("Team:\t\t\t\t%v\n", playerInfo.Team)
		fmt.Printf("Frags:\t\t\t\t%v\n", playerInfo.Frags)
		// here you could get more things like deaths, health, armour, the player state (dead/alive/spectator/...), and so on
		fmt.Printf("Privilege:\t\t\t%v\n", playerInfo.Privilege)
		fmt.Printf("IP:\t\t\t\t\t%v\n", playerInfo.IP)
	}

Output would look like this:

	Player Information:
	Name:				oo|berk
	Client Number:		4
	Ping:				25
	Team:				evil
	Frags:				37
	Privilege:			none
	IP:					123.234.345.0
