package wasm_go

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsigned(t *testing.T) {
	cases := map[uint64]string{
		0:                    "00000000",
		0x7f:                 "01111111",
		0x80:                 "00000001 10000000",
		0xff:                 "00000001 11111111",
		0x91d:                "00010010 10011101",
		0xef17:               "00000011 11011110 10010111",
		624485:               "00100110 10001110 11100101",
		0xffff:               "00000011 11111111 11111111",
		18446744073709551615: "00000001 11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
	}

	for expect, binaryString := range cases {
		r := leb128Reader{bytes: binaryStringToBytes(binaryString), pos: 0}
		v, err := r.eatU64()
		assert.NoError(t, err)
		assert.Equal(t, expect, v)
	}
}

func TestSigned(t *testing.T) {
	cases := map[int64]string{
		-9223372036854775808: "01111111 10000000 10000000 10000000 10000000 10000000 10000000 10000000 10000000 10000000",
		-624485:              "01011001 11110001 10011011",
		^0x40:                "01111111 10111111",
		^0x3f:                "01000000",
		-1:                   "01111111",
		0:                    "00000000",
		1:                    "00000001",
		0x3f:                 "00111111",
		0x40:                 "00000000 11000000",
		0xef17:               "00000011 11011110 10010111",
		9223372036854775807:  "00000000 11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
	}

	for expect, binaryString := range cases {
		r := leb128Reader{bytes: binaryStringToBytes(binaryString), pos: 0}
		v, err := r.eatI64()
		assert.NoError(t, err)
		assert.Equal(t, expect, v)
	}
}

func binaryStringToBytes(s string) []byte {
	parts := strings.Split(s, " ")
	l := len(parts)
	b := make([]byte, l)

	for i, v := range parts {
		fmt.Sscanf(v, "%b", &b[l-i-1])
	}
	return b
}
