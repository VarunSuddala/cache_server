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

	server, err := net.Listen("tcp", "0.0.0.0:8000")
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
			db.set(parts[1], parts[2])
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'GET' command\r\n"))
				continue
			}
			value, ok := db.get(parts[1])
			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				buf := make([]byte, 0, 64)
				buf = fmt.Appendf(buf, "*%d\r\n%s\r\n", len(value), value)
				conn.Write(buf)
			}
		case "DEL":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'DEL' command\r\n"))
				continue
			}
			delete := db.delete(parts[1])
			if delete {
				conn.Write([]byte(":1\r\n"))
			} else {
				conn.Write([]byte(":0\r\n"))
			}

		case "EXISTS":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'EXISTS' command\r\n"))
				continue
			}
			exists := db.exists(parts[1])
			if exists {
				conn.Write([]byte(":1\r\n"))
			} else {
				conn.Write([]byte(":0\r\n"))
			}

		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
func (d *Database) set(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.store[key] = value
}
func (d *Database) get(key string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	value, ok := d.store[key]
	return value, ok
}
func (d *Database) delete(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.store[key]; ok {
		delete(d.store, key)
		return true
	}
	return false
}

func (d *Database) exists(key string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, ok := d.store[key]
	return ok
}
