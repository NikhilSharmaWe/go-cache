package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// REASON WE ARE USING THIS APPROACH IS BECAUSE FOR ENCODING AND DECODING
// STRUCT INTO []BYTE THIS IS THE BEST APPROACH
// BECAUSE WHEN WE CONVERT THE STRUCT INTO []BYTE THEN HOW CAN WE
// ACCESS INDIVIDUAL FIELDS, THROUGH BINARY.READ WE GET THE VALUES
// ADDED TO THE BUFFER IN SEQUENCE IT WAD ADDED.

// ** and we need to do that because tcp conn works with data in []byte, the read
// and write functions work in []byte typed data **

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
func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	err := binary.Read(r, binary.LittleEndian, &cmd)
	if err != nil {
		return nil, err
	}

	switch cmd {
	case CmdSet:
		return parseSetCommand(r), nil
	case CmdGet:
		return parseGetCommand(r), nil
	default:
		return nil, fmt.Errorf("invalid command")
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
