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

type Bloom struct {
	m     int
	k     int
	bits  []byte
	adds  int
	nInit int
}

var bloom Bloom

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

func (b *Bloom) postions(key []byte) []int {
	ha, hb := hashes(key)
	pos := make([]int, b.k)
	for i := 0; i < b.k; i++ {
		p := (int(ha) + i*int(hb)) % b.m
		pos[i] = int(p)
	}
	return pos
}

func (b *Bloom) bloomInit(n, p []byte) (int, int) {
	n1, _ := strconv.ParseFloat(string(n), 64)
	p1, _ := strconv.ParseFloat(string(p), 64)
	bloom.nInit = int(n1)
	m := math.Ceil(-n1 * math.Log(p1) / (math.Ln2 * math.Ln2))
	k := math.Floor(m/n1*math.Ln2 + 0.5)
	k = max(k, 1)
	b.m = int(m)
	b.k = int(k)
	b.adds = 0
	b.bits = make([]byte, (int(b.m)+7)/8)
	return b.m, b.k
}

func handle(cmd string, key []byte) {
	switch cmd {
	case "ADD":
		for _, p := range bloom.postions(key) {
			bloom.bits[p/8] |= 1 << (p % 8)
		}
		bloom.adds++
		fmt.Println("OK")
	case "CHECK":
		isTaken := true
		for _, p := range bloom.postions(key) {
			if bloom.bits[p/8]&(1<<(p%8)) == 0 {
				isTaken = false
				break
			}
		}
		if isTaken {
			fmt.Println("MAYBE")
		} else {
			fmt.Println("NO")
		}
	case "BITS":
		fmt.Printf("%x\n", bloom.bits)
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
		res := bloom.postions(key)
		str := fmt.Sprint(res)
		str = strings.Trim(str, "[]")
		str = strings.ReplaceAll(str, " ", ",")
		fmt.Println(str)
	case "INIT":
		fields := bytes.Fields(key)
		m, k := bloom.bloomInit(fields[0], fields[1])
		fmt.Printf("OK m=%d k=%d\n", m, k)
	case "STATS":
		set := 0
		load := float64(bloom.adds) / float64(bloom.nInit)
		for _, v := range bloom.bits {
			set += bits.OnesCount8(v)
		}
		fmt.Printf("m=%d k=%d n=%d load=%.4f set=%d\n", bloom.m, bloom.k, bloom.adds, load, set)
	case "DUMP":
		fmt.Printf("%x\n", bloom.bits)
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
