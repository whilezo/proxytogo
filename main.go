package main

import (
	"log"
	"net"
)

func main() {
	server, err := net.Listen("tcp", "0.0.0.0:4040")
	checkError(err)

	log.Println("server created")

	for {
		conn, err := server.Accept()
		checkError(err)

		log.Printf("accepted client")

		handleClient(conn)
	}
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	bufer := make([]byte, 4096)

	for {
		_, err := conn.Read(bufer)
		checkError(err)
		log.Printf("Bytes: %d", len(bufer))

		conn.Write(bufer)
		log.Print("connection closed")
		break
	}
}
