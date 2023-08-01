package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/NikhilSharmaWe/go-cache/cache"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()
	opts := ServerOpts{
		ListenAddr: *listenAddr,
		LeaderAddr: *leaderAddr,
		IsLeader:   true,
	}

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", *leaderAddr) // this is client side conn
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("SET Foo Bar 25000000000"))

		time.Sleep(time.Second * 2)
		conn.Write([]byte("GET Foo"))

		buf := make([]byte, 1000)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buf[:n]))
	}()

	server := NewServer(opts, *cache.New())
	server.Start()
}
