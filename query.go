// +build go1.1

package extinfo

import (
	"bufio"
	"errors"
	"log"
	"net"
	"time"
)

// builds a request
func buildRequest(infoType int, extendedInfoType int, clientNum int) []byte {
	request := []byte{}

	// extended info request
	if infoType == EXTENDED_INFO {
		request = append(request, byte(infoType), byte(extendedInfoType))

		// player stats has to include the clientNum
		if extendedInfoType == EXTENDED_INFO_CLIENT_INFO {
			request = append(request, byte(clientNum))
		}
	}

	// basic info request
	if infoType == BASIC_INFO {
		request = append(request, byte(infoType))
	}

	return request
}

// queries the given server and returns the raw response as []byte and an error in case something went wrong. clientNum is optional, put 0 if not needed.
func (s *Server) queryServer(request []byte) ([]byte, error) {
	// connect to server at port+1 (port is the port you connect to in game, sauerbraten listens on the one higher port for BasicInfo queries
	conn, err := net.DialUDP("udp", nil, s.addr)
	defer conn.Close()
	if err != nil {
		return []byte{}, err
	}

	// set up a buffered reader
	bufconn := bufio.NewReader(conn)

	log.Println(request)

	// send the request to server
	_, err = conn.Write(request)
	if err != nil {
		return []byte{}, err
	}

	// receive response from server with 5 second timeout
	response := make([]byte, 64)
	var bytesRead int
	conn.SetReadDeadline(time.Now().Add(s.timeOut))
	bytesRead, err = bufconn.Read(response)
	if err != nil {
		return []byte{}, err
	}

	// trim response to what's actually from the server
	response = response[:bytesRead]

	if bytesRead < 2 {
		return []byte{}, errors.New("extinfo: invalid response")
	}

	// if not a response to EXTENDED_INFO_CLIENT_INFO, we are done
	if response[0] != EXTENDED_INFO || (response[0] == EXTENDED_INFO && response[1] != EXTENDED_INFO_CLIENT_INFO) {
		return response, nil
	}

	// handle response to EXTENDED_INFO_CLIENT_INFO

	// some server mods silently fail to implement responses → fail gracefully
	if len(response) < 7 {
		return []byte{}, errors.New("extinfo: invalid response")
	}

	if response[5] == EXTENDED_INFO_ERROR {
		return []byte{}, errors.New("extinfo: invalid cn")
	}

	if response[6] != EXTENDED_INFO_CLIENT_INFO_RESPONSE_CNS {
		return []byte{}, errors.New("extinfo: invalid response")
	}

	// get CNs out of the reponse, ignore 7 first bytes, which are:
	// EXTENDED_INFO, EXTENDED_INFO_CLIENT_INFO, CN from request, EXTENDED_INFO_ACK, EXTENDED_INFO_VERSION, EXTENDED_INFO_NO_ERROR, EXTENDED_INFO_CLIENT_INFO_RESPONSE_CNS
	clientNums := response[7:]

	log.Println("clientNums:", clientNums)

	// for each client, receive a packet and append it to the response
	clientInfos := make([]byte, 0)
	for _ = range clientNums {
		// read from connection
		clientInfo := make([]byte, 64)
		conn.SetReadDeadline(time.Now().Add(s.timeOut))
		_, err = bufconn.Read(clientInfo)
		if err != nil {
			return clientInfos, err
		}

		// append to slice
		clientInfos = append(clientInfos, clientInfo...)
	}

	return clientInfos, nil
}
