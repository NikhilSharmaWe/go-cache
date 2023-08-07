package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NikhilSharmaWe/go-cache/client"
)

func main() {
	SendStuff()
}

func SendStuff() {
	c, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		var (
			key   = fmt.Sprintf("key_%d", i)
			value = fmt.Sprintf("val_%d", i)
		)

		if err := c.Set(context.Background(), []byte(key), []byte(value), 0); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}
