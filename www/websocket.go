package main

// websocket のデータを bash スクリプトでやり取りする
// http://ssklogs.blogspot.jp/2012/10/websockets-handshake-using-netcatbash.html

// 参考サイト
// http://gary.burd.info/go-websocket-chat

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"text/template"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"log"
)

const (
	PORT = 8000
	HOST = "localhost"
)

// Echo Server
func Echo(ws *websocket.Conn) {
	for {
		in := make([]byte, 1024)
		n, err := ws.Read(in)
		if  err != nil {
			fmt.Println(err.Error())
			break
		}
		received := strings.Trim(string(in[:n]), "\n")
		fmt.Printf("Received: %s\n", received)
		ws.Write([]byte(received + "\r\n"))
	}
}

func Render(filename string, name string, data map[string]string) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile(filename)
		if err != nil {log.Fatal(err)}
		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {log.Fatal(err)}
		tmpl.Execute(w, data)
	}
	return http.HandlerFunc(handler)
}

func StaticFile(name string, filename string) (url string, handler http.Handler) {
	url = name
	h := func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile(filename)
		if err != nil {log.Fatal(err)}
		w.Write(content)
	}
	handler = http.HandlerFunc(h)
	return
}

func Start(port int, host string) {
	p := strconv.Itoa(port)
	http.Handle("/echo", websocket.Handler(Echo));
	http.Handle(StaticFile("/client", "client.html"))
	http.Handle("/conn.js", Render("./static/conn.js", "conn",
		map[string]string{`port`:p, `host`:host}))
	http.Handle("/", http.FileServer(http.Dir("./static/")));
	fmt.Printf("listen: *:%s\n", p)
	err := http.ListenAndServe(":" + p, nil);
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	Start(PORT, HOST)
}

