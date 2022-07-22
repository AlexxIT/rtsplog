package main

import (
	"fmt"
	"os"
	"rtsplog/rtsp"
)

func main() {
	conn := new(rtsp.Conn)
	conn.Out = func(msg interface{}) {
		fmt.Printf("%s\n\n", msg)
	}
	if err := conn.Dial(os.Args[1]); err != nil {
		panic(err)
	}
}
