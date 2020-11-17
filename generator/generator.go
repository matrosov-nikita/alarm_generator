package generator

import (
	"log"
	"time"

	"github.com/matrosov-nikita/alarm_generator/event"
)

const batchSize = 200

type DBClient interface {
	BulkInsert(items []event.Item) error
}

type Generator struct {
	*timeGenerator
	count  int
	client DBClient
}

func New(client DBClient, startDate, endDate time.Time, count int) *Generator {
	return &Generator{
		timeGenerator: NewTimeGenerator(startDate, endDate),
		count:         count,
		client:        client,
	}
}

func (g *Generator) LoadEvents() error {
	batch := make([]event.Item, 0, batchSize)
	inBatch, total := 0, 0
	for i := 0; i < g.count; i++ {
		// TODO: create alert events
		batch = append(batch, event.DummyFaceAppearedEvent(g.GetTime()))
		inBatch++
		if inBatch == batchSize {
			if err := g.client.BulkInsert(batch); err != nil {
				return err
			}
			total += batchSize
			log.Printf("batch with [%d] entities was written, total: [%d]\n", batchSize, total)
			batch = make([]event.Item, 0, batchSize)
			inBatch = 0
		}
	}

	if inBatch > 0 {
		if err := g.client.BulkInsert(batch); err != nil {
			return err
		}
		total += inBatch
		log.Printf("final batch with [%d] entities was written, total: [%d]\n", inBatch, total)
	}

	return nil
}
