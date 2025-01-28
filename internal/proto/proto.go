package proto

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"
	"wow/internal/hashcash"
)

type RequestType string
type ResponseType string

const (
	RequestActionType          RequestType = "action"
	RequestActionExecutionType RequestType = "execution"

	ResponseOnActionType          ResponseType = "on_action"
	ResponseOnActionExecutionType ResponseType = "on_execution"
)

var (
	ErrHashValidationZeros = errors.New("wrong number of zeros")
	ErrHashValidationCheck = errors.New("error in checking number of zeros")
)

type ClientMeta struct {
	ClientId      string
	ClientVersion string
}

type Request struct {
	Meta   ClientMeta
	Type   RequestType
	LineId int
	Hash   hashcash.Hashcash
}

type Response struct {
	Type   ResponseType
	Hash   hashcash.Hashcash
	LineId int
	Result string
}

func NewRequestAction(meta ClientMeta) Request {
	return Request{
		Type: RequestActionType,
		Meta: meta,
	}
}

func NewRequestActionExecution(meta ClientMeta, lineId int, hash hashcash.Hashcash) Request {
	return Request{
		Meta:   meta,
		Type:   RequestActionExecutionType,
		LineId: lineId,
		Hash:   hash,
	}
}

func NewResponseOnAction(workLoadFactor int16) Response {
	var randBytes [32]byte

	n, _ := rand.Read(randBytes[:])

	c, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		c = big.NewInt(1)
	}

	return Response{
		Type: ResponseOnActionType,
		Hash: hashcash.Hashcash{
			Version:  1,
			Zeros:    int(workLoadFactor),
			Date:     time.Now().UTC().Unix(),
			Resource: "action-execute",
			Rand:     string(randBytes[:n]),
			Counter:  int(c.Int64()),
		},
	}
}

func NewResponseOnActionExecution(lineId int, result string) Response {
	return Response{
		Type:   ResponseOnActionExecutionType,
		LineId: lineId,
		Result: result,
	}
}

func (request *Request) Validate(workLoadFactor int16) (bool, error) {
	if request.Hash.Zeros != int(workLoadFactor) {
		return false, ErrHashValidationZeros
	}

	sha1, err := request.Hash.GetSha1()
	if err != nil {
		return false, ErrHashValidationCheck
	}

	if !request.Hash.Check(sha1) {
		return false, ErrHashValidationCheck
	}

	return true, nil
}
