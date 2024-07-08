package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	udpServer, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	fmt.Println("UDP KV store server listening on port 8080...")

	db := map[string]string{}

	ch := make(chan int)

	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go response(udpServer, addr, buf[:n], db, ch)
		<-ch
	}

}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte, db map[string]string, ch chan int) {
	message := string(buf)

	fmt.Println("message:", message)

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
		db[key] = value
		fmt.Println("inserted:", key, "=", db[key])
	} else if message == "version" {
		udpServer.WriteTo([]byte("version=1"), addr)
	} else {
		key := message
		value := db[key]
		fmt.Println("retrieved:", key, "=", value)
		udpServer.WriteTo([]byte(key+"="+value), addr)
	}
	ch <- 1
}
