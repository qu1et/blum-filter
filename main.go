package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const M = 64

var bitArr [M]byte
var out strings.Builder

func h1(key []byte) int {
	sum := 0
	for _, b := range key {
		sum += int(b)
	}
	return sum % 64
}

func h2(key []byte) int {
	sum := 0
	for i, b := range key {
		sum += (i + 1) * int(b)
	}
	return sum % 64
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
		bitArr[h1(key)] = 1
		bitArr[h2(key)] = 1
		fmt.Println("OK")
	case "CHECK":
		if bitArr[h1(key)] == 1 && bitArr[h2(key)] == 1 {
			fmt.Println("MAYBE")
		} else {
			fmt.Println("NO")
		}
	case "BITS":
		for _, v := range bitArr {
			fmt.Print(v)
		}
		fmt.Println()
	}
}

func main() {
	forEachLine(handle)
	os.Stdout.WriteString(out.String())
}
