package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	js "github.com/itimofeev/go-util/json"
)

const (
	domainID     = 45
	teamID       = "d72452fc-1587-4786-b28b-9fe449d026ed"
	macroEventID = "19178c2a-894b-4f06-9f0f-213ec15ebaa3"
)

var alertsColumns = []string{"id",
	"team__id",
	"domain__id",
	"type",
	"time_utc",
	"alert_id",
	"version",
	"datetime",
	"year",
	"month",
	"date",
	"time",
	"hour",
	"hour_x2",
	"hour_x4",
	"min_of_day",
	"min_of_day_x5",
	"min_of_day_x10",
	"min_of_day_x15",
	"min_of_day_x30",
	"nanoseconds",
	"utc_nanoseconds",
	"weekday",
	"server_id",
	"server_name",
	"camera_id",
	"camera_name",
	"archive_id",
	"archive_name",
	"reviewer",
	"reviewer_type",
	"severity",
	"state",
	"detector_event_id",
	"detector_event_type",
	"macro_event_id",
	"bookmark_id",
	"bookmark_message",
	"bookmark_time_utc",
	"bookmark_time_datetime",
	"bookmark_server_id",
	"bookmark_server_name",
	"bookmark_camera_id",
	"bookmark_camera_name",
	"bookmark_archive_id",
	"bookmark_archive_name",
	"bookmark_alert_id",
	"bookmark_group_id",
	"bookmark_user",
	"bookmark_is_protected",
	"bookmark_boundary_x",
	"bookmark_boundary_y",
	"bookmark_boundary_w",
	"bookmark_boundary_h",
	"bookmark_boundary_index",
	"bookmark_range_time_begin",
	"bookmark_range_time_end",
	"bookmark_geometry_id",
	"bookmark_geometry_alpha",
	"bookmark_geometry_type",
	"bookmark_geometry_ellipse_center_x",
	"bookmark_geometry_ellipse_center_y",
	"bookmark_geometry_ellipse_yr",
	"bookmark_geometry_ellipse_xr",
	"bookmark_geometry_point_x",
	"bookmark_geometry_point_y",
	"bookmark_geometry_rectangle_x",
	"bookmark_geometry_rectangle_y",
	"bookmark_geometry_rectangle_w",
	"bookmark_geometry_rectangle_h",
	"bookmark_geometry_rectangle_index",
	"initiator",
	"initiator_type",
	"reason_mask",
}
var columns = []string{
	"team__id",
	"domain__id",
	"type",
	"time_utc",
	"id",
	"version",
	"datetime",
	"year",
	"month",
	"date",
	"time",
	"hour",
	"hour_x2",
	"hour_x4",
	"min_of_day",
	"min_of_day_x5",
	"min_of_day_x10",
	"min_of_day_x15",
	"min_of_day_x30",
	"nanoseconds",
	"utc_nanoseconds",
	"weekday",
	"detector_type",
	"detector_id",
	"detector_name",
	"server_id",
	"server_name",
	"camera_id",
	"camera_name",
	"phase",
	"multi_phase_id",
	"recognition_quality",
	"detector_face_age",
	"detector_face_gender",
	"detector_face_time_begin",
	"detector_face_time_begin_utc",
	"rectangle_h",
	"rectangle_w",
	"rectangle_x",
	"rectangle_y",
	"rectangle_index",
	"camera_display_id",
	"camera_group_id",
	"camera_group_name",
	"utc_offset",
	"object_class",
}

type Item interface {
	InsertStatement(tableName string, columns []string) string
	Columns() []string
	Values(columns []string) []interface{}
	TableName() string
}

type Event struct {
	js.Object
}

func NewEvent(obj js.Object) Event {
	return Event{obj}
}

func (e Event) TableName() string {
	return "events"
}

func (e Event) Values(columns []string) []interface{} {
	values := make([]interface{}, 0, len(columns))
	for _, column := range columns {
		values = append(values, e.Object[column])
	}
	return values
}

func (e Event) InsertStatement(tableName string, columns []string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ","),
		placeholdersString(len(columns)),
	)
}

func (e Event) Columns() []string {
	return columns
}

type AlertEvent struct {
	Event
}

func NewAlertEvent(evt js.Object) AlertEvent {
	return AlertEvent{Event{Object: evt}}
}

func (e AlertEvent) TableName() string {
	return "alerts"
}
func (e AlertEvent) Columns() []string {
	return alertsColumns
}

func DummyAlertEvent(id string, evtTime time.Time, detectorEventID string, serverID int) (AlertEvent, string) {
	evtTimeStr := evtTime.Format("2006-01-02T15:04:05.000000")
	evtTimeStrUtc := evtTime.UTC().Format("2006-01-02T15:04:05.000000")
	evt := js.NewObject()
	evt["id"] = id
	evt["domain__id"] = domainID
	evt["team__id"] = teamID
	evt["version"] = 1
	evt["type"] = "alert"
	evt["datetime"] = evtTimeStr
	evt["time_utc"] = evtTimeStrUtc
	addDummySource(evt, serverID, false)
	evt["initiator"] = "root"
	evt["initiator_type"] = "AIT_USER"
	evt["reviewer"] = "root"
	evt["reason_mask"] = 4
	evt["detector_event_id"] = detectorEventID
	evt["detector_event_type"] = "faceAppeared"
	evt["macro_event_id"] = macroEventID
	enrichTime(evt)
	return NewAlertEvent(evt), id
}

