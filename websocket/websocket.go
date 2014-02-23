package websocket

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
	"os"
	"log"
	"time"
)

const ROOT = "./websocket/"

type Filter func(string) string

type Message struct {
	Face string `json:"face"`
	Emotion string `json:"emotion"`
}

// Echo Server
func Echo(ws *websocket.Conn) {
	defer ws.Close()
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

type handleConn func(*websocket.Conn)

// *TODO* Log フォーマットを分離させる
func HandleJson(filter Filter) handleConn {
	return func(ws *websocket.Conn) {
		log, err := os.OpenFile("./log/faces.txt", os.O_RDWR|os.O_APPEND, 0600);
		access, err := os.OpenFile("./log/access.txt", os.O_RDWR|os.O_APPEND, 0600);
		req := ws.Request()
		access.WriteString(
			fmt.Sprintf("ip:%s\t" + "host:%s\t" + "ua:%s\t" + "time:%s\t\n",
				req.RemoteAddr, req.Host, req.Header["User-Agent"][0], time.Now()))
		access.Close()
		if err != nil {panic(err)}
		defer func() {
			ws.Close()
			log.Close()
		}()
		for {
			var request, response Message
			err := websocket.JSON.Receive(ws, &request)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				break
			}
			log.WriteString(request.Face + "\n")
			response.Emotion = filter(request.Face)
			response.Face = request.Face
			fmt.Printf("receive: %v\n", request)
			fmt.Printf("send: %v\n", response)
			err = websocket.JSON.Send(ws, response)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				break
			}
		}
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

func Start(port int, host string, filter Filter) {
	p := strconv.Itoa(port)
	http.Handle("/echo", websocket.Handler(HandleJson(filter)));
	http.Handle(StaticFile("/", ROOT + "client.html"))
	http.Handle("/conn.js", Render(ROOT + "conn.js", "conn",
		map[string]string{`port`:p, `host`:host}))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(ROOT + "static/"))));
	fmt.Printf("listen: *:%s\n", p)
	err := http.ListenAndServe(":" + p, nil);
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

