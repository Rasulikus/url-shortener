package alias

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const (
	DefaultLength = 10
	alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var ErrInvalidLength = errors.New("alias generator: invalid length")

type Generator struct {
	length int
}

func New(length int) (*Generator, error) {
	if length <= 0 {
		return nil, ErrInvalidLength
	}

	return &Generator{
		length: length,
	}, nil
}

func (g *Generator) NewAlias() (string, error) {
	b := make([]byte, g.length)
	lenAlpha := big.NewInt(int64(len(alphabet)))

	for i := 0; i < g.length; i++ {
		n, err := rand.Int(rand.Reader, lenAlpha)
		if err != nil {
			return "", fmt.Errorf("alias generator: rand.Int: %w", err)
		}
		b[i] = alphabet[n.Int64()]
	}

	return string(b), nil
}
