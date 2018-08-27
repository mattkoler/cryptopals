package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func hexTo64(h string) (string, error) {
	data, err := hex.DecodeString(h)
	if err != nil {
		return "", fmt.Errorf("error decoding string to hex: %s", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func fixedXOR(h, x string) (string, error) {
	data, err := hex.DecodeString(h)
	if err != nil {
		return "", fmt.Errorf("error decoding string to hex: %s", err)
	}
	xor, err := hex.DecodeString(x)
	if err != nil {
		return "", fmt.Errorf("error decoding string to hex: %s", err)
	}
	n := len(data)
	if len(xor) < n {
		n = len(xor)
	}
	var val []byte
	for i := 0; i < n; i++ {
		val = append(val, data[i]^xor[i])
	}

	dst := hex.EncodeToString(val)

	return dst, nil
}

func singleByteXOR(h string) (string, float64, int, error) {
	data, err := hex.DecodeString(h)
	if err != nil {
		return "", 0, 0, fmt.Errorf("error decoding string to hex: %s", err)
	}
	score := 10000.00
	var best string
	var decoder int
	for i := 0; i < 256; i++ {
		var decode []byte
		for j := 0; j < len(data); j++ {
			decode = append(decode, data[j]^byte(i))
		}
		decodedStr := string(decode)
		if s := scoreDecode(decodedStr); s < score {
			score = s
			best = decodedStr
			decoder = i
		}
	}
	return best, score, decoder, nil
}

func scoreDecode(s string) float64 {
	letterDist := map[byte]float64{
		'a': 8.167,
		'b': 1.492,
		'c': 2.782,
		'd': 4.253,
		'e': 12.702,
		'f': 2.228,
		'g': 2.015,
		'h': 6.094,
		'i': 6.966,
		'j': 0.153,
		'k': 0.772,
		'l': 4.025,
		'm': 2.406,
		'n': 6.749,
		'o': 7.507,
		'p': 1.929,
		'q': 0.095,
		'r': 5.987,
		's': 6.327,
		't': 9.056,
		'u': 2.758,
		'v': 0.978,
		'w': 2.360,
		'x': 0.150,
		'y': 1.974,
		'z': 0.074,
	}
	dist := make(map[byte]float64)
	s = strings.ToLower(s)
	length := len(s)
	for i := 0; i < length; i++ {
		char := s[i]
		if _, ok := dist[char]; ok {
			dist[char]++
		} else {
			dist[char] = 1
		}
	}
	var score float64
	for char, got := range dist {
		if perc, ok := letterDist[char]; ok {
			off := got / float64(length)
			score += math.Abs(perc - off)
		} else {
			if char == ' ' {
				continue
			}
			score += got * 100
		}
	}
	return score // lower is better
}

func findXORstring(p string) (string, int, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", 0, fmt.Errorf("error opening file path %s: %q", p, err)
	}
	defer f.Close()

	score := 10000.00
	var bestStr string
	var decoded int

	r := bufio.NewScanner(f)

	for r.Scan() {
		if text, sco, uni, err := singleByteXOR(r.Text()); err == nil {
			if sco < score {
				score = sco
				bestStr = text
				decoded = uni
			}
		} else {
			return "", 0, fmt.Errorf("error converting: %q", err)
		}
	}
	if err := r.Err(); err != nil {
		return "", 0, fmt.Errorf("error reading file: %q", err)
	}
	return bestStr, decoded, nil
}

func repeatingKeyXOR(input, key string) (string, error) {
	var val []byte
	for i := 0; i < len(input); i++ {
		val = append(val, input[i]^key[i%len(key)])
	}
	return hex.EncodeToString(val), nil
}

func hammingDist(s1, s2 string) (int, error) {
	if len(s1) != len(s2) {
		return 0, fmt.Errorf("must have 2 strings of equal length: %s %s", s1, s2)
	}
	var dist int
	for i := 0; i < len(s1); i++ {
		diff := s1[i] ^ s2[i]
		dist += bits.OnesCount(uint(diff))
	}
	return dist, nil
}

type key struct {
	size  int
	score float64
}

func findKeysize(cipherFile string) ([]int, error) {
	f, err := os.Open(cipherFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %q", err)
	}
	s := bufio.NewScanner(f)
	var cipher64 string

	for s.Scan() {
		cipher64 += s.Text()
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(cipher64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base 64: %q", err)
	}

	var keys []key
	for i := 2; i < 41; i++ {
		var blocks [][]byte
		for j := 0; j < 5; j++ {
			blocks = append(blocks, cipherBytes[j*i:j*i+i]) // generate 5 blocks of i size
		}
		var keyScore float64
		for b := range blocks {
			if b == 0 {
				continue
			}
			score, err := hammingDist(string(blocks[b-1]), string(blocks[b]))
			if err != nil {
				return nil, fmt.Errorf("unable to compute hammign dist: %q", err)
			}
			keyScore += float64(score)
		}
		normalized := keyScore / (float64(i) * 4.0)
		keys = append(keys, key{size: i, score: normalized})
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i].score < keys[j].score })

	var sizes []int
	for _, keySize := range keys {
		sizes = append(sizes, keySize.size)
	}
	return sizes, nil
}

func transposeBlocks(cipherFile string, size int) ([][]byte, error) {
	f, err := os.Open(cipherFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %q", err)
	}
	s := bufio.NewScanner(f)
	var cipher64 string

	for s.Scan() {
		cipher64 += s.Text()
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(cipher64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base 64: %q", err)
	}
	blocks := make([][]byte, size)
	for i := 0; i < len(cipherBytes); i++ {
		block := i % size
		blocks[block] = append(blocks[block], cipherBytes[i])
	}
	return blocks, nil
}

func solveRepeatingKey(cipherFile string) (string, error) {
	type decodeAttempt struct {
		keySize     int
		transBlocks [][]byte
		cipher      string
		score       float64
	}
	sizes, err := findKeysize(cipherFile)
	if err != nil {
		return "", fmt.Errorf("unable to find key size: %q", err)
	}
	var attempts []decodeAttempt
	for i := 0; i < 20; i++ {
		blocks, err := transposeBlocks(cipherFile, sizes[i])
		if err != nil {
			return "", fmt.Errorf("unable to transpose blocks: %q", err)
		}
		attempts = append(attempts, decodeAttempt{keySize: sizes[i], transBlocks: blocks})
	}
	for _, attempt := range attempts {
		for _, block := range attempt.transBlocks {
			_, score, uni, err := singleByteXOR(hex.EncodeToString(block))
			if err != nil {
				return "", fmt.Errorf("unable to find xor: %q", err)
			}
			attempt.cipher += string(uni)
			attempt.score += score
		}
		fmt.Printf("Attempt at size %d gave %s\n", attempt.keySize, attempt.cipher)
	}
	return "", nil
}
