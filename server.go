package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/NikhilSharmaWe/go-cache/cache"
	"github.com/NikhilSharmaWe/go-cache/client"
	"github.com/NikhilSharmaWe/go-cache/proto"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	members map[*client.Client]struct{} // deleting from a map is much easier than from a slice
	cache   cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		members:    make(map[*client.Client]struct{}),
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Fatal(err)
			}
		}()
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

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.LeaderAddr)
	}

	log.Println("connected to leader:", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, proto.CmdJoin)

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	// fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := proto.ParseCommand(conn) // it is a blocking call and since net.Conn implements the Reader and Writer interfaces
		if err != nil {
			if err != io.EOF {
				log.Println("parse command error:", err)
			}
			fmt.Println(err)
			break
		}

		go s.handleCommand(conn, cmd)
	}

	// fmt.Println("connection closed:", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) { // v is a variable declared inside the switch statement, and its type is determined by the type assertion in each case.
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
		s.handleGetCommand(conn, v)
	case *proto.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *proto.CommandJoin) error {
	fmt.Println("member just joined the cluster:", conn.RemoteAddr())
	client := client.NewClientFromConn(conn)
	s.members[client] = struct{}{}
	return nil
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	fmt.Printf("SET %s to %s\n", cmd.Key, cmd.Value)

	go func() {
		for member := range s.members {
			err := member.Set(context.Background(), cmd.Key, cmd.Value, cmd.TTL)
			if err != nil {
				log.Println("forward to member error:", err)
			}
		}
	}()

	resp := proto.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		if err != nil {
			return err
		}
		return err
	}

	resp.Status = proto.StatussOK
	_, err := conn.Write(resp.Bytes())
	return err
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {
	resp := proto.ResponseGet{}
	val, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusKeyNotFound
		_, err := conn.Write(resp.Bytes())
		if err != nil {
			return err
		}
	}

	resp.Status = proto.StatussOK
	resp.Value = val
	_, err = conn.Write(resp.Bytes())

	return err
}
