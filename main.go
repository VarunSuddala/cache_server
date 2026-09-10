package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Database struct {
	mu    sync.RWMutex
	store map[string]string
}

func new_database() *Database {
	return &Database{
		mu:    sync.RWMutex{},
		store: make(map[string]string),
	}
}

var db = new_database()

func main() {

	server, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("server is listening on 8000")
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		fmt.Println("client connected", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		b := make([]byte, 512)
		n, err := conn.Read(b)
		if err != nil {
			fmt.Println("client disconnected:", err)
			return
		}
		cmd := string(b[:n])
		fmt.Println("Received command:", cmd)
		parts := strings.Split(strings.TrimSpace(cmd), " ")

		switch strings.ToUpper(parts[0]) {
		case "SET":
			if len(parts) != 3 {
				conn.Write([]byte("-ERR wrong number of arguments for 'SET' command\r\n"))
				continue
			}
			db.mu.Lock()
			set(parts[1], parts[2])
			db.mu.Unlock()
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'GET' command\r\n"))
				continue
			}
			db.mu.RLock()
			value, ok := get(parts[1])
			db.mu.RUnlock()
			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
			}

		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}

func set(key, value string) {
	db.store[key] = value
}

func get(key string) (string, bool) {
	value, ok := db.store[key]
	if !ok {
		return "", false
	}
	return value, true
}
