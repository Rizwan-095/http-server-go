package main

import (
	"fmt"
	"regexp"

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

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	req := make([]byte, 1024)
	conn.Read(req)
	re := regexp.MustCompile(`/echo/([^ ]+)`)

	// Find the substring that matches the pattern
	match := re.FindStringSubmatch(string(req))
	requestStr := match[1]
	// if !strings.HasPrefix(string(req), "GET /echo/ HTTP/1.1") {
	// 	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	// }
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nContent-Type: text/plain\r\nContent-Length: " + string(len(requestStr)) + "\r\n\r\n" + requestStr))
}
