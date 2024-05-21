package main

import (
	"context"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	go func() {
		time.Sleep(time.Second * 5)
		//ctx.Done()
		//cancel()
	}()
	//var ch = make(chan struct{}, 1)
	<-ctx.Done()
	log.Println(ctx.Err())
	println("end")
	cancel()
}
