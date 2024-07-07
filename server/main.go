package main

import (
	"log"
	"net"
)

func main() {
	udpServer, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	db := map[string]string{}

	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go response(udpServer, addr, buf[:n], &db)
	}

}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte, db *map[string]string) {
	message := string(buf)

	equalIndex := -1
	for i, c := range message {
		if c == '=' {
			equalIndex = i
			break
		}
	}

	if equalIndex >= 0 {
		key := message[:equalIndex]
		value := message[equalIndex+1:]
		(*db)[key] = value
	} else if message == "version" {
		udpServer.WriteTo([]byte("version=1"), addr)
	} else {
		key := message
		value := (*db)[key]
		udpServer.WriteTo([]byte(key+"="+value), addr)
	}

}
