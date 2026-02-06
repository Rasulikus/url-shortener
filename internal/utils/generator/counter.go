package generator

import (
	"sync/atomic"
)

type CounterGenerator struct {
	nextID atomic.Uint64
	secret uint64
	length int
	base   uint64
}

func NewCounter(startFrom uint64, secret uint64, length int) (*CounterGenerator, error) {
	if length <= 0 {
		return nil, ErrInvalidLength
	}

	g := &CounterGenerator{
		secret: secret,
		length: length,
		base:   uint64(len(alphabet)),
	}
	g.nextID.Store(startFrom)
	return g, nil
}

func (g *CounterGenerator) NewAlias() (string, error) {
	// Берём следующий уникальный ID атомарно
	id := g.nextID.Add(1)

	// Защита от простого угадывания
	mixed := id ^ g.secret

	// Кодировка числа в строку
	return encode(mixed, g.length, g.base)
}

func encode(n uint64, length int, base uint64) (string, error) {
	// Кол-во символов на выход, если не все заполняется то будут 0 вначале, в данном случае это 'a'
	out := make([]byte, length)

	// Заполняем справа налево: берём остаток от деления на мощность алфавита, делим на мощность алфавита и записывает результат в выход
	for i := length - 1; i >= 0; i-- {
		rem := n % base
		n /= base
		out[i] = alphabet[rem]
	}

	// Если осталось число после length операций, значит оно не поместилось в заданную длину
	if n != 0 {
		return "", ErrOverflow
	}

	return string(out), nil
}
