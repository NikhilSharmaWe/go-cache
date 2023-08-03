package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	r := bytes.NewReader(cmd.Bytes())
	pcmd := ParseCommand(r)

	assert.Equal(t, cmd, pcmd)
}

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	r := bytes.NewReader(cmd.Bytes())
	pcmd := ParseCommand(r)

	assert.Equal(t, cmd, pcmd)
}

func BenchmarkParseCommand(b *testing.B) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   6,
	}

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(cmd.Bytes())
		ParseCommand(r)
	}
}
