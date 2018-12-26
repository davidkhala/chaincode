package main

import (
	. "github.com/davidkhala/goutils"
	. "github.com/davidkhala/goutils/crypto"
)

func Hash(data []byte) string {
	return HexEncode(HashSha512(data))
}
