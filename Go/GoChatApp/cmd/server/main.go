package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	addr := "127.0.0.1:9000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listening on %s", addr)
	}
	defer ln.Close()
	log.Printf("listening on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("accept error %v", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	log.Printf("connected: &s", remote)

	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("read error (%s): %v", remote, err)
			}
			log.Printf("disconnected: %s", remote)
			return
		}

		reply := fmt.Sprintf("echo: %s", line)
		if _, err := conn.Write([]byte(reply)); err != nil {
			log.Printf("write error (%s): %v", remote, err)
			return
		}
	}
}
