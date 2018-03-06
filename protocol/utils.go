package protocol

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/crypto/sha3"
	"math"
	"math/rand"
	"time"
)

//Serializes the input in big endian and returns the sha3 hash function applied on ths input
func SerializeHashContent(data interface{}) (hash [32]byte) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, data)

	return sha3.Sum256(buf.Bytes())
}

func Encode(data [][]byte, sliceSize int) []byte {
	encodedData := make([]byte, len(data)*sliceSize)
	index := 0

	for _, item := range data {
		copy(encodedData[index:index+sliceSize], item)
		index += sliceSize
	}

	return encodedData
}

func Decode(encodedData []byte, sliceSize int) (data [][]byte) {
	index := 0
	cnt := 1

	for len(encodedData) >= cnt*sliceSize {
		data = append(data, encodedData[index:index+sliceSize])

		index += sliceSize
		cnt++
	}

	return data
}

func calculateBloomFilterParams(n float64, p float64) (uint, uint) {
	mFloat := math.Ceil((n * math.Log(p)) / math.Log(1/math.Pow(2.00, math.Log(2.00))))
	kFloat := int(math.Floor(math.Log(2.00) * mFloat / n))

	m := uint(mFloat)
	k := uint(kFloat)

	return m, k
}

func CreateRandomSeed() [32]byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seed [32]byte
	for i := range seed {
		seed[i] = chars[r.Intn(len(chars))]
	}
	return seed
}
