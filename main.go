package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var memory map[string]string

func handleconnection(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)

	if err != nil {
		log.Println(err)
	}
	input := strings.TrimSpace(string(buff[:n]))
	parts := strings.Fields(input)
	switch parts[0] {
	case "SET":
		set(parts[1], parts[2])
		fmt.Fprintf(conn, "ECHO - %s", "ok")
	case "GET":
		fmt.Fprintf(conn, "ECHO - %s", get(parts[1]))

	}

}

func get(key string) string {

	val, err := memory[key]

	if err == false {
		return "NOT FOUND"

	}

	return val

}
func set(key string, val string) {
	memory[key] = val
}
func main() {

	memory = make(map[string]string)

	server, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("server is listening on 8000")
	defer server.Close()
	for {
		connection, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("connected", connection)
		go handleconnection(connection)
	}
}
