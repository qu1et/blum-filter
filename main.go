package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

const (
	M  = 64
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
)

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

func murmur3_32(key []byte, seed uint32) uint32 {
	h := seed
	nblocks := len(key) / 4
	for i := range nblocks {
		k := binary.LittleEndian.Uint32(key[i*4 : i*4+4])
		k *= c1
		k = bits.RotateLeft32(k, 15)
		k *= c2
		h ^= k
		h = bits.RotateLeft32(h, 13)
		h = h*5 + 0xe6546b64
	}

	tail := key[nblocks*4:]
	if len(tail) > 0 {
		var k uint32
		switch len(tail) {
		case 3:
			k ^= uint32(tail[2]) << 16
			fallthrough
		case 2:
			k ^= uint32(tail[1]) << 8
			fallthrough
		case 1:
			k ^= uint32(tail[0])
		}

		k *= c1
		k = bits.RotateLeft32(k, 15)
		k *= c2
		h ^= k
	}

	h ^= uint32(len(key))

	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

func hashes(key []byte) (uint32, uint32) {
	ha := murmur3_32(key, 0)
	hb := murmur3_32(key, 0x9747b28c)
	return ha, hb
}

func postions(m, k, key []byte) []string {
	var pos []string
	k_i, _ := strconv.Atoi(string(k))
	m_i, _ := strconv.Atoi(string(m))
	ha, hb := hashes(key)

	for i := range k_i {
		p := (int(ha) + i*int(hb)) % m_i
		pos = append(pos, strconv.FormatInt(int64(p), 10))
	}

	return pos
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
	case "MURMUR":
		var k []byte
		var seed []byte
		key = bytes.TrimRight(key, "\r\n")
		if ind := bytes.IndexByte(key, ' '); ind != -1 {
			seed = key[:ind]
			k = key[ind+1:]
		} else {
			seed = key
		}
		n, _ := strconv.ParseUint(string(seed), 16, 32)
		fmt.Printf("%08x\n", murmur3_32(k, uint32(n)))
	case "HASHES":
		ha, hb := hashes(key)
		fmt.Printf("ha=%08x hb=%08x\n", ha, hb)
	case "POSITIONS":
		fileds := bytes.Fields(key)
		res := postions(fileds[0], fileds[1], fileds[2])
		fmt.Println(strings.Join(res, ","))
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
