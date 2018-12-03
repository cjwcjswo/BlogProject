package main

import (
	"fmt"
	"net"
	"log"
	"encoding/json"
	"time"
)

type Packet struct {
	Action 			string 	`json:"Action"`
	Message 		string 	`json:"Message"`
}

func main() {

	var serverAddr string
	fmt.Print("Enter HolePunching Server Address(ex: 192.168.0.3:243): ")
	fmt.Scanln(&serverAddr)

	serverResolveAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	clientAddr := conn.LocalAddr().String()

	sendMessage := Packet{
		Action:  "New",
		Message: clientAddr,
	}
	sendByte, err := json.Marshal(&sendMessage)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteToUDP(sendByte, serverResolveAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Local Address Send Success!: ", clientAddr)

	readPacket(conn, serverResolveAddr)
}

func readPacket(conn *net.UDPConn, serverResolveAddr *net.UDPAddr) {
	for {
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer[0:])
		if err != nil {
			log.Fatal(err)
		}
		message := buffer[:length]
		var packet Packet
		err = json.Unmarshal(message, &packet)
		if err != nil {
			log.Fatal(err)
		}
		processPacket(conn, serverResolveAddr, packet.Action, packet.Message)

	}
}

func processPacket(conn *net.UDPConn, serverResolveAddr *net.UDPAddr, protocol string, message string) {
	switch protocol {
	case "New" : {
		log.Println("My Public Addr: ", string(message))
		sendFindMessage(conn, serverResolveAddr)
	}
	case "Find" : {
		if message == "" {
			time.Sleep(5 * time.Second)
			sendFindMessage(conn, serverResolveAddr)
		} else {
			otherClientAddr, err := net.ResolveUDPAddr("udp", message)
			if err != nil {
				log.Fatal(err.Error())
			}
			go chattingProcess(conn, otherClientAddr)
		}
	}
	case "Message" : {
		fmt.Println("OtherClientMessage: ", string(message))
	}
	}
}

func sendFindMessage(conn *net.UDPConn, serverResolveAddr *net.UDPAddr) {
	log.Println("Find Another Client...")
	sendMessage := Packet{
		Action:  "Find",
	}
	sendByte, err := json.Marshal(sendMessage)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteToUDP(sendByte, serverResolveAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func chattingProcess(conn *net.UDPConn, otherClientAddr *net.UDPAddr) {
	for {
		var chat string
		fmt.Printf("Enter SendMessage: ")
		fmt.Scanln(&chat)

		packet := Packet{
			Action:  "Message",
			Message: chat,
		}
		packetData, err := json.Marshal(packet)
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = conn.WriteToUDP(packetData, otherClientAddr)
		if err != nil {
			log.Fatal(err.Error())
		}

	}
}
