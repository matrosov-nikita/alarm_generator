package generator

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/matrosov-nikita/smart-generator/pkg/config"

	"github.com/matrosov-nikita/smart-generator/generator/job"

	js "github.com/itimofeev/go-util/json"

	"github.com/google/uuid"
)

type Client interface {
	BulkInsert(items []js.Object) error
}

type supplierTask struct {
	detector     string
	eventsAmount int
	teamID       string
	domainID     int
}

type Supplier struct {
	*randomTimeGenerator
	serversCount      int
	teams             []string
	domains           []int
	client            Client
	detectorConfig    map[string]int
	timeGeneratorType string
	batchSize         int

	tasks chan *supplierTask
}

func New(client Client, cfg *config.Config) *Supplier {
	teams := make([]string, 0, cfg.TeamsCount)
	domains := make([]int, 0, cfg.TeamsCount)
	for i := 0; i < cfg.TeamsCount; i++ {
		teams = append(teams, uuid.New().String())
		domains = append(domains, i)
	}

	return &Supplier{
		randomTimeGenerator: newRandomTimeGenerator(cfg.StartDate, cfg.EndDate),
		serversCount:        cfg.ServersCount,
		client:              client,
		teams:               teams,
		domains:             domains,
		detectorConfig:      cfg.Detectors,
		tasks:               make(chan *supplierTask, len(cfg.Detectors)),
		timeGeneratorType:   cfg.TimeGeneratorType,
		batchSize:           cfg.BatchSize,
	}
}

func (s *Supplier) Run() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go s.worker(&wg)
	}

	for i, team := range s.teams {
		for detector, eventsAmount := range s.detectorConfig {
			j := &supplierTask{
				detector:     detector,
				eventsAmount: eventsAmount,
				teamID:       team,
				domainID:     s.domains[i],
			}

			s.tasks <- j
		}
	}

	close(s.tasks)
	wg.Wait()
}

func (s *Supplier) worker(wg *sync.WaitGroup) {
	for j := range s.tasks {
		if err := s.pushEvents(j); err != nil {
			log.Printf("failed to load events for job: %+v: %+v", err, j)
		}
	}

	wg.Done()
}

func (s *Supplier) pushEvents(task *supplierTask) error {
	batch := make([]js.Object, 0, s.batchSize)
	inBatch, total := 0, 0
	timeGenerator := NewGenerator(s.timeGeneratorType, s.from, s.to, int64(task.eventsAmount))
	for i := 0; i < task.eventsAmount; i++ {
		j := job.NewJob(task.detector, task.teamID, task.domainID, s.serversCount)
		serversEvents := j.GenerateEvents(timeGenerator.GetTime())
		batch = append(batch, serversEvents...)
		inBatch += len(serversEvents)

		if inBatch >= s.batchSize {
			if err := s.client.BulkInsert(batch); err != nil {
				return err
			}
			total += inBatch
			log.Printf("batch with [%d] entities for team [%s] was written, total: [%d]\n", inBatch, task.teamID, total)
			batch = make([]js.Object, 0, s.batchSize)
			inBatch = 0
		}
	}

	if inBatch > 0 {
		if err := s.client.BulkInsert(batch); err != nil {
			return err
		}
		total += inBatch
		fmt.Printf("final batch with [%d] entities for team [%s] was written, total: [%d]\n", inBatch, task.teamID, total)
	}

	return nil
}
