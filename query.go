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
	if infoType == InfoTypeExtended {
		request = append(request, infoType, extendedInfoType)

		// client stats has to include the clientNum (-1 for all)
		if extendedInfoType == ExtInfoTypeClientInfo {
			request = append(request, byte(clientNum))
		}
	}

	// basic info request
	if infoType == InfoTypeBasic {
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
	rawResponse := make([]byte, MaxPacketLength)
	var bytesRead int
	conn.SetReadDeadline(time.Now().Add(s.timeOut))
	bytesRead, err = bufconn.Read(rawResponse)
	if err != nil {
		return
	}

	// trim response to what's actually from the server
	packet := cubecode.NewPacket(rawResponse[:bytesRead])

	// response must include the entire request, ExtInfoAck, ExtInfoVersion, and either ExtInfoError or uptime
	if bytesRead < len(request)+3 {
		err = errors.New("extinfo: invalid response: too short")
		return
	}

	infoType, err := packet.ReadByte()
	if err != nil {
		return
	}

	// end of basic info response handling
	if infoType == InfoTypeBasic {
		response, err = packet.SubPacketFromRemaining()
		return
	}

	// make sure the entire request is correctly replayed

	for i, b := range request {
		if rawResponse[i] != b {
			err = errors.New("extinfo: invalid response: response does not match request")
			return
		}
	}

	command, err := packet.ReadByte()
	if err != nil {
		return
	}

	// skip rest of request
	for i := 2; i < len(request); i++ {
		_, err = packet.ReadByte()
		if err != nil {
			return
		}
	}

	// validate ack
	ack, err := packet.ReadByte()
	if err != nil {
		return
	}
	if ack != ExtInfoACK {
		err = errors.New("extinfo: invalid response: expected " + strconv.Itoa(int(ExtInfoACK)) + " (ACK), got " + strconv.Itoa(int(ack)))
		return
	}

	// validate version
	version, err := packet.ReadByte()
	if err != nil {
		return
	}
	// this package only supports protocol version 105
	if version != ExtInfoVersion {
		err = errors.New("extinfo: wrong version: expected " + strconv.Itoa(int(ExtInfoVersion)) + ", got " + strconv.Itoa(int(version)))
		return
	}

	// end of uptime request handling
	if command == ExtInfoTypeUptime {
		response, err = packet.SubPacketFromRemaining()
		return
	}

	commandError, err := packet.ReadByte()
	if err != nil {
		return
	}

	if commandError == ExtInfoError {
		switch command {
		case ExtInfoTypeClientInfo:
			err = errors.New("extinfo: no client with cn " + strconv.Itoa(int(request[2])))
		case ExtInfoTypeTeamScores:
			err = errors.New("extinfo: server is not running a team mode")
		}
		return
	}

	// end of team scores request handling
	if command == ExtInfoTypeTeamScores {
		response, err = packet.SubPacketFromRemaining()
		return
	}

	// handle response to ExtInfoTypeClientInfo

	// some server mods silently fail to implement responses → fail gracefully
	clientNumsHeader, err := packet.ReadByte()
	if err != nil {
		return
	}
	if clientNumsHeader != ClientInfoResponseTypeCNs {
		err = errors.New("extinfo: invalid response: expected " + strconv.Itoa(int(ExtInfoVersion)) + ", got " + strconv.Itoa(int(version)))
		return
	}

	// count (but discard) CNs
	numberOfClients := 0
	for packet.HasRemaining() {
		_, err = packet.ReadInt()
		if err != nil {
			return
		}

		numberOfClients++
	}

	// for each client, receive a packet and append it to a new slice
	clientInfos := make([]byte, 0, MaxPacketLength*numberOfClients)
	for i := 0; i < numberOfClients; i++ {
		// read from connection
		clientInfo := make([]byte, MaxPacketLength)
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
