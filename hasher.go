package main

import (
	"crypto/sha256"
	"io"
)

type Hasher interface {
	Hash(r io.Reader) ([]byte, error)
}

type DefaultHasher struct{}

func (hasher *DefaultHasher) Hash(r io.Reader) ([]byte, error) {
	h := sha256.New()

	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
