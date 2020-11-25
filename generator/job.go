package generator

import (
	"log"
	"math/rand"
	"time"

	js "github.com/itimofeev/go-util/json"

	"github.com/google/uuid"
	"github.com/matrosov-nikita/smart-generator/events"
)

// TODO: create different job types (e.g. Detector Job, Alerts Job, Vehicle Job)
type job struct {
	detector     string
	eventsAmount int
	serversCount int
	teamID       string
	domainID     int
}

func (j *job) generateDetectorEvents(raiseTime time.Time) []js.Object {
	var ev js.Object
	eventsGenerator := events.NewGenerator(j.teamID, j.domainID)
	var eventGenerators = map[string]func(time.Time, int) (js.Object, string){
		"faceAppeared":         eventsGenerator.CreateFaceAppearedEvent,
		"plateRecognized":      eventsGenerator.CreatePlateRecognizedEvent,
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
	case "alerts":
		alerts := j.generateAlerts(raiseTime)
		log.Println("LEN", len(alerts), j.serversCount)
		return alerts
	case "NoBodySegment", "NoHipsSegment", "NoHandSegment", "NoFootSegment", "NoOtherSegment", "NoHeadSegment":
		ev, _ = eventsGenerator.CreateEquipmentEvent(raiseTime, 0, j.detector)
	default:
		generator, ok := eventGenerators[j.detector]
		if !ok {
			log.Printf("event of %s is not supported", j.detector)
			return nil
		}
		ev, _ = generator(raiseTime, 0)
	}

	return []js.Object{ev}
}

func (j *job) generateAlerts(alertRaiseTime time.Time) []js.Object {
	eventsGenerator := events.NewGenerator(j.teamID, j.domainID)

	alertSeverities := []string{"True", "False", "Missed", "Suspicious"}
	alertStateSeverity := alertSeverities[rand.Intn(len(alertSeverities))]
	timeElapsedBeforeDetectorEvent := time.Second
	timeElapsedBeforeAlertStateChanged := 2 * time.Second

	faceAppearedTime := alertRaiseTime.Add(timeElapsedBeforeDetectorEvent)
	alertStateTime := alertRaiseTime.Add(timeElapsedBeforeAlertStateChanged)

	alertsEvents := make([]js.Object, 0, 2*j.serversCount+1)
	faceAppearedEvent, faceAppearedEventID := eventsGenerator.CreateFaceAppearedEvent(faceAppearedTime, 0)

	alertsEvents = append(alertsEvents, faceAppearedEvent)
	alertID := uuid.New().String()
	for i := 0; i < j.serversCount; i++ {
		alertEvent, alertID := eventsGenerator.CreateAlertEvent(alertID, alertRaiseTime, faceAppearedEventID, i)
		alertStateEvent := eventsGenerator.CreateAlertEventState(alertStateTime, alertStateSeverity, alertID, i)
		alertsEvents = append(alertsEvents, []js.Object{alertEvent, alertStateEvent}...)
	}

	return alertsEvents
}
