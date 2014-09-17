package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// Listen on TCP port 12300 on all interfaces
	l, err := net.Listen("tcp", ":12300")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new routine
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			defer c.Close()

			// Generate a random key
			key := make([]byte, 24)
			rand.Read(key)
			if err != nil {
				log.Fatal(err)
				return
			}
			scrambled_key := Scramble(key)
			log.Printf("Key: %s", hex.EncodeToString(key))
			log.Printf("Scrambled key: %s", hex.EncodeToString(scrambled_key))
			// Send key to alarm system
			n, err := c.Write(scrambled_key)
			// TODO: compare n with size of key
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Printf("Sent %d bytes to alarm (key)", n)
			buf := make([]byte, 1024)
			n, err = c.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Read error: ", err)
					return
				}
			}
			log.Printf("Read %d bytes", n)
			encryptedData := buf[:n]
			log.Printf("Data: %s", hex.EncodeToString(encryptedData))
			data := Decrypt3DESECB(encryptedData, key)
			fmt.Println("Message(byte): ", hex.EncodeToString(data))
			fmt.Println("Message: ", string(data[:]))
			ack := []byte("ACK")
			ack = append(ack, make([]byte, 5)...) // padding.
			encryptedAck := Encrypt3DESECB(ack, key)
			log.Printf("Encrypted ACK: %s", hex.EncodeToString(encryptedAck))
			n, err = c.Write(encryptedAck)
			if err != nil {
				log.Fatal(err)
			}

		}(conn)
	}
}
