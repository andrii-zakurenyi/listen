package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const usage = `Usage:

	listen <port number>

Example:

	listen 5678
`

const tmplt = `
═══ %s %s %d bytes ═══
%s
═══ END

`

const (
	network = "tcp"
	minPort = 1
	maxPort = 65535
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	port := os.Args[1]
	if n, err := strconv.Atoi(port); err != nil || n < minPort || n > maxPort {
		fmt.Println("Error: port must be a number ranging from 1 to 65535.")
		fmt.Println(usage)
		os.Exit(1)
	}

	fmt.Printf("Starting TCP listener: 127.0.0.1:%s\n\n", port)

	addr, err := net.ResolveTCPAddr(network, fmt.Sprintf("127.0.0.1:%s", port))
	if err != nil {
		fmt.Printf("ResolveTCPAddr: %v\n", err)
		os.Exit(1)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Printf("ListenTCP: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	buf := make([]byte, 1024*4)
	for {
		fmt.Println("Accepting new connections")
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Printf("AcceptTCP: %v\n", err)
			conn.Close()
			continue
		}
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Read: %v\n", err)
			conn.Close()
			continue
		}
		result := string(buf[:n])
		fmt.Printf(tmplt, time.Now().Format("15:04:05.999"), conn.RemoteAddr().String(), n, result)
		conn.Close()
	}
}
