package hashcash

import (
	"crypto/sha1"
	"errors"
	"fmt"
)

// https://en.wikipedia.org/wiki/Hashcash

var (
	ErrHashcashMaxIteration = errors.New("max iterations")
)

type Hashcash struct {
	Version  int
	Date     int64
	Resource string
	Zeros    int
	Rand     string
	Counter  int
}

func (hash Hashcash) Compute(maxIterations int) (Hashcash, error) {
	emptyHash := Hashcash{}

	for hash.Counter <= maxIterations || maxIterations <= 0 {
		sha1, err := hash.GetSha1()
		if err != nil {
			return emptyHash, err
		}

		if hash.Check(sha1) {
			return hash, nil
		}

		hash.Counter++
	}

	return hash, ErrHashcashMaxIteration
}

func (hash Hashcash) ToString() string {
	return fmt.Sprintf(
		"%d:%d:%d:%s::%s:%d",
		hash.Version, hash.Zeros, hash.Date, hash.Resource, hash.Rand, hash.Counter,
	)
}

func (hash Hashcash) Check(shaHash string) bool {
	if hash.Zeros > len(shaHash) {
		return false
	}

	for _, ch := range shaHash[0:hash.Zeros] {
		if ch != 0x30 {
			return false
		}
	}

	return true
}

func (hash Hashcash) GetSha1() (string, error) {
	shaHasher := sha1.New()

	_, err := shaHasher.Write([]byte(hash.ToString()))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", shaHasher.Sum(nil)), nil
}
