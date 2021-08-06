package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// dial
	conn, err := net.Dial("tcp", ":2020")
	if err != nil {
		log.Fatalln("Dial err: %v", err)
	}
	defer conn.Close()

	// io分离
	done := make(chan struct{}) // 控制goroutine顺序
	// 读
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{}
	}()

	// 写
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}

	<-done
}
