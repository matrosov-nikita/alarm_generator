package generator

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/matrosov-nikita/smart-generator/events"
)

const batchSize = 200

type DBClient interface {
	BulkInsert(items []*events.Event) error
}

type Generator struct {
	*timeGenerator
	serversCount   int
	client         DBClient
	detectorConfig map[string]int

	jobs chan *job
}

type job struct {
	detector     string
	eventsAmount int
}

func New(client DBClient, startDate, endDate time.Time, serversCount int, detectorConfig map[string]int) *Generator {
	return &Generator{
		timeGenerator:  NewTimeGenerator(startDate, endDate),
		serversCount:   serversCount,
		client:         client,
		detectorConfig: detectorConfig,
		jobs:           make(chan *job, len(detectorConfig)),
	}
}

func (g *Generator) Run() {
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go g.worker(&wg)
	}

	for detector, eventsAmount := range g.detectorConfig {
		g.jobs <- &job{
			detector:     detector,
			eventsAmount: eventsAmount,
		}
	}

	close(g.jobs)
	wg.Wait()
}

func (g *Generator) worker(wg *sync.WaitGroup) {
	for j := range g.jobs {
		log.Printf("Got new job: %+v\n", j)
		if err := g.loadEvents(j.detector, j.eventsAmount); err != nil {
			log.Printf("failed to load events: %+v", err)
		}
	}

	wg.Done()
}

func (g *Generator) loadEvents(detectorType string, eventsAmount int) error {
	batch := make([]*events.Event, 0, batchSize)
	inBatch, total := 0, 0
	for i := 0; i < eventsAmount; i++ {
		serversEvents := g.generateEvents(detectorType)
		batch = append(batch, serversEvents...)
		inBatch += len(serversEvents)

		if inBatch >= batchSize {
			if err := g.client.BulkInsert(batch); err != nil {
				return err
			}
			total += inBatch
			log.Printf("batch with [%d] entities was written, total: [%d]\n", inBatch, total)
			batch = make([]*events.Event, 0, batchSize)
			inBatch = 0
		}
	}

	if inBatch > 0 {
		if err := g.client.BulkInsert(batch); err != nil {
			return err
		}
		total += inBatch
		fmt.Printf("final batch with [%d] entities was written, total: [%d]\n", inBatch, total)
	}

	return nil
}

func (g *Generator) generateEvents(detectorType string) []*events.Event {
	raiseTime := g.GetTime()
	var ev *events.Event
	var eventGenerators = map[string]func(time.Time, int) (*events.Event, string){
		"faceAppeared":        events.NewFaceAppearedEvent,
		"plateRecognized":     events.NewPlateRecognizedEvent,
		"listed_lpr_detected": events.NewListedLprEvent,
		"QueueDetected":       events.NewQueueDetectedEvent,
		"People":              events.NewPeopleEvent,
	}

	switch detectorType {
	case "alerts":
		return g.generateAlerts(raiseTime)
	default:
		generator, ok := eventGenerators[detectorType]
		if !ok {
			log.Printf("event of %s is not supported", detectorType)
			return nil
		}
		ev, _ = generator(raiseTime, 0)
	}

	return []*events.Event{ev}
}

func (g *Generator) generateAlerts(alertRaiseTime time.Time) []*events.Event {
	alertSeverities := []string{"True", "False", "Missed", "Suspicious"}
	alertStateSeverity := alertSeverities[rand.Intn(len(alertSeverities))]
	timeElapsedBeforeDetectorEvent := time.Second
	timeElapsedBeforeAlertStateChanged := 2 * time.Second

	faceAppearedTime := alertRaiseTime.Add(timeElapsedBeforeDetectorEvent)
	alertStateTime := alertRaiseTime.Add(timeElapsedBeforeAlertStateChanged)

	alertsEvents := make([]*events.Event, 0, 2*g.serversCount+1)
	faceAppearedEvent, faceAppearedEventID := events.NewFaceAppearedEvent(faceAppearedTime, 0)

	alertsEvents = append(alertsEvents, faceAppearedEvent)
	alertID := uuid.New().String()
	alertStateEventID := uuid.New().String()
	for i := 0; i < g.serversCount; i++ {
		alertEvent, alertID := events.NewAlertEvent(alertID, alertRaiseTime, faceAppearedEventID, i)
		alertStateEvent := events.NewAlertEventState(alertStateEventID, alertStateTime, alertStateSeverity, alertID, i)
		alertsEvents = append(alertsEvents, []*events.Event{alertEvent, alertStateEvent}...)
	}

	return alertsEvents
}
