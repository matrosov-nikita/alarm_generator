package generator

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	js "github.com/itimofeev/go-util/json"

	"github.com/google/uuid"
)

const batchSize = 200

type Client interface {
	BulkInsert(items []js.Object) error
}

type Supplier struct {
	*randomTimeGenerator
	serversCount      int
	teams             []string
	domains           []int
	client            Client
	detectorConfig    map[string]int
	timeGeneratorType string

	jobs chan *job
}

func New(client Client, startDate, endDate time.Time, serversCount, teamsCount int, timeGeneratorType string, detectorConfig map[string]int) *Supplier {
	teams := make([]string, 0, teamsCount)
	domains := make([]int, 0, teamsCount)
	for i := 0; i < teamsCount; i++ {
		teams = append(teams, uuid.New().String())
		domains = append(domains, i)
	}

	return &Supplier{
		randomTimeGenerator: newRandomTimeGenerator(startDate, endDate),
		serversCount:        serversCount,
		client:              client,
		teams:               teams,
		domains:             domains,
		detectorConfig:      detectorConfig,
		jobs:                make(chan *job, len(detectorConfig)),
		timeGeneratorType:   timeGeneratorType,
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
		if err := s.pushEvents(j); err != nil {
			log.Printf("failed to load events for job: %+v: %+v", err, j)
		}
	}

	wg.Done()
}

func (s *Supplier) pushEvents(job *job) error {
	batch := make([]js.Object, 0, batchSize)
	inBatch, total := 0, 0
	timeGenerator := NewGenerator(s.timeGeneratorType, s.from, s.to, int64(job.eventsAmount))
	for i := 0; i < job.eventsAmount; i++ {
		serversEvents := job.generateDetectorEvents(timeGenerator.GetTime())
		batch = append(batch, serversEvents...)
		inBatch += len(serversEvents)

		if inBatch >= batchSize {
			if err := s.client.BulkInsert(batch); err != nil {
				return err
			}
			total += inBatch
			log.Printf("batch with [%d] entities for team [%s] was written, total: [%d]\n", inBatch, job.teamID, total)
			batch = make([]js.Object, 0, batchSize)
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
