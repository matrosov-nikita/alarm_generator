package generator

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/matrosov-nikita/smart-generator/events"
)

const batchSize = 200

type Client interface {
	BulkInsert(items []*events.Event) error
}

type Supplier struct {
	*timeGenerator
	serversCount   int
	teams          []string
	domains        []int
	client         Client
	detectorConfig map[string]int

	jobs chan *job
}

func New(client Client, startDate, endDate time.Time, serversCount, teamsCount int, detectorConfig map[string]int) *Supplier {
	teams := make([]string, 0, teamsCount)
	domains := make([]int, 0, teamsCount)
	for i := 0; i < teamsCount; i++ {
		teams = append(teams, uuid.New().String())
		domains = append(domains, i)
	}

	return &Supplier{
		timeGenerator:  NewTimeGenerator(startDate, endDate),
		serversCount:   serversCount,
		client:         client,
		teams:          teams,
		domains:        domains,
		detectorConfig: detectorConfig,
		jobs:           make(chan *job, len(detectorConfig)),
	}
}

func (s *Supplier) Run() {
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go s.worker(&wg)
	}

	for i, team := range s.teams {
		for detector, eventsAmount := range s.detectorConfig {
			s.jobs <- &job{
				detector:     detector,
				eventsAmount: eventsAmount,
				teamID:       team,
				domainID:     s.domains[i],
				serversCount: s.serversCount,
			}
		}
	}

	close(s.jobs)
	wg.Wait()
}

func (s *Supplier) worker(wg *sync.WaitGroup) {
	for j := range s.jobs {
		log.Printf("Got new job: %+v\n", j)
		if err := s.pushEvents(j); err != nil {
			log.Printf("failed to load events: %+v", err)
		}
	}

	wg.Done()
}

func (s *Supplier) pushEvents(job *job) error {
	batch := make([]*events.Event, 0, batchSize)
	inBatch, total := 0, 0
	for i := 0; i < job.eventsAmount; i++ {
		serversEvents := job.generateDetectorEvents(s.GetTime())
		batch = append(batch, serversEvents...)
		inBatch += len(serversEvents)

		if inBatch >= batchSize {
			if err := s.client.BulkInsert(batch); err != nil {
				return err
			}
			total += inBatch
			log.Printf("batch with [%d] entities for team [%s] was written, total: [%d]\n", inBatch, job.teamID, total)
			batch = make([]*events.Event, 0, batchSize)
			inBatch = 0
		}
	}

	if inBatch > 0 {
		if err := s.client.BulkInsert(batch); err != nil {
			return err
		}
		total += inBatch
		fmt.Printf("final batch with [%d] entities for team [%s] was written, total: [%d]\n", inBatch, job.teamID, total)
	}

	return nil
}
