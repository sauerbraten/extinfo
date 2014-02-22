package extinfo

// GetUptime returns the uptime of the server in seconds.
func (s *Server) GetUptime() (int, error) {
	response, err := s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_UPTIME, 0))
	if err != nil {
		return -1, err
	}

	positionInResponse = 0

	// first int is EXTENDED_INFO
	_ = dumpInt(response)

	// next int is EXT_EXTENDED_INFO_UPTIME = 0
	_ = dumpInt(response)

	// next int is EXT_ACK = -1
	_ = dumpInt(response)

	// next int is EXT_VERSION
	_ = dumpInt(response)

	// next int is the actual uptime
	return dumpInt(response), nil
}
