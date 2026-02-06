package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type RandomGenerator struct {
	length int
}

func NewRandom(length int) (*RandomGenerator, error) {
	if length <= 0 {
		return nil, ErrInvalidLength
	}

	return &RandomGenerator{
		length: length,
	}, nil
}

func (g *RandomGenerator) NewAlias() (string, error) {
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