func DummyAlertEventState(id string, evtTime time.Time, alertType, alertID string, serverID int) AlertEvent {
	evtTimeStr := evtTime.Format("2006-01-02T15:04:05.000000")
	evtTimeStrUtc := evtTime.UTC().Format("2006-01-02T15:04:05.000000")
	evt := js.NewObject()
	evt["id"] = id
	evt["version"] = 1
	evt["domain__id"] = domainID
	evt["team__id"] = teamID
	evt["type"] = "alert_state"
	evt["datetime"] = evtTimeStr
	evt["time_utc"] = evtTimeStrUtc
	addDummySource(evt, serverID, false)
	evt["severity"] = alertType
	evt["reviewer_type"] = "RT_USER"
	evt["reviewer"] = "root"
	evt["state"] = "ST_CLOSED"
	evt["alert_id"] = alertID
	addDummyBookmark(evt, serverID, alertID)
	enrichTime(evt)
	return NewAlertEvent(evt)
}

func DummyFaceAppearedEvent(evtTime time.Time, serverID int) (Event, string) {
	evtTimeStr := evtTime.Format("2006-01-02T15:04:05.000000")
	evtTimeStrUtc := evtTime.UTC().Format("2006-01-02T15:04:05.000000")
	evt := js.NewObject()
	id := uuid.New().String()
	evt["id"] = id
	evt["version"] = 1
	evt["domain__id"] = domainID
	evt["team__id"] = teamID
	evt["type"] = "detector"
	evt["datetime"] = evtTimeStr
	evt["time_utc"] = evtTimeStrUtc
	addDummySource(evt, serverID, true)
	evt["detector_type"] = "faceAppeared"
	evt["detector_face_age"] = 34
	evt["detector_face_gender"] = 2
	evt["detector_face_time_begin"] = evtTimeStr
	evt["detector_face_time_begin_utc"] = evtTimeStrUtc
	evt["detector_people_state"] = "in"
	evt["phase"] = "happened"
	evt["detector_queue_max"] = 2
	evt["recognition_quality"] = 0.45
	evt["detector_listedItem_list_id"] = "bfefb72f-235a-414f-afe7-5303f9d2e50e"
	evt["detector_listedItem_item_id"] = "7a9b3866-8901-4260-9826-470ce34c4219"
	addDummyRectangle(evt)
	evt["detector_listedItem_matched_event_time_utc"] = time.Date(2019, 9, 26, 9, 9, 9, 729000000, time.UTC).Format("2006-01-02T15:04:05.000000")
	evt["multi_phase_id"] = uuid.New().String()
	enrichTime(evt)
	return NewEvent(evt), id
}

func addDummySource(event js.Object, serverID int, withDetector bool) {
	event["server_id"] = fmt.Sprintf("SERVER%d", serverID)
	event["server_name"] = fmt.Sprintf("someServer:%d", serverID)

	event["camera_id"] = "SERVER0/DeviceIpint.1/SourceEndpoint.video:0:0"
	event["camera_name"] = "someCamera"

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
	evtTimeStr := evtTime.Format("2006-01-02T15:04:05.000000")
	evtTimeStrUtc := evtTime.UTC().Format("2006-01-02T15:04:05.000000")
	evt["bookmark_time_datetime"] = evtTimeStr
	evt["bookmark_time_utc"] = evtTimeStrUtc
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
	evt["bookmark_geometry_ellipse_center_x"] = 2
	evt["bookmark_geometry_ellipse_center_y"] = 4
	evt["bookmark_geometry_ellipse_yr"] = 4.5
	evt["bookmark_geometry_ellipse_xr"] = 4.5
	evt["bookmark_range_time_begin"] = evtTimeStrUtc
	evt["bookmark_range_time_end"] = evtTimeStrUtc
}

func enrichTime(event js.Object) {
	localTime := event.GetFieldAsTime("datetime", "2006-01-02T15:04:05")
	event.Put("weekday", strings.ToLower(localTime.Weekday().String()))
	event.Put("year", localTime.Year())
	event.Put("month", int(localTime.Month()))
	event.Put("date", localTime.Format("2006-01-02"))
	event.Put("time", localTime.Format("15:04:05.000"))
	event.Put("hour", localTime.Hour())
	event.Put("nanoseconds", localTime.Nanosecond())
	event.Put("utc_nanoseconds", event.GetFieldAsTime("time_utc").Nanosecond())
	_, utcOffset := localTime.Local().Zone()
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

func placeholdersString(count int) string {
	if count == 0 {
		return ""
	}

	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(placeholders, ",")
}
