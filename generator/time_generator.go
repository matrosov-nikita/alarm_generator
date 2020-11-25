package generator

import (
	"math/rand"
	"time"
)

type TimeGenerator interface {
	GetTime() time.Time
}

func NewGenerator(generatorType string, from, to time.Time, eventsCount int64) TimeGenerator {
	switch generatorType {
	case "normal":
		return newNormalTimeGenerator(from, to, eventsCount)
	case "random":
		return newRandomTimeGenerator(from, to)
	default:
		return nil
	}
}

type randomTimeGenerator struct {
	from       time.Time
	to         time.Time
	randSource rand.Source
}

func newRandomTimeGenerator(from, to time.Time) *randomTimeGenerator {
	return &randomTimeGenerator{
		from:       from,
		to:         to,
		randSource: rand.NewSource(time.Now().UnixNano()),
	}
}

func (gen *randomTimeGenerator) GetTime() time.Time {
	min := gen.from.Unix()
	max := gen.to.Unix()

	sec := rand.Int63n(max-min) + min
	return time.Unix(sec, 0)
}

type normalTimeGenerator struct {
	from       time.Time
	to         time.Time
	ticksCount int64
	offset     int64
}

func newNormalTimeGenerator(from, to time.Time, ticksCount int64) *normalTimeGenerator {
	return &normalTimeGenerator{
		from:       from,
		to:         to,
		ticksCount: ticksCount,
	}
}

func (gen *normalTimeGenerator) GetTime() time.Time {
	defer func() {
		gen.offset += 1
	}()

	min := gen.from.Unix()
	max := gen.to.Unix()

	interval := (max - min) / gen.ticksCount
	return time.Unix(min+interval*gen.offset, 0)
}
