package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// handleConnection reads lines from client and echoes them in uppercase
func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("[TCP Server] Client connected: %s\n", remoteAddr)

	reader := bufio.NewReader(conn)
	for {
		// Read line terminated by '\n'
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("[TCP Server] Read error from %s: %v\n", remoteAddr, err)
			}
			break
		}

		cleanMsg := strings.TrimSpace(message)
		if cleanMsg == "QUIT" {
			conn.Write([]byte("Goodbye!\n"))
			break
		}

		response := fmt.Sprintf("ECHO: %s\n", strings.ToUpper(cleanMsg))
		conn.Write([]byte(response))
	}

	fmt.Printf("[TCP Server] Client disconnected: %s\n", remoteAddr)
}

func main() {
	// 1. Listen on local TCP port
	// Specifying port ":0" asks the operating system to allocate an ephemeral free port.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("Failed to bind TCP listener: %v\n", err)
		return
	}
	// Idiom: defer Close() immediately after checking error so socket is freed on return
	defer listener.Close()

	// ── Type Assertion on Interfaces ──
	// listener.Addr() returns the generic interface `net.Addr` (which only exposes Network() and String()).
	// But since we created a TCP listener, the underlying concrete type is `*net.TCPAddr`.
	// The syntax `val.(*net.TCPAddr)` is a TYPE ASSERTION:
	// It extracts the concrete `*net.TCPAddr` pointer, allowing us to read its `.Port` integer field!
	tcpAddr := listener.Addr().(*net.TCPAddr)
	port := tcpAddr.Port
	fmt.Printf("=== TCP Echo Server listening on 127.0.0.1:%d ===\n", port)


	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn) // Handle concurrently!
		}
	}()

	// 2. Test Client Connection via net.Dial
	time.Sleep(50 * time.Millisecond)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Printf("Client dial error: %v\n", err)
		return
	}
	defer conn.Close()

	// Send messages to server
	messages := []string{"hello go", "tcp networking fundamentals", "QUIT"}
	clientReader := bufio.NewReader(conn)

	for _, msg := range messages {
		fmt.Printf("[Client -> Server] %s\n", msg)
		fmt.Fprintf(conn, "%s\n", msg)

		reply, _ := clientReader.ReadString('\n')
		fmt.Printf("[Server -> Client] %s", reply)
	}
}
