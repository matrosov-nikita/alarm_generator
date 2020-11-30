package db

import (
	js "github.com/itimofeev/go-util/json"
	"github.com/matrosov-nikita/smart-generator/events"
)

type Event struct {
	Columns   []string
	Values    []interface{}
	TableName string
}

func ConvertEvents(eventsList []js.Object) []Event {
	dbEvents := make([]Event, 0, len(eventsList))
	var eventColumns []string
	var tableName string
	for _, event := range eventsList {
		eventType := event.GetFieldAsString("type")
		switch eventType {
		case "alert", "alert_state":
			tableName = "alerts"
			eventColumns = alertsColumns
		default:
			tableName = "events"
			eventColumns = columns
		}
		values := make([]interface{}, 0, len(eventColumns))
		for _, column := range eventColumns {
			columnValue := event[column]
			switch v := columnValue.(type) {
			case events.Time:
				values = append(values, v.Time)
			default:
				values = append(values, v)
			}
		}

		dbEvents = append(dbEvents, Event{
			Columns:   eventColumns,
			Values:    values,
			TableName: tableName,
		})
	}

	return dbEvents
}

// 74 columns
var alertsColumns = []string{
	"id",
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

// 68 columns
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
	"detector_people_state",
	"detector_queue_max",
	"detector_queue_min",
	"detector_face_age",
	"detector_face_gender",
	"detector_face_time_begin",
	"detector_face_time_begin_utc",
	"rectangle_h",
	"rectangle_w",
	"rectangle_x",
	"rectangle_y",
	"rectangle_index",
	"detector_lpr_direction",
	"detector_lpr_country",
	"detector_lpr_plate",
	"detector_lpr_best_datetime",
	"detector_lpr_best_utc",
	"detector_lpr_begin_datetime",
	"detector_lpr_begin_utc",
	"detector_lpr_end_datetime",
	"detector_lpr_end_utc",
	"detector_temperature_time_begin",
	"detector_temperature_time_begin_utc",
	"detector_temperature_value",
	"detector_temperature_unit",
	"camera_display_id",
	"camera_group_id",
	"camera_group_name",
	"utc_offset",
	"object_class",
	"detector_listedItem_matched_event_time_utc",
	"detector_listedItem_matched_event_time_datetime",
	"detector_listedItem_matched_event_id",
	"detector_listedItem_list_id",
	"detector_listedItem_item_id",
	"detector_listedFace_score",
}
