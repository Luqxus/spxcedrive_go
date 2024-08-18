package main

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

type PathTransformFunc func(key string) string

func DefaultPathTransformFunc(key string) string {
	parts := strings.Split(key, "/")
	h := sha1.Sum([]byte(parts[len(parts)-1]))
	return hex.EncodeToString(h[:])
}
