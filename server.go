package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/NikhilSharmaWe/go-cache/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cache
}

func NewServer(opts ServerOpts, c cache.Cache) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s", err)
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	var (
		rawStr = string(rawCmd)
		parts  = strings.Split(rawStr, " ")
	)

	if len(parts) == 0 {
		// respond
		log.Println("invalid command")
		return
	}
	// cmd := parts[0]
}
