package main

import (
	//"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/xtaci/kcp-go/v5"
	//"golang.org/x/crypto/pbkdf2"
)

func main() {
	if listener, err := kcp.ListenWithOptions("127.0.0.1:12345", nil, 0, 0); err == nil {
		// spin-up the client
		go client()
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}
}

// handleEcho send back everything it received
func handleEcho(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func client() {
	time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions("127.0.0.1:12345", nil, 0, 0); err == nil {
		data := time.Now().String()
		log.Println("sent kcp:", data)
		if _, err := sess.Write([]byte(data)); err == nil {
			fmt.Println("write first kcp")
		} else {
			log.Fatal(err)
		}

		for {
			data := time.Now().String()
			buf := make([]byte, len(data))
			// read back the data
			if n, err := sess.Read(buf); err == nil {
				log.Println("recv:", string(buf), n)
			} else {
				log.Fatal(err)
			}

			if _, err := sess.WritePacket(0, []byte(data)); err == nil {
				log.Println("sent udp:", data)
			} else {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Fatal(err)
	}
}
