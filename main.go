package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

var memory map[string]string

type ValueType string

const (
	ARRAY  ValueType = "*"
	BULK   ValueType = "$"
	STRING ValueType = "+"
)

type Value struct {
	typ   ValueType
	bulk  string
	str   string
	array []Value
}

func (v *Value) readValue(reader io.Reader) {
	buff := make([]byte, 4)
	reader.Read(buff)
	arrLen, err := strconv.Atoi(string(buff[0]))
	if err != nil {
		log.Println(err)
		return
	}
	for range arrLen {
		bulk := v.readBulk(reader)
		v.array = append(v.array, bulk)
	}
}

func (v *Value) readBulk(reader io.Reader) Value {
	buff := make([]byte, 4)
	reader.Read(buff)
	n, err := strconv.Atoi(string(buff[1]))

	if err != nil {
		log.Print(err)
		return Value{}
	}

	bulkBuff := make([]byte, n)
	reader.Read(bulkBuff)

	bulk := string(bulkBuff[:n])

	return Value{typ: BULK, bulk: bulk}
}

func main() {

	memory = make(map[string]string)

	server, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("server is listening on 8000")
	defer server.Close()

	conn, err := server.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	for {
		v := Value{typ: ARRAY}
		v.readValue(conn)
		fmt.Println(v.array)
		conn.Write([]byte("+ok\r\n"))
		return
	}
}
