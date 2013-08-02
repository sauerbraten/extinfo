// +build !go1.1

package extinfo

import (
	"bufio"
	"bytes"
	"net"
)

// builds a request
func buildRequest(informationType int, extendedInfoType int, clientNum int) []byte {
	request := make([]byte, 0)

	// extended info request
	if informationType == extendedInfo {
		// player stats has to include the clientNum
		if extendedInfoType == playerStatsInfo {
			request = append(request, byte(informationType), byte(extendedInfoType), byte(clientNum))
		} else {
			request = append(request, byte(informationType), byte(extendedInfoType))
		}
	}

	// basic info request
	if informationType == basicInfo {
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
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{ipaddr.IP, port + 1})
	defer conn.Close()
	if err != nil {
		return []byte{}, err
	}

	// set up a buffered reader
	bufconn := bufio.NewReader(conn)

	// send the request to server
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

	// first byte = 0, second byte = 1, 4th byte 0 (no error) --> player info response with no error, wait for following packages
	if response[0] == 0x00 && response[1] == 0x01 && response[5] == 0x00 {
		// if third byte = -1, information was queried for all players --> multiple packages following
		if response[2] == 0xFF {
			// trim null bytes
			response = bytes.TrimRight(response, "\x00")

			// get player cns out of the reponse: 7 first bytes are extendedInfo, playerStatsInfo, clientNum, server ACK byte, server VERSION byte, server NO_ERROR byte, server playerStatsInfo_RESP_STATS byte
			playerCns := response[7:]

			// for each client, receive a packet and append it to the response
			playerInfos := make([]byte, 0)
			for _ = range playerCns {
				// read from connection
				response = make([]byte, 64)
				_, err = bufconn.Read(response)

				// append to slice
				playerInfos = append(playerInfos, response...)

				// on error, return what we already have
				if err != nil {
					return playerInfos, err
				}
			}

			return playerInfos, nil
		}

		// else, only one cn was asked for --> one package following
		playerInfoResponse := make([]byte, 64)
		_, err = bufconn.Read(playerInfoResponse)
		if err != nil {
			return []byte{}, err
		}

		return playerInfoResponse, nil
	}

	return response, nil
}
