package job

import (
	"math/rand"
	"time"

	js "github.com/itimofeev/go-util/json"
	"github.com/matrosov-nikita/smart-generator/events"
)

type VehicleJob struct {
	teamID   string
	domainID int
}

func NewVehicleJob(teamID string, domainID int) *VehicleJob {
	return &VehicleJob{
		teamID:   teamID,
		domainID: domainID,
	}
}

func (j *VehicleJob) GenerateEvents(raiseTime time.Time) []js.Object {
	eventsGenerator := events.NewGenerator(j.teamID, j.domainID)
	cameraIn := 0
	cameraOut := 1
	plate := randString(8)
	hoursOnParking := rand.Intn(4) + 1
	leftParking := rand.Intn(2) != 0

	plateInEvent, _ := eventsGenerator.CreatePlateRecognizedEvent(&events.EventValues{
		RaiseTime: raiseTime,
		Plate:     plate,
		CameraID:  cameraIn,
	})
	plateOutEvent, _ := eventsGenerator.CreatePlateRecognizedEvent(
		&events.EventValues{
			RaiseTime: raiseTime.Add(time.Duration(hoursOnParking) * time.Hour),
			CameraID:  cameraOut,
			Plate:     plate,
		})

	vehicleEvents := []js.Object{plateInEvent, plateOutEvent}
	if leftParking {
		secondTimeIn := hoursOnParking + 2
		plateInSecondTimeEvent, _ := eventsGenerator.CreatePlateRecognizedEvent(
			&events.EventValues{
				RaiseTime: raiseTime.Add(time.Duration(secondTimeIn) * time.Hour),
				CameraID:  cameraIn,
				Plate:     plate,
			})
		vehicleEvents = append(vehicleEvents, plateInSecondTimeEvent)
	}

	return vehicleEvents
}
