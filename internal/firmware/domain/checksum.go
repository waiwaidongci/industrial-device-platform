package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Checksum(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func VerifyChecksum(data []byte, want string) error {
	if Checksum(data) != want {
		return fmt.Errorf("firmware checksum mismatch")
	}
	return nil
}
func Chunk(data []byte, size int) [][]byte {
	if size <= 0 {
		size = 64 * 1024
	}
	out := [][]byte{}
	for len(data) > 0 {
		n := size
		if n > len(data) {
			n = len(data)
		}
		out = append(out, append([]byte(nil), data[:n]...))
		data = data[n:]
	}
	return out
}

func CloneBytes(data []byte) []byte { return data }
