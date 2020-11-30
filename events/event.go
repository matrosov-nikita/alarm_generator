package events

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	js "github.com/itimofeev/go-util/json"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

type Generator struct {
	teamID   string
	domainID int
}

type EventValues struct {
	RaiseTime       time.Time
	ServerID        int
	CameraID        int
	DetectorEventID string
	AlertID         string
	Plate           string
	DetectorType    string
}

func NewGenerator(teamID string, domainID int) *Generator {
	return &Generator{
		teamID:   teamID,
		domainID: domainID,
	}
}

func (eg *Generator) CreateAlertEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["version"] = 1
	evt["type"] = "alert"
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, false)
	evt["initiator"] = "root"
	evt["initiator_type"] = "AIT_USER"
	evt["reviewer"] = "root"
	evt["reason_mask"] = 4
	evt["detector_event_id"] = values.DetectorEventID
	evt["detector_event_type"] = "faceAppeared"
	evt["macro_event_id"] = uuid.New().String()
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateAlertEventState(values *EventValues) js.Object {
	alertSeverities := []string{"True", "False", "Missed", "Suspicious"}
	alertStateSeverity := alertSeverities[rand.Intn(len(alertSeverities))]
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	evt["id"] = uuid.New().String()
	evt["version"] = 1
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["type"] = "alert_state"
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, false)
	evt["severity"] = alertStateSeverity
	evt["reviewer_type"] = "RT_USER"
	evt["reviewer"] = "root"
	evt["state"] = "ST_CLOSED"
	evt["alert_id"] = values.AlertID
	addDummyBookmark(evt, values.ServerID, values.AlertID)
	enrichTime(evt)
	return evt
}

func (eg *Generator) CreatePeopleEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["type"] = "detector"
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	evt["detector_type"] = "People"
	evt["detector_people_state"] = "in"
	evt["phase"] = "happened"
	addDummySource(evt, values.ServerID, values.CameraID, true)
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateQueueDetectedEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["type"] = "detector"
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	evt["detector_type"] = "QueueDetected"
	evt["detector_queue_max"] = 3
	evt["detector_queue_min"] = 3
	evt["phase"] = "happened"
	addDummySource(evt, values.ServerID, values.CameraID, true)
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateQueueLengthEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["type"] = "detector"
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	evt["detector_type"] = "QueueLength"
	evt["detector_queue_max"] = 3
	evt["detector_queue_min"] = 3
	evt["phase"] = "happened"
	addDummySource(evt, values.ServerID, values.CameraID, true)
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreatePlateRecognizedEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "plateRecognized"
	evt["detector_lpr_country"] = "ru"
	evt["detector_lpr_direction"] = 1
	evt["detector_lpr_plate"] = values.Plate
	evt["phase"] = "happened"
	evt["detector_lpr_best_utc"] = utc
	evt["detector_lpr_begin_datetime"] = datetime
	evt["detector_lpr_begin_utc"] = utc
	evt["detector_lpr_begin_datetime"] = datetime
	evt["detector_lpr_end_datetime"] = datetime
	evt["detector_lpr_end_utc"] = utc
	addDummyRectangle(evt)
	evt["recognition_quality"] = 0.45
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateListedLprEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "listed_lpr_detected"
	evt["detector_lpr_plate"] = randString(8)
	evt["phase"] = "happened"
	evt["detector_listedItem_list_id"] = "25580957-8639-459d-86a1-724f4e772956"
	evt["detector_listedItem_item_id"] = uuid.New().String()
	evt["detector_listedItem_matched_event_id"] = uuid.New().String()
	evt["detector_listedItem_matched_event_time_datetime"] = datetime
	evt["detector_listedItem_matched_event_time_utc"] = utc

	addDummyRectangle(evt)
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateListedFaceEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "listed_face_detected"
	evt["phase"] = "happened"
	evt["detector_listedItem_list_id"] = "4b63054f-4b82-40b9-88dc-800ae26e76f9"
	evt["detector_listedItem_item_id"] = uuid.New().String()
	evt["detector_listedItem_matched_event_id"] = uuid.New().String()
	evt["detector_listedItem_matched_event_time_datetime"] = datetime
	evt["detector_listedItem_matched_event_time_utc"] = utc
	evt["detector_listedFace_score"] = 0.5
	addDummyRectangle(evt)
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateBodyTemperatureEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "bodyTemperature"
	evt["phase"] = "happened"
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreatePeopleDistanceEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "peopleDistance"
	evt["phase"] = "happened"
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateLotsObjectsEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "lotsObjects"
	evt["phase"] = "happened"
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateFaceMaskAbsenceEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "EvasionDetected"
	evt["phase"] = "happened"
	evt["multi_phase_id"] = uuid.New().String()
	enrichTime(evt)
	addDummyRectangle(evt)
	return evt, id
}

