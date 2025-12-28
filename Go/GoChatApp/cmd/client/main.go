package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr := "127.0.0.1:9000"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()
	log.Printf("connected to %s", addr)

	stdin := bufio.NewReader(os.Stdin)
	server := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		line, err := stdin.ReadString('\n')
		if err != nil {
			log.Fatalf("stdin read error: %v", err)
		}
		if _, err := conn.Write([]byte(line)); err != nil {
			log.Printf("write error : %v", err)
			return
		}

		reply, err := server.ReadString('\n')
		if err != nil {
			log.Fatalf("server read error : %v", err)
		}
		fmt.Print(reply)
	}
}
