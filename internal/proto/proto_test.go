package proto

import (
	"testing"
	"wow/internal/hashcash"

	"github.com/stretchr/testify/assert"
)

func TestOnNewResponse(t *testing.T) {
	responseOnAction := NewResponseOnAction(10)
	assert.NotEqual(t, 0, responseOnAction.Hash.Counter)
	assert.NotEqual(t, "", responseOnAction.Hash.Rand)
}

func TestOnValidation(t *testing.T) {
	hc, err := hashcash.Hashcash{
		Version:  1,
		Zeros:    3,
		Date:     123456,
		Resource: "maxam.hacker@gmail.com",
		Counter:  101,
		Rand:     "dgbernler",
	}.Compute(1000000)
	assert.NoError(t, err)

	r := NewRequestActionExecution(
		ClientMeta{},
		3,
		hc,
	)

	ok, err := r.Validate(3)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
}
