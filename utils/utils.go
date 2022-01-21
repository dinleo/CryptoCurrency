package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// ToByte encode interface data to []byte and return it
func ToByte(i interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleErr(encoder.Encode(i))
	return buffer.Bytes()
}

// FromBytes decode []byte to interface data
func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
}

// Hash interface by sha256 and return it
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	return hash
}
