package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	req := make([]byte, 1024)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
	request := string(req)
	path := strings.Split(string(req), " ")[1]
	switch {
	case path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	case strings.Split(path, "/")[1] == "echo":
		message := strings.Split(path, "/")[2]
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))
		return
	case strings.Split(path, "/")[1] == "user-agent":
		headers := strings.Split(request, "\r\n")[1:]
		userAgent := ""
		for _, header := range headers {
			if strings.HasPrefix(header, "User-Agent") {
				userAgent = strings.Split(header, ": ")[1]
			}
		}
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)))
		return
	case strings.Split(path, "/")[1] == "files":
		fileName := strings.Split(path, "/")[2]
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), string(fileContent))))
		return
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}
}
