package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// REASON WE ARE USING THIS APPROACH IS BECAUSE FOR ENCODING AND DECODING
// STRUCT INTO []BYTE THIS IS THE BEST APPROACH

// and we need to do that because tcp conn works with data in []byte

//ALSO THIS TEACH USE HOW PROTOBUF AND MESSAGEPACKER ARE WORKING IN THE BACKEND

type Command byte

const (
	CmdNone Command = iota
	CmdSet
	CmdGet
	CmdDel
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)

	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	valLen := int32(len(c.Value))
	binary.Write(buf, binary.LittleEndian, valLen)
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, int32(c.TTL))

	return buf.Bytes()
}

type CommandGet struct {
	Key []byte
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)

	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

// the binary.Read reads the contents of the reader sequentially which
// has not been read, when we first read the reader in this case through the
// reader made from the []byte returned from the c.Bytes() function, it first
// returns the CmdSet value written, then keyLen and so on.
func ParseCommand(r io.Reader) any {
	var cmd Command
	binary.Read(r, binary.LittleEndian, &cmd)

	switch cmd {
	case CmdSet:
		return parseSetCommand(r)
	case CmdGet:
		return parseGetCommand(r)
	default:
		panic("invalid command")
	}
}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valLen int32
	binary.Read(r, binary.LittleEndian, &valLen)
	cmd.Value = make([]byte, valLen)
	binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)

	return cmd
}

func parseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	fmt.Println(keyLen)

	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	return cmd
}
