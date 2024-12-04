package utils

import (
	"crypto/rand"
	"math/big"
)

type RandomGenerator interface {
	GenerateNumber(min, max int64) (int64, error)
}

type cryptoRandom struct{}

func NewRandomGenerator() RandomGenerator {
	return &cryptoRandom{}
}

func (r *cryptoRandom) GenerateNumber(min, max int64) (int64, error) {
	// Криптографически безопасная генерация случайных чисел
	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return 0, err
	}
	return n.Int64() + min, nil
}

//ПЕПЕПУПУНАНА
