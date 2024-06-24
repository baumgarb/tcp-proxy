package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/baumgarb/tcp-proxy/vlog"
)

func main() {

	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Println("Usage: ./tcp-proxy <src addr> <dest addr> [-v]")
		fmt.Println()
		fmt.Println("Examples: ")
		fmt.Println("       ./tcp-proxy :3000 :5000                         # forwards all incoming TCP connections on port 3000 to 5000.")
		fmt.Println("       ./tcp-proxy :3000 :5000 &                       # forwards all incoming TCP connections on port 3000 to 5000")
		fmt.Println("                                                       # while running silently in the background.")
		fmt.Println("       ./tcp-proxy :5000 :31774 -v                     # forwards all incoming TCP connections on port 5000 to 31774.")
		fmt.Println("                                                       # Verbose logging is enabled.")
		fmt.Println("       ./tcp-proxy 127.0.0.1:3000 192.168.100.10:5000  # forwards all incoming TCP connections on 127.0.0.1 on port 3000")
		fmt.Println("                                                       # to target with IP 192.168.100.10 on port 5000.")
		os.Exit(1)
	}

	portA := os.Args[1]
	portB := os.Args[2]

	listener, err := net.Listen("tcp", portA)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	vlog.Printf("Proxy listening on port %v and forwarding to %v...\n", portA, portB)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		vlog.Printf("Accepted new connection on %v\n", portA)

		// Handle each connection in a separate goroutine
		go handleConnection(conn, portB)
	}
}

func handleConnection(clientConn net.Conn, portB string) {
	data := make([]byte, 256)
	n, err := clientConn.Read(data)
	if err != nil {
		log.Printf("Error peeking into client connection data: %v", err)
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	vlog.Printf("First line: %v", lines[0])

	// Connect to the target server (port B)
	serverConn, err := net.Dial("tcp", fmt.Sprintf("localhost%s", portB))
	if err != nil {
		log.Println("Error connecting to server:", err)
		clientConn.Close()
		return
	}
	defer func() {
		vlog.Println("Closing proxy connection...")
		serverConn.Close()
	}()

	vlog.Printf("Successfully dialed up to %v...\n", portB)

	serverConn.Write(data[0:n])
	// Copy data between client and server
	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}
