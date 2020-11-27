package job

import (
	"time"

	js "github.com/itimofeev/go-util/json"
	"github.com/matrosov-nikita/smart-generator/events"
)

type AlertJob struct {
	teamID       string
	domainID     int
	serversCount int
}

func NewAlertJob(teamID string, domainID, serversCount int) *AlertJob {
	return &AlertJob{
		teamID:       teamID,
		domainID:     domainID,
		serversCount: serversCount,
	}
}

func (j *AlertJob) GenerateEvents(alertRaiseTime time.Time) []js.Object {
	eventsGenerator := events.NewGenerator(j.teamID, j.domainID)
	timeElapsedBeforeDetectorEvent := time.Second
	timeElapsedBeforeAlertStateChanged := 2 * time.Second

	faceAppearedTime := alertRaiseTime.Add(timeElapsedBeforeDetectorEvent)
	alertStateTime := alertRaiseTime.Add(timeElapsedBeforeAlertStateChanged)

	alertsEvents := make([]js.Object, 0, 2*j.serversCount+1)
	faceAppearedEvent, faceAppearedEventID := eventsGenerator.CreateFaceAppearedEvent(&events.EventValues{
		RaiseTime: faceAppearedTime,
	})

	alertsEvents = append(alertsEvents, faceAppearedEvent)
	for serverID := 0; serverID < j.serversCount; serverID++ {
		alertEvent, alertID := eventsGenerator.CreateAlertEvent(
			&events.EventValues{
				RaiseTime:       alertRaiseTime,
				ServerID:        serverID,
				DetectorEventID: faceAppearedEventID,
			})
		alertStateEvent := eventsGenerator.CreateAlertEventState(
			&events.EventValues{
				RaiseTime: alertStateTime,
				ServerID:  serverID,
				AlertID:   alertID,
			})
		alertsEvents = append(alertsEvents, []js.Object{alertEvent, alertStateEvent}...)
	}

	return alertsEvents
}
