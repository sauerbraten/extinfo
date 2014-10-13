package extinfo

// GetUptime returns the uptime of the server in seconds.
func (s *Server) GetUptime() (uptime int, err error) {
	var response []byte
	response, err = s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_UPTIME, 0))
	if err != nil {
		return -1, err
	}

	positionInResponse = 0

	// first int is EXTENDED_INFO
	_, err = dumpInt(response)
	if err != nil {
		return
	}

	// next int is EXT_EXTENDED_INFO_UPTIME = 0
	_, err = dumpInt(response)
	if err != nil {
		return
	}

	// next int is EXT_ACK = -1
	_, err = dumpInt(response)
	if err != nil {
		return
	}

	// next int is EXT_VERSION
	_, err = dumpInt(response)
	if err != nil {
		return
	}

	// next int is the actual uptime
	uptime, err = dumpInt(response)
	if err != nil {
		return
	}
	return
}
