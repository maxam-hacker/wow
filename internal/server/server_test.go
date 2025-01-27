package server

import (
	"testing"
	"wow/internal/pkg/transport/workers"

	"github.com/stretchr/testify/assert"
)

func TestOnEmptyLoadBalancer(t *testing.T) {
	s := Server{}
	err := s.Start()
	assert.Equal(t, ErrEmptyLoadBalancer, err)
}

func TestMessageHandler(t *testing.T) {
	s := Server{}

	err := s.messageHandler(nil, nil, 0, nil)
	assert.Equal(t, ErrEmptyWriter, err)

	err = s.messageHandler(nil, &workers.Writer{}, 0, nil)
	assert.Equal(t, ErrEmptyCloser, err)

	err = s.messageHandler(nil, &workers.Writer{}, 0, &workers.Closer{})
	assert.Error(t, err)
}
