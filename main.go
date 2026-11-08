package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

const M = 64

var bitArr [M]byte
var out strings.Builder

func h1(key []byte) int {
	// TODO: sum of the key's bytes, mod 64
	_ = key
	return 0
}

func h2(key []byte) int {
	// TODO: sum of byte * (index + 1), mod 64
	_ = key
	return 0
}

func forEachLine(handle func(cmd string, rest []byte)) {
	data, _ := ioutil.ReadAll(os.Stdin)
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		sp := bytes.IndexByte(line, ' ')
		if sp < 0 {
			handle(string(line), nil)
		} else {
			handle(string(line[:sp]), line[sp+1:])
		}
	}
}

func handle(cmd string, key []byte) {
	switch cmd {
	case "ADD":
		// TODO: set bits h1(key) and h2(key), then print OK
	case "CHECK":
		// TODO: MAYBE only if both bits are 1, otherwise NO
	case "BITS":
		// TODO: print the 64 bits as a string of 0 and 1, bit 0 first
	}
}

func main() {
	forEachLine(handle)
	os.Stdout.WriteString(out.String())
}
