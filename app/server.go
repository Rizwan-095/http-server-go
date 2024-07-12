package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	//Server is now listening on tcp port 4221.
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	//Server accepting concurrent http requests.
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

// Function to handle requests.
func handleConnection(conn net.Conn) {
	req := make([]byte, 1024)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}
	request := string(req)
	path := strings.Split(string(req), " ")[1]
	headers := strings.Split(request, "\r\n")[1:]
	absPath := strings.Split(path, "/")[1]
	switch {
	//Server responding to http request with status code (200 Ok).
	case path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	//Server can extract URL path of request.
	case absPath == "echo":
		message := strings.Split(path, "/")[2]
		encodingType := ""
		for _, header := range headers {
			if strings.HasPrefix(header, "Accept-Encoding") {
				encodingType = strings.Split(header, ": ")[1]
			}
		}
		if len(encodingType) < 1 {
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))
			return
		} else {
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: %s\r\nContent-Length: %d\r\n\r\n%s", len(message), encodingType, message)))
		}
	//Server Reading headers.
	case absPath == "user-agent":
		userAgent := ""
		for _, header := range headers {
			if strings.HasPrefix(header, "User-Agent") {
				userAgent = strings.Split(header, ": ")[1]
			}
		}
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)))
		return
	//Server returing file content for Http GET and creating new file with content from request body for Http POST request respectively.
	case absPath == "files":
		dir := os.Args[2]
		fileName := strings.Split(path, "/")[2]
		if strings.Split(string(req), " ")[0] == "GET" {
			fileContent, err := os.ReadFile(dir + fileName)
			if err != nil {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), fileContent)))
			return
		} else {
			fileContent := []byte(strings.Trim(strings.SplitN(request, "\r\n\r\n", 2)[1], "\x00"))
			err := os.WriteFile(dir+fileName, fileContent, os.ModePerm)
			if err != nil {
				fmt.Println("Error reading request: ", err.Error())
				return
			}
			conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
		}
	//Server responding with status code (404 Not Found) if request is wrong
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}
}