func (eg *Generator) CreateEquipmentEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = values.DetectorType
	evt["phase"] = "happened"
	enrichTime(evt)
	addDummyRectangle(evt)
	return evt, id
}

func (eg *Generator) CreateFaceAppearedEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["type"] = "detector"
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "faceAppeared"
	evt["detector_face_age"] = 34
	evt["detector_face_gender"] = "male"
	evt["detector_face_time_begin"] = datetime
	evt["detector_face_time_begin_utc"] = utc
	evt["detector_queue_max"] = 2
	evt["recognition_quality"] = 0.45
	evt["detector_listedItem_list_id"] = "bfefb72f-235a-414f-afe7-5303f9d2e50e"
	evt["detector_listedItem_item_id"] = "7a9b3866-8901-4260-9826-470ce34c4219"
	addDummyRectangle(evt)
	evt["detector_listedItem_matched_event_time_utc"] = utc
	evt["phase"] = "happened"
	evt["multi_phase_id"] = uuid.New().String()
	enrichTime(evt)
	return evt, id
}

func (eg *Generator) CreateOneLineEvent(values *EventValues) (js.Object, string) {
	datetime, utc := timeToLocalAndUTCTimestamps(values.RaiseTime)
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["type"] = "detector"
	evt["domain__id"] = eg.domainID
	evt["team__id"] = eg.teamID
	evt["datetime"] = datetime
	evt["time_utc"] = utc
	addDummySource(evt, values.ServerID, values.CameraID, true)
	evt["detector_type"] = "oneLine"
	evt["phase"] = "happened"
	evt["multi_phase_id"] = uuid.New().String()
	enrichTime(evt)
	addDummyRectangle(evt)
	return evt, id
}

func addDummySource(event js.Object, serverID int, cameraID int, withDetector bool) {
	event["server_id"] = fmt.Sprintf("SERVER%d", serverID)
	event["server_name"] = fmt.Sprintf("someServer:%d", serverID)

	event["camera_id"] = fmt.Sprintf("SERVER%d/DeviceIpint.%d/SourceEndpoint.video:0:0", serverID, cameraID)
	event["camera_name"] = fmt.Sprintf("someCamera:%d", cameraID)

	if withDetector {
		event["detector_id"] = "HOST/AVDetector.1/EventSupplier"
		event["detector_name"] = "someDetector"
	}
}

func addDummyRectangle(event js.Object) {
	event["rectangle_h"] = 0.3
	event["rectangle_w"] = 0.3
	event["rectangle_x"] = 0.4
	event["rectangle_y"] = 0.3
	event["rectangle_index"] = 1
}

