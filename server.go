package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/NikhilSharmaWe/go-cache/cache"
	"github.com/NikhilSharmaWe/go-cache/proto"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
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
	defer conn.Close()

	fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := proto.ParseCommand(conn) // it is a blocking call and since net.Conn implements the Reader and Writer interfaces
		if err != nil {
			if err != io.EOF {
				log.Println("parse command error:", err)
			}
			break
		}

		go s.handleCommand(conn, cmd)
	}

	fmt.Println("connection closed:", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) { // v is a variable declared inside the switch statement, and its type is determined by the type assertion in each case.
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
		// case *CommandGet:
		// 	s.handleGetCommand(conn, v)
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)
	return s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL))
}

// func (s *Server) handleGetCommand(conn net.Conn, cmd *CommandGet) error {
// 	log.Printf("SET %s to %s", cmd.Key)
// 	return s.cache.Get(cmd.Key)
// }
