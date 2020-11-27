package job

import (
	"log"
	"math/rand"
	"time"

	js "github.com/itimofeev/go-util/json"
	"github.com/matrosov-nikita/smart-generator/events"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

type DetectorJob struct {
	detector string
	teamID   string
	domainID int
}

func NewDetectorJob(detector, teamID string, domainID int) *DetectorJob {
	return &DetectorJob{
		detector: detector,
		teamID:   teamID,
		domainID: domainID,
	}
}

func (j *DetectorJob) GenerateEvents(raiseTime time.Time) []js.Object {
	var ev js.Object
	eventsGenerator := events.NewGenerator(j.teamID, j.domainID)
	var eventGenerators = map[string]func(values *events.EventValues) (js.Object, string){
		"faceAppeared":         eventsGenerator.CreateFaceAppearedEvent,
		"listed_lpr_detected":  eventsGenerator.CreateListedLprEvent,
		"listed_face_detected": eventsGenerator.CreateListedFaceEvent,
		"QueueDetected":        eventsGenerator.CreateQueueDetectedEvent,
		"QueueLength":          eventsGenerator.CreateQueueLengthEvent,
		"People":               eventsGenerator.CreatePeopleEvent,
		"bodyTemperature":      eventsGenerator.CreateBodyTemperatureEvent,
		"peopleDistance":       eventsGenerator.CreatePeopleDistanceEvent,
		"EvasionDetected":      eventsGenerator.CreateFaceMaskAbsenceEvent,
		"oneLine":              eventsGenerator.CreateOneLineEvent,
		"lotsObjects":          eventsGenerator.CreateLotsObjectsEvent,
	}

	switch j.detector {
	case "NoBodySegment", "NoHipsSegment", "NoHandSegment", "NoFootSegment", "NoOtherSegment", "NoHeadSegment":
		ev, _ = eventsGenerator.CreateEquipmentEvent(&events.EventValues{
			RaiseTime:    raiseTime,
			DetectorType: j.detector,
		})
	case "plateRecognized":
		ev, _ = eventsGenerator.CreatePlateRecognizedEvent(&events.EventValues{
			RaiseTime: raiseTime,
			Plate:     randString(8),
		})
	default:
		generator, ok := eventGenerators[j.detector]
		if !ok {
			log.Printf("event of %s is not supported", j.detector)
			return nil
		}
		ev, _ = generator(&events.EventValues{
			RaiseTime: raiseTime,
		})
	}

	return []js.Object{ev}
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
