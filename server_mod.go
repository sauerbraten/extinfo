package extinfo

import "github.com/sauerbraten/cubecode"

// GetServerMod returns the name of the mod in use at this server.
func (s *Server) GetServerMod() (serverMod string, err error) {
	var response *cubecode.Packet
	uptimeRequest := buildRequest(InfoTypeExtended, ExtInfoTypeUptime, 0)
	modRequest := append(uptimeRequest, 0x01)
	response, err = s.queryServer(modRequest)
	if err != nil {
		return
	}

	// read & discard uptime
	_, err = response.ReadInt()
	if err != nil {
		return
	}

	// try to read one more byte
	mod, err := response.ReadInt()

	// if there is none, it's not a detectable mod (probably vanilla), so we will return ""
	if err == cubecode.ErrBufferTooShort {
		err = nil
	} else if err == nil {
		serverMod = getServerModName(mod)
	}

	return
}
