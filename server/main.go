package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type request struct {
	addr net.Addr
	data []byte
}

func main() {
	udpServer, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	fmt.Println("UDP KV store server listening on port 8080...")

	db := map[string]string{}

	var mu sync.Mutex

	ch := make(chan request)

	// for i := 0; i < 5; i++ {
	// 	go worker(udpServer, db, &mu, ch)
	// }
	go worker(udpServer, db, &mu, ch) // single worker thread to process requests sequentially

	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		// go response(udpServer, addr, buf[:n], db, ch)
		ch <- request{addr: addr, data: buf[:n]}
	}
}

func worker(udpServer net.PacketConn, db map[string]string, mu *sync.Mutex, ch chan request) {
	for request := range ch {
		response(udpServer, request.addr, request.data, db, mu)
	}
}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte, db map[string]string, mu *sync.Mutex) {
	message := string(buf)

	fmt.Println("message:", message)

	equalIndex := -1
	for i, c := range message {
		if c == '=' {
			equalIndex = i
			break
		}
	}

	// not need for locking if single threaded
	// mu.Lock()
	// defer mu.Unlock()

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
}
