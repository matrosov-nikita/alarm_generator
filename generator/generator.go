package generator

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/matrosov-nikita/alarm_generator/event"
)

const batchSize = 200

type DBClient interface {
	BulkInsert(items []event.Item) error
}

type Generator struct {
	*timeGenerator
	alertsCount  int
	serversCount int
	client       DBClient
}

func New(client DBClient, startDate, endDate time.Time, alertsCount, serversCount int) *Generator {
	return &Generator{
		timeGenerator: NewTimeGenerator(startDate, endDate),
		alertsCount:   alertsCount,
		serversCount:  serversCount,
		client:        client,
	}
}

func (g *Generator) LoadEvents() error {
	batch := make([]event.Item, 0, batchSize)
	inBatch, total := 0, 0
	for i := 0; i < g.alertsCount; i++ {
		serversEvents := g.prepareEvents()
		batch = append(batch, serversEvents...)
		inBatch += len(serversEvents)

		if inBatch >= batchSize {
			if err := g.client.BulkInsert(batch); err != nil {
				return err
			}
			total += inBatch
			log.Printf("batch with [%d] entities was written, total: [%d]\n", inBatch, total)
			batch = make([]event.Item, 0, batchSize)
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

func (g *Generator) prepareEvents() []event.Item {
	alertSeverities := []string{"True", "False", "Missed", "Suspicious"}
	alertRaiseTime := g.GetTime()
	alertStateSeverity := alertSeverities[rand.Intn(len(alertSeverities))]
	timeElapsedBeforeDetectorEvent := time.Second
	timeElapsedBeforeAlertStateChanged := 2 * time.Second

	faceAppearedTime := alertRaiseTime.Add(timeElapsedBeforeDetectorEvent)
	alertStateTime := alertRaiseTime.Add(timeElapsedBeforeAlertStateChanged)

	events := make([]event.Item, 0, 2*g.serversCount+1)
	faceAppearedEvent, faceAppearedEventID := event.DummyFaceAppearedEvent(faceAppearedTime, 0)

	events = append(events, faceAppearedEvent)
	alertID := uuid.New().String()
	alertStateEventID := uuid.New().String()
	for i := 0; i < g.serversCount; i++ {
		alertEvent, alertID := event.DummyAlertEvent(alertID, alertRaiseTime, faceAppearedEventID, i)
		alertStateEvent := event.DummyAlertEventState(alertStateEventID, alertStateTime, alertStateSeverity, alertID, i)
		events = append(events, []event.Item{alertEvent, alertStateEvent}...)
	}

	return events
}
