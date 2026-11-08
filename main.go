package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
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

func optimal(p, n float64) (m, k float64) {
	m = math.Ceil(-n * math.Log(p) / (math.Ln2 * math.Ln2))
	k = math.Floor((m/n)*math.Ln2 + 0.5)
	return
}

func fp(m, n, k float64) float64 {
	return math.Pow(1-math.Exp(-k*n/m), k)
}

func bpi(p float64) float64 {
	return -math.Log(p) / (math.Ln2 * math.Ln2)
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
	case "OPTIMAL":
		fields := strings.Fields(string(key))
		p, _ := strconv.ParseFloat(fields[0], 64)
		n, _ := strconv.ParseFloat(fields[1], 64)
		m, k := optimal(p, n)
		fmt.Printf("m=%d k=%d\n", int(m), int(k))
	case "FP":
		fields := strings.Fields(string(key))
		m, _ := strconv.ParseFloat(fields[0], 64)
		n, _ := strconv.ParseFloat(fields[1], 64)
		k, _ := strconv.ParseFloat(fields[2], 64)
		fmt.Printf("%.6f\n", fp(m, n, k))
	case "BPI":
		p, _ := strconv.ParseFloat(strings.TrimSpace(string(key)), 64)
		fmt.Printf("%.4f\n", bpi(p))
	}
}

func forEachLine(handle func(cmd string, rest []byte)) {
	data, _ := io.ReadAll(os.Stdin)
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

func main() {
	forEachLine(handle)
}
