package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const filepath = "D://text.html"

var isFinish = false

func main() {

	go writeFile()
	startHttp(80)
}

func startHttp(port int) {
	//start server
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	//mux.HandleFunc("/", replay)

	mux.HandleFunc("/download1", download1)

	fmt.Printf("监听%d端口\n", port)
	err = http.Serve(l, mux)
	if err != nil {
		panic(err)
	}
}

func writeFile() {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		file.Write([]byte(fmt.Sprintf("%d\n", i)))
		time.Sleep(time.Millisecond * 10)
	}

	isFinish = true
	file.Close()
}
func download1(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != io.EOF {
			w.Write(buf[:n])
		} else {
			if isFinish {
				break
			} else {
				//fmt.Println("EOF， 但未完成， 休息下")
				time.Sleep(time.Millisecond * 10)
				w.Write([]byte("EOF， 但未完成， 休息下\n"))
			}

		}
	}
	fmt.Printf("全部发送完毕")

}
