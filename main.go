package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/NikhilSharmaWe/go-cache/cache"
	"github.com/NikhilSharmaWe/go-cache/client"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 10)
		if opts.IsLeader {
			SendStuff()
		}
	}()

	server := NewServer(opts, cache.New())
	server.Start()
}

func SendStuff() {
	for i := 0; i < 100; i++ {
		go func(i int) {
			c, err := client.New(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}
			var (
				key   = fmt.Sprintf("key_%d", i)
				value = fmt.Sprintf("val_%d", i)
			)

			err = c.Set(context.Background(), []byte(key), []byte(value), 0)
			if err != nil {
				log.Fatal(err)
			}

			fetchedValue, err := c.Get(context.Background(), []byte(key))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(fetchedValue))

			time.Sleep(time.Second)
		}(i)
	}
}

// func randomBytes(n int) []byte {
// 	buf := make([]byte, n)
// 	io.ReadFull(rand.Reader, buf)
// 	return buf
// }
