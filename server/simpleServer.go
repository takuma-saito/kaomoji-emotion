
package server
import (
	"fmt"
	"net"
	"bufio"
	"io"
	"os"
	"strings"
	"strconv"
)

type Filter func(string) string

func Exit(content string, name string) {
	println(content, name)
	os.Exit(1)
}

func echoServer(conn net.Conn, filter Filter) {
	bufr := bufio.NewReader(conn)
	io.WriteString(conn, "If you want to exit, type 'exit'.\r\n > ")
	for {
		line, _ := bufr.ReadString('\n')
		line = strings.Trim(line, "\r\n")
		if line[0:len(line)] == "exit" {break}
		fmt.Printf("in: %s\n", line)
		fmt.Printf("out: %s\n", filter(line))
		io.WriteString(conn, fmt.Sprintf("%s\r\n > ", filter(line)))
	}
	conn.Close()
}

func Start(port int, filter Filter) {
	print(fmt.Sprintf("listen: 0.0.0.0:%s\n", strconv.Itoa(port)))
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(port)))
	if err != nil {Exit("Error Listen: ", err.Error())}
	for {
		conn, err := listener.Accept()
		if err != nil {Exit("Error Accept: ", err.Error())}
		go echoServer(conn, filter)
	}
}


