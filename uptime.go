package extinfo

// GetUptime returns the uptime of the server in seconds.
func (s *Server) GetUptime() (uptime int, err error) {
	var response *extinfoResponse
	response, err = s.queryServer(buildRequest(EXTENDED_INFO, EXTENDED_INFO_UPTIME, 0))
	if err != nil {
		return
	}

	uptime, err = response.ReadInt()

	return
}
