package job

import (
	"time"

	js "github.com/itimofeev/go-util/json"
)

type Job interface {
	GenerateEvents(raiseTime time.Time) []js.Object
}

func NewJob(detector, teamID string, domainID, serversCount int) Job {
	switch detector {
	case "alerts":
		return NewAlertJob(teamID, domainID, serversCount)
	case "vehicles":
		return NewVehicleJob(teamID, domainID)
	default:
		return NewDetectorJob(detector, teamID, domainID)
	}
}
