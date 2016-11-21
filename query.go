package extinfo

import (
	"bufio"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/sauerbraten/cubecode"
)

// builds a request
func buildRequest(infoType byte, extendedInfoType byte, clientNum int) []byte {
	request := []byte{}

	// extended info request
	if infoType == EXTENDED_INFO {
		request = append(request, infoType, extendedInfoType)

		// client stats has to include the clientNum (-1 for all)
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

// queries the given server and returns the response and an error in case something went wrong. clientNum is optional, put 0 if not needed.
func (s *Server) queryServer(request []byte) (response *cubecode.Packet, err error) {
	// connect to server at port+1 (port is the port you connect to in game, sauerbraten listens on the one higher port for BasicInfo queries
	var conn *net.UDPConn
	conn, err = net.DialUDP("udp", nil, s.addr)
	if err != nil {
		return
	}
	defer conn.Close()

	// set up a buffered reader
	bufconn := bufio.NewReader(conn)

	// send the request to server
	_, err = conn.Write(request)
	if err != nil {
		return
	}

	// receive response from server with 5 second timeout
	rawResponse := make([]byte, MAX_PACKET_SIZE)
	var bytesRead int
	conn.SetReadDeadline(time.Now().Add(s.timeOut))
	bytesRead, err = bufconn.Read(rawResponse)
	if err != nil {
		return
	}

	// trim response to what's actually from the server
	packet := cubecode.NewPacket(rawResponse[:bytesRead])

	if bytesRead < 5 {
		err = errors.New("extinfo: invalid response: too short")
		return
	}

	// do some basic checks on the response

	infoType, err := packet.ReadByte()
	if err != nil {
		return
	}
	command, err := packet.ReadByte() // only valid if infoType == EXTENDED_INFO
	if err != nil {
		return
	}

	if infoType == EXTENDED_INFO {
		var version, commandError byte

		if command == EXTENDED_INFO_CLIENT_INFO {
			if bytesRead < 6 {
				err = errors.New("extinfo: invalid response: too short")
				return
			}
			version = rawResponse[4]
			commandError = rawResponse[5]
		} else {
			version = rawResponse[3]
			commandError = rawResponse[4]
		}

		if infoType != request[0] || command != request[1] {
			err = errors.New("extinfo: invalid response: response does not match request")
			return
		}

		// this package only support extinfo protocol version 105
		if version != EXTENDED_INFO_VERSION {
			err = errors.New("extinfo: wrong version: expected " + strconv.Itoa(int(EXTENDED_INFO_VERSION)) + ", got " + strconv.Itoa(int(version)))
			return
		}

		if commandError == EXTENDED_INFO_ERROR {
			switch command {
			case EXTENDED_INFO_CLIENT_INFO:
				err = errors.New("extinfo: no client with cn " + strconv.Itoa(int(request[2])))
			case EXTENDED_INFO_TEAMS_SCORES:
				err = errors.New("extinfo: server is not running a team mode")
			}
			return
		}
	}

	// if not a response to EXTENDED_INFO_CLIENT_INFO, we are done
	if infoType != EXTENDED_INFO || command != EXTENDED_INFO_CLIENT_INFO {
		offset := 0

		if infoType == EXTENDED_INFO {
			switch command {
			case EXTENDED_INFO_UPTIME:
				offset = 4
			case EXTENDED_INFO_TEAMS_SCORES:
				offset = 5
			}
		}

		response = cubecode.NewPacket(rawResponse[offset:])
		return
	}

	// handle response to EXTENDED_INFO_CLIENT_INFO

	// some server mods silently fail to implement responses → fail gracefully
	if len(rawResponse) < 7 || rawResponse[6] != EXTENDED_INFO_CLIENT_INFO_RESPONSE_CNS {
		err = errors.New("extinfo: invalid response")
		return
	}

	// get CNs out of the reponse, ignore 7 first bytes, which were:
	// EXTENDED_INFO, EXTENDED_INFO_CLIENT_INFO, CN from request, EXTENDED_INFO_ACK, EXTENDED_INFO_VERSION, EXTENDED_INFO_NO_ERROR, EXTENDED_INFO_CLIENT_INFO_RESPONSE_CNS
	clientNums := rawResponse[7:]

	numberOfClients, err := countClientNums(clientNums)
	if err != nil {
		return
	}

	// for each client, receive a packet and append it to a new slice
	clientInfos := make([]byte, 0, MAX_PACKET_SIZE*numberOfClients)
	for i := 0; i < numberOfClients; i++ {
		// read from connection
		clientInfo := make([]byte, MAX_PACKET_SIZE)
		conn.SetReadDeadline(time.Now().Add(s.timeOut))
		_, err = bufconn.Read(clientInfo)
		if err != nil {
			return
		}

		// append bytes to slice
		clientInfos = append(clientInfos, clientInfo...)
	}

	response = cubecode.NewPacket(clientInfos)
	return
}

func countClientNums(buf []byte) (count int, err error) {
	p := cubecode.NewPacket(buf)

	for p.HasRemaining() {
		_, err = p.ReadInt()
		if err != nil {
			return
		}

		count++
	}

	return
}
