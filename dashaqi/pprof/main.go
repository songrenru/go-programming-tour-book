package main

import (
	"github.com/songrenru/dashaqi"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)


func main() {
	go func() {
		for {
			log.Printf("len: %d", dashaqi.Add("go-programming-exercise"))
			time.Sleep(time.Millisecond*10)
		}
	}()

	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}
