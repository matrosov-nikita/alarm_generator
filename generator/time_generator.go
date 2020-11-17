package generator

import (
	"math/rand"
	"time"
)

type timeGenerator struct {
	from       time.Time
	to         time.Time
	randSource rand.Source
}

func NewTimeGenerator(from, to time.Time) *timeGenerator {
	return &timeGenerator{
		from:       from,
		to:         to,
		randSource: rand.NewSource(time.Now().UnixNano()),
	}
}

func (gen *timeGenerator) GetTime() time.Time {
	min := gen.from.Unix()
	max := gen.to.Unix()

	sec := rand.Int63n(max-min) + min
	return time.Unix(sec, 0)
}