func addDummyBookmark(evt js.Object, serverID int, alertID string) {
	evtTime := time.Now()
	datetime := evt.GetFieldAsTime("datetime", "2006-01-02T15:04:05")
	evt["bookmark_time_datetime"] = datetime
	evt["bookmark_time_utc"] = evtTime.UTC()
	evt["bookmark_server_id"] = fmt.Sprintf("SERVER%d", serverID)
	evt["bookmark_server_name"] = fmt.Sprintf("someServer%d", serverID)
	evt["bookmark_camera_id"] = "SERVER0/DeviceIpint.1/SourceEndpoint.video:0:0"
	evt["bookmark_camera_name"] = "someCamera"
	evt["bookmark_id"] = uuid.New().String()
	evt["bookmark_message"] = "test message"
	evt["bookmark_is_protected"] = 1
	evt["bookmark_user"] = "root"
	evt["bookmark_alert_id"] = alertID
	evt["bookmark_group_id"] = uuid.New().String()
	evt["bookmark_boundary_x"] = 0.2
	evt["bookmark_boundary_y"] = 0.3
	evt["bookmark_boundary_w"] = 0.4
	evt["bookmark_boundary_h"] = 0.5
	evt["bookmark_boundary_index"] = 1
	evt["bookmark_geometry_alpha"] = 147
	evt["bookmark_geometry_id"] = uuid.New().String()
	evt["bookmark_geometry_type"] = "PT_ELLIPSE"
	evt["bookmark_geometry_ellipse_center_x"] = 2.0
	evt["bookmark_geometry_ellipse_center_y"] = 4.0
	evt["bookmark_geometry_ellipse_yr"] = 4.5
	evt["bookmark_geometry_ellipse_xr"] = 4.5
	evt["bookmark_range_time_begin"] = evtTime.UTC()
	evt["bookmark_range_time_end"] = evtTime.Add(3 * time.Second).UTC()
}

func enrichTime(event js.Object) {
	datetime := event.GetField("datetime")
	localTime, _ := datetime.(Time)
	event.Put("weekday", strings.ToLower(localTime.Weekday().String()))
	event.Put("year", localTime.Year())
	event.Put("month", int(localTime.Month()))
	event.Put("date", localTime.Format("2006-01-02"))
	event.Put("time", localTime.Format("15:04:05.000"))
	event.Put("hour", localTime.Hour())
	event.Put("nanoseconds", localTime.Nanosecond())
	_, utcOffset := localTime.Local().Zone()
	utctime := event.GetField("time_utc")
	utc, _ := utctime.(Time)
	event.Put("utc_nanoseconds", utc.Nanosecond())
	event.Put("utc_offset", utcOffset)

	dayHour, dayMin, daySec := localTime.Clock()
	secondsInDay := dayHour*3600 + dayMin*60 + daySec
	minOfDayX5 := secondsInDay / (5 * 60)
	minOfDayX10 := secondsInDay / (10 * 60)
	minOfDayX15 := secondsInDay / (15 * 60)
	minOfDayX30 := secondsInDay / (30 * 60)
	event.Put("min_of_day", fmt.Sprintf("%02d:%02d", dayHour, dayMin))
	event.Put("min_of_day_x5", fmt.Sprintf("%02d:%02d", minOfDayX5*5/60, minOfDayX5*5%60))
	event.Put("min_of_day_x10", fmt.Sprintf("%02d:%02d", minOfDayX10*10/60, minOfDayX10*10%60))
	event.Put("min_of_day_x15", fmt.Sprintf("%02d:%02d", minOfDayX15*15/60, minOfDayX15*15%60))
	event.Put("min_of_day_x30", fmt.Sprintf("%02d:%02d", minOfDayX30*30/60, minOfDayX30*30%60))

	hourX2 := dayHour / 2
	hourX4 := dayHour / 4
	event.Put("hour_x2", hourX2*2)
	event.Put("hour_x4", hourX4*4)
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func timeToLocalAndUTCTimestamps(raiseTime time.Time) (Time, Time) {
	evtTimeStr := raiseTime.Local().Format("2006-01-02T15:04:05.000000")
	datetime, _ := time.Parse("2006-01-02T15:04:05", evtTimeStr)
	return Time{datetime}, Time{raiseTime.UTC()}
}
