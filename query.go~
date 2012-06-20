package extinfo

import (
	"net"
	"bufio"
)

// builds a request
func buildRequest(informationType int, extendedInformationType int, clientNum int) []byte {
	request := make([]byte, 0)

	// extended info request
	if informationType == 0 {
		// player stats has to include the clientNum
		if extendedInformationType == 1 {
			request = append(request, byte(informationType), byte(extendedInformationType), byte(clientNum))
		} else {
			request = append(request, byte(informationType), byte(extendedInformationType))
		}
	}

	// basic info request
	if informationType == 1 {
		request = append(request, byte(informationType))
	}

	return request
}

// queries the given server and returns the raw response as []byte and an error in case something went wrong. clientNum is optional, put 0 if not needed.
func queryServer(addr string, port int, request []byte) ([]byte, error) {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return []byte{}, err
	}

	// connect to server at port+1 (port is the port you connect to in game, sauerbraten listens on the one higher port for BasicInfo queries
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{ipaddr.IP, port+1})
	if err != nil {
		return []byte{}, err
	}

	// set up a buffered reader
	bufconn := bufio.NewReader(conn)

	// send '1' to server
	_, err = conn.Write(request)
	if err != nil {
		return []byte{}, err
	}

	// receive response from server
	response := make([]byte, 64)
	_, err = bufconn.Read(response)
	if err != nil {
		return []byte{}, err
	}

	if response[0] == 0x00 && response[1] == 0x01 && response[2] != 0xFF {
		playerStatsResponse := make([]byte, 64)
		_, err = bufconn.Read(playerStatsResponse)
		if err != nil {
			return []byte{}, err
		}
		return playerStatsResponse, nil
	}
	return response, nil
}
