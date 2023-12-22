package hashing

import (
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

type PathKey struct {
	PathName string
	Original string
}

func HashFile(key string) PathKey {
	hash := sha3.Sum256([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

func (file PathKey) Filename() string {
	return fmt.Sprintf("%s/%s", file.PathName, file.Original)
}
