package server

import (
	"io"
	"log"
	"net"

	"github.com/VarunSuddala/redis/internal/config"
)

func readCommand(c net.Conn) (string, error) {

	var buf []byte = make([]byte, 512)

	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil

}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func RunSyncTCPServer() {
	cfg := config.MustLoad()
	log.Println("starting a TCP server on ", cfg.HOST+":"+cfg.PORT)
	var con_clients int = 0

	lsnr, err := net.Listen("tcp", cfg.HOST+":"+cfg.PORT)
	if err != nil {
		panic(err)
	}
	for {
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		con_clients += 1
		log.Println("client connected with address :", c.RemoteAddr(), "concurrent clients", con_clients)

		for {
			cmd, err := readCommand(c)

			if err != nil {

				log.Println("client disconnected", c.RemoteAddr(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				} else {
					log.Println("error", err)
				}
				con_clients -= 1
				c.Close()
				break

			}
			log.Println("command", cmd)

			if err = respond(cmd, c); err != nil {
				log.Println("err write", err)
			}

		}
	}
}
