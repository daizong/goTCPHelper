package main

import (
	"TCPHelper"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

var lock sync.Mutex

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listener, err := net.Listen("tcp4", "127.0.0.1:8787")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}

			helper := TCPHelper.NewHelper(conn, TCPHelper.DefaultPacketProtocol())
			go helper.ReadLoop(func(msg []byte) {
				lock.Lock()
				defer lock.Unlock()
				fmt.Printf("server received: %s\n", msg)
				helper.Write([]byte("server says: hi client!"))
			})
		}
	}()

	go SimulateClient()
	go SimulateClient()
	go SimulateClient()
	go SimulateClient()
	go SimulateClient()

	exit := make(chan struct{})
	<-exit
}

func SimulateClient() {
	conn, err := net.Dial("tcp4", "127.0.0.1:8787")
	if err != nil {
		panic(err)
	}
	helper := TCPHelper.NewHelper(conn, TCPHelper.DefaultPacketProtocol())
	go helper.ReadLoop(func(msg []byte) {
		lock.Lock()
		defer lock.Unlock()
		fmt.Printf("client received: %s\n", msg)
	})

	_, err = helper.Write([]byte("client says: hi server!"))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			i := 10
			for i > 0 {
				helper.Write([]byte("client says: hi server!"))
				i -= 1
			}

		}
	}
}
