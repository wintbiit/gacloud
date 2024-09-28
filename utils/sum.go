package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"io"
)

// Md5Sum returns the MD5 checksum of the file, usually has a length of 32 characters.
func Md5Sum(f io.Reader) string {
	h := md5.New()
	_, err := io.Copy(h, f)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Md5SumBytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Sha256Sum returns the SHA-256 checksum of the file, usually has a length of 64 characters.
func Sha256Sum(f io.Reader) string {
	h := sha256.New()
	_, err := io.Copy(h, f)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha256SumBytes(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Fnv1a64Sum(reader io.Reader) string {
	h := fnv.New64a()
	_, err := io.Copy(h, reader)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Fnv1a64SumBytes(b []byte) string {
	h := fnv.New64a()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Fnv1a32Sum(reader io.Reader) string {
	h := fnv.New32a()
	_, err := io.Copy(h, reader)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Fnv1a32SumBytes(b []byte) string {
	h := fnv.New32a()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}
