package main

import (
	"net"
	"os"
	"time"
)

func main() {
	udpServer, err := net.ResolveUDPAddr("udp", ":8080")

	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}

	//close the connection
	defer conn.Close()

	_, err = conn.Write([]byte("a=1"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	// buffer to get data
	received := make([]byte, 1024)

	_, err = conn.Write([]byte("b=1"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	_, err = conn.Write([]byte("a=2"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	_, err = conn.Write([]byte("b"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println(string(received))

	_, err = conn.Write([]byte("a"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println(string(received))

	time.Sleep(1 * time.Second)

	_, err = conn.Write([]byte("version"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println(string(received))
}
