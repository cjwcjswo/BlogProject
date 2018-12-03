package main

import (
	"net"
	"fmt"
	"log"
	"encoding/json"
)

type Packet struct {
	Action 			string 	`json:"Action"`
	Message 		string 	`json:"Message"`
}

var clientInfoMap map[string]string // userName, publicAddr

func main() {
	clientInfoMap = make(map[string]string, 128)

	var bindIPAddress string
	fmt.Print("Enter bind address(ex: 192.168.0.3:243): ")
	fmt.Scanln(&bindIPAddress)

	log.Println("Start Server: ", bindIPAddress)
	addr, err := net.ResolveUDPAddr("udp", bindIPAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer udpConn.Close()

	for {
		var buffer = make([]byte, 1024)
		length, client, err:= udpConn.ReadFromUDP(buffer[:])
		if err != nil {
			log.Println("Read Error!", err)
			continue
		}
		if length < 1 {
			log.Println("Read Length Error!")
			continue
		}
		var packet Packet
		err = json.Unmarshal(buffer[:length], &packet)
		if err != nil {
			log.Println(err)
			continue
		}
		processClientPacket(udpConn, client, packet.Action, packet.Message)
	}
}

func processClientPacket(conn *net.UDPConn, client *net.UDPAddr, protocol string, message string) {
	log.Println("Read Packet! Protocol:", protocol, "Message: ", message)
	switch protocol {
		case "New":{
			userName := message
			clientPublicAddr := fmt.Sprintf("%s:%d", client.IP.String(), client.Port)
			clientInfoMap[userName] = clientPublicAddr

			packet := Packet{
				Action:  "New",
				Message: clientPublicAddr,
			}
			sendByte, err := json.Marshal(packet)
			if err != nil {
				log.Println(err)
				return
			}
			conn.WriteToUDP(sendByte, client)

			printClientInfoMap()
		}
		case "Find": {
			clientPublicAddr := fmt.Sprintf("%s:%d", client.IP.String(), client.Port)
			packet := Packet{
				Action:  "Find",
			}
			for _, publicAddr := range clientInfoMap {
				if publicAddr != clientPublicAddr {
					packet.Message = publicAddr
					break
				}
			}
			sendByte, err := json.Marshal(packet)
			if err != nil {
				log.Println(err)
				return
			}
			conn.WriteToUDP(sendByte, client)
		}
	}
}

func printClientInfoMap() {
	fmt.Println("----------Current Client Info----------")
		fmt.Println("UserName\t\tPublicAddr")
	for userName, publicAddr:= range clientInfoMap {
		fmt.Printf("%s\t\t%s\t\t", userName, publicAddr)
		fmt.Println()
	}
	fmt.Println("----------End	Client Info----------")
}