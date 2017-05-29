package extinfo

import (
	"errors"

	"github.com/sauerbraten/cubecode"
)

// BasicInfoRaw contains the information sent back from the server in their raw form, i.e. no translation from ints to strings, even if possible.
type BasicInfoRaw struct {
	NumberOfClients    int    `json:"numberOfClients"`    // the number of clients currently connected to the server (players and spectators)
	ProtocolVersion    int    `json:"protocolVersion"`    // version number of the protocol in use by the server
	GameMode           int    `json:"gameMode"`           // current game mode
	SecsLeft           int    `json:"secsLeft"`           // the time left until intermission in seconds
	MaxNumberOfClients int    `json:"maxNumberOfClients"` // the maximum number of clients the server allows
	MasterMode         int    `json:"masterMode"`         // the current master mode of the server
	Paused             bool   `json:"paused"`             // wether the game is paused or not
	GameSpeed          int    `json:"gameSpeed"`          // the gamespeed
	Map                string `json:"map"`                // current map
	Description        string `json:"description"`        // server description
}

// BasicInfo contains the parsed information sent back from the server, i.e. game mode and master mode are translated into human readable strings.
type BasicInfo struct {
	BasicInfoRaw
	GameMode   string `json:"gameMode"`   // current game mode
	MasterMode string `json:"masterMode"` // the current master mode of the server
}

// GetBasicInfoRaw queries a Sauerbraten server at addr on port and returns the raw response or an error in case something went wrong. Raw response means that the int values sent as game mode and master mode are NOT translated into the human readable name.
func (s *Server) GetBasicInfoRaw() (basicInfoRaw BasicInfoRaw, err error) {
	var response *cubecode.Packet
	response, err = s.queryServer(buildRequest(InfoTypeBasic, 0, 0))
	if err != nil {
		return
	}

	basicInfoRaw.NumberOfClients, err = response.ReadInt()
	if err != nil {
		err = errors.New("extinfo: error reading number of connected clients: " + err.Error())
		return
	}

	// next int is always 5 or 7, the number of additional attributes after the clientcount and before the strings for map and description
	sevenAttributes := false
	numberOfAttributes, err := response.ReadInt()
	if err != nil {
		err = errors.New("extinfo: error reading number of following values: " + err.Error())
		return
	}

	if numberOfAttributes == 7 {
		sevenAttributes = true
	}

	basicInfoRaw.ProtocolVersion, err = response.ReadInt()
	if err != nil {
		err = errors.New("extinfo: error reading protocol version: " + err.Error())
		return
	}

	basicInfoRaw.GameMode, err = response.ReadInt()
	if err != nil {
		err = errors.New("extinfo: error reading game mode: " + err.Error())
		return
	}

	basicInfoRaw.SecsLeft, err = response.ReadInt()
	if err != nil {
		err = errors.New("extinfo: error reading time left: " + err.Error())
		return
	}

	basicInfoRaw.MaxNumberOfClients, err = response.ReadInt()
	if err != nil {
		err = errors.New("extinfo: error reading maximum number of clients: " + err.Error())
		return
	}

	basicInfoRaw.MasterMode, err = response.ReadInt()
	if err != nil {
		err = errors.New("extinfo: error reading master mode: " + err.Error())
		return
	}

	if sevenAttributes {
		var isPausedValue int
		isPausedValue, err = response.ReadInt()
		if err != nil {
			err = errors.New("extinfo: error reading paused value: " + err.Error())
			return
		}

		if isPausedValue == 1 {
			basicInfoRaw.Paused = true
		}

		basicInfoRaw.GameSpeed, err = response.ReadInt()
		if err != nil {
			err = errors.New("extinfo: error reading game speed: " + err.Error())
			return
		}
	} else {
		basicInfoRaw.GameSpeed = 100
	}

	basicInfoRaw.Map, err = response.ReadString()
	if err != nil {
		err = errors.New("extinfo: error reading map name: " + err.Error())
		return
	}

	basicInfoRaw.Description, err = response.ReadString()
	if err != nil {
		err = errors.New("extinfo: error reading server description: " + err.Error())
		return
	}

	return
}

// GetBasicInfo queries a Sauerbraten server at addr on port and returns the parsed response or an error in case something went wrong. Parsed response means that the int values sent as game mode and master mode are translated into the human readable name, e.g. '12' -> "insta ctf".
func (s *Server) GetBasicInfo() (BasicInfo, error) {
	basicInfo := BasicInfo{}

	basicInfoRaw, err := s.GetBasicInfoRaw()
	if err != nil {
		return basicInfo, err
	}

	basicInfo.BasicInfoRaw = basicInfoRaw
	basicInfo.GameMode = getGameModeName(basicInfo.BasicInfoRaw.GameMode)
	basicInfo.MasterMode = getMasterModeName(basicInfo.BasicInfoRaw.MasterMode)
	basicInfo.Map = cubecode.SanitizeString(basicInfo.Map)
	basicInfo.Description = cubecode.SanitizeString(basicInfo.Description)
	return basicInfo, nil
}
