package client

import (
	"context"
	"fmt"
	"net"

	"github.com/NikhilSharmaWe/go-cache/proto"
)

type Options struct {
}

type Client struct {
	conn net.Conn
	opts Options
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		opts: opts,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) (any, error) {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to write to the server")
	}

	return nil, nil
}
