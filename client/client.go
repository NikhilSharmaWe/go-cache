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

func NewClientFromConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
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

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet{
		Key: key,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := proto.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}

	if resp.Status == proto.StatusKeyNotFound {
		return nil, fmt.Errorf("could not find key [%s]", key)
	}

	if resp.Status != proto.StatussOK {
		return nil, fmt.Errorf("server responded with a non OK status [%s]", resp.Status)
	}

	return resp.Value, nil
}

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) error {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write to the server")
	}

	resp, err := proto.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != proto.StatussOK {
		return fmt.Errorf("server responded with a non OK status [%s]", resp.Status)
	}

	return nil
}
