package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
)

// Md5Sum returns the MD5 checksum of the file, usually has a length of 32 characters.
func Md5Sum(f io.Reader) string {
	h := md5.New()
	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Sha256Sum returns the SHA-256 checksum of the file, usually has a length of 64 characters.
func Sha256Sum(f io.Reader) string {
	h := sha256.New()
	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))
}
