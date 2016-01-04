package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

func readSecretKey(filename string) string {
	// Open the file.
	f, _ := os.Open(filename)
	defer f.Close()

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		return line
	}
	return ""
}

func computeHashOfSecretKey(key string, number int) string {

	return ""
}

func isHashStartingWithNZeros(N int, hash []byte) bool {
	for i := 0; i < N; i++ {
		c := hash[i/2]
		if i&1 == 0 {
			if c&0xF0 != 0 {
				return false
			}
		} else {
			if c&0x0F != 0 {
				return false
			}
		}
	}
	return true
}

func main() {
	var secretKey = readSecretKey("input.text")

	hasher := md5.New()
	var saltedSecretKey = fmt.Sprintf("%s%v", secretKey, 609043)
	hasher.Write([]byte(saltedSecretKey))
	hash := hasher.Sum(nil)
	if isHashStartingWithNZeros(5, hash) {
		fmt.Printf("Hash: %s\n", hex.EncodeToString(hash))
	}

	for i := 0; i < 10000000; i++ {
		hasher.Reset()
		saltedSecretKey = fmt.Sprintf("%s%v", secretKey, i)
		hasher.Write([]byte(saltedSecretKey))
		hash = hasher.Sum(nil)
		if isHashStartingWithNZeros(6, hash) {
			fmt.Printf("The number that produces the hash: %v\n", i)
			break
		}
	}
}
