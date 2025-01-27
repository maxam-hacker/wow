package hashcash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnHashcash(t *testing.T) {
	var err error

	hc := Hashcash{
		Version:  1,
		Zeros:    1,
		Date:     123456,
		Resource: "maxam.hacker@gmail.com",
		Counter:  101,
		Rand:     "dgbernler",
	}

	hc, err = hc.Compute(1000)
	assert.NoError(t, err)
	assert.Equal(t, "1:1:123456:maxam.hacker@gmail.com::dgbernler:102", hc.ToString())

	sha, err := hc.GetSha1()
	assert.NoError(t, err)
	assert.Equal(t, "0d31fe344c1c0dc998fccade6c7b290044208353", sha)

	checked := hc.Check(sha)
	assert.Equal(t, true, checked)
}

func TestOnHashcashWith3Zeros(t *testing.T) {
	var err error

	hc := Hashcash{
		Version:  1,
		Zeros:    3,
		Date:     123456,
		Resource: "maxam.hacker@gmail.com",
		Counter:  101,
		Rand:     "dgbernler",
	}

	hc, err = hc.Compute(1000000)
	assert.NoError(t, err)
	assert.Equal(t, "1:3:123456:maxam.hacker@gmail.com::dgbernler:2028", hc.ToString())

	sha, err := hc.GetSha1()
	assert.NoError(t, err)
	assert.Equal(t, "0006b095a1b57cf653088b05de7c5d48a847c9c6", sha)

	checked := hc.Check(sha)
	assert.Equal(t, true, checked)
}

func TestOnHashcashWith3ZerosAndError(t *testing.T) {
	var err error

	hc := Hashcash{
		Version:  1,
		Zeros:    3,
		Date:     123456,
		Resource: "maxam.hacker@gmail.com",
		Counter:  101,
		Rand:     "dgbernler",
	}

	_, err = hc.Compute(1000)
	assert.Error(t, err)
}
