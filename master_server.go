package extinfo

import (
        "net"
        "time"
        "fmt"
        "strings"
)

const (
        MASTER_SERVER_ADDRESS = "sauerbraten.org"
        MASTER_SERVER_PORT = 28787
)

func GetMasterServerList(defaultTimeOut time.Duration) ([]*Server, error){
     serverList := make([]*Server, 100)
     bufList, err := queryMasterList()
     if err != nil{
          return serverList, err
     }

     for {
          status, readError := bufList.ReadString('\n')
          if readError != nil{
               return serverList, readError
          }
          parts := strings.Split(status, " ")
          serverIp := parts[1] + ":" + parts[2]
          address, err := net.ResolveUDPAddr("udp", serverIp)
          if err != nil{
               //Don't make one bad server parsing fatal
               fmt.Printf("Could not parse server from string: %s", serverIp)
               fmt.Printf("Got address: %v", address)
               return serverList, err
          }
          serverList = append(serverList, NewServer(address, defaultTimeOut))
     }

     return serverList, nil
}
