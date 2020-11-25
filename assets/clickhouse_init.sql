CREATE TABLE IF NOT EXISTS alerts (
    team__id UUID,
    domain__id Int64,
    type String,
    time_utc DateTime,
    id Nullable(UUID),
    alert_id Nullable(UUID),
    version Nullable(UInt16),
    datetime Nullable(DateTime),
    year Nullable(UInt16),
    month Nullable(UInt8),
    date Nullable(Date),
    "time" Nullable(String),
    hour Nullable(UInt8),
    hour_x2 Nullable(UInt8),
    hour_x4 Nullable(UInt8),
    min_of_day Nullable(String),
    min_of_day_x5 Nullable(String),
    min_of_day_x10 Nullable(String),
    min_of_day_x15 Nullable(String),
    min_of_day_x30 Nullable(String),
    nanoseconds Nullable(Int64),
    utc_nanoseconds Nullable(Int64),
    weekday Nullable(String),
    server_id Nullable(String),
    server_name Nullable(String),
    camera_id Nullable(String),
    camera_name Nullable(String),
    archive_id Nullable(String),
    archive_name Nullable(String),
    reviewer Nullable(String),
    reviewer_type Nullable(String),
    severity Nullable(String),
    state Nullable(String),
    detector_event_type Nullable(String),
    macro_event_id Nullable(String),
    bookmark_id Nullable(String),
    bookmark_message Nullable(String),
    bookmark_time_utc Nullable(DateTime),
    bookmark_time_datetime Nullable(DateTime),
    bookmark_server_id Nullable(String),
    bookmark_server_name Nullable(String),
    bookmark_camera_id Nullable(String),
    bookmark_camera_name Nullable(String),
    bookmark_archive_id Nullable(String),
    bookmark_archive_name Nullable(String),
    bookmark_alert_id Nullable(String),
    bookmark_group_id Nullable(String),
    bookmark_user Nullable(String),
    bookmark_is_protected Nullable(UInt8),
    bookmark_boundary_x Nullable(Float64),
    bookmark_boundary_y Nullable(Float64),
    bookmark_boundary_w Nullable(Float64),
    bookmark_boundary_h Nullable(Float64),
    bookmark_boundary_index Nullable(UInt16),
    bookmark_range_time_begin Nullable(DateTime),
    bookmark_range_time_end Nullable(DateTime),
    bookmark_geometry_id Nullable(String),
    bookmark_geometry_alpha Nullable(UInt16),
    bookmark_geometry_type Nullable(String),
    bookmark_geometry_ellipse_center_x Nullable(Float64),
    bookmark_geometry_ellipse_center_y Nullable(Float64),
    bookmark_geometry_ellipse_yr Nullable(Float64),
    bookmark_geometry_ellipse_xr Nullable(Float64),
    bookmark_geometry_point_x Nullable(Float64),
    bookmark_geometry_point_y Nullable(Float64),
    bookmark_geometry_rectangle_x Nullable(Float64),
    bookmark_geometry_rectangle_y Nullable(Float64),
    bookmark_geometry_rectangle_w Nullable(Float64),
    bookmark_geometry_rectangle_h Nullable(Float64),
    bookmark_geometry_rectangle_index Nullable(UInt16),
    initiator Nullable(String),
    initiator_type Nullable(String),
    reason_mask Nullable(UInt16),
    detector_event_id Nullable(UUID),
    utc_offset Nullable(UInt16),
    camera_display_id Nullable(String),
    camera_group_id Nullable(String),
    camera_group_name Nullable(String)
)
ENGINE = MergeTree()
PARTITION BY (type, toYYYYMM(toDate(time_utc)))
ORDER BY (team__id, type, time_utc)
SETTINGS index_granularity=8192;

CREATE TABLE IF NOT EXISTS events (
    team__id UUID,
    domain__id Int64,
    type String,
    time_utc DateTime,
    id Nullable(UUID),
    version Nullable(UInt16),
    "datetime" Nullable(DateTime),
    year Nullable(UInt16),
    month Nullable(UInt8),
    date Nullable(Date),
    "time" Nullable(String),
    hour Nullable(UInt8),
    hour_x2 Nullable(UInt8),
    hour_x4 Nullable(UInt8),
    min_of_day Nullable(String),
    min_of_day_x5 Nullable(String),
    min_of_day_x10 Nullable(String),
    min_of_day_x15 Nullable(String),
    min_of_day_x30 Nullable(String),
    nanoseconds Nullable(Int64),
    utc_nanoseconds Nullable(Int64),
    weekday Nullable(String),
    detector_type Nullable(String),
    detector_id Nullable(String),
    detector_name Nullable(String),
    server_id Nullable(String),
    server_name Nullable(String),
    camera_id Nullable(String),
    camera_name Nullable(String),
    phase Nullable(String),
    multi_phase_id Nullable(UUID),
    recognition_quality Nullable(Float64),
    detector_people_state Nullable(String),
    detector_queue_max Nullable(Int64),
    detector_queue_min Nullable(Int64),
    detector_face_age Nullable(UInt8),
    detector_face_gender Nullable(String),
    detector_face_time_begin Nullable(DateTime),
    detector_face_time_begin_utc Nullable(DateTime),
    rectangle_h Nullable(Float64),
    rectangle_w Nullable(Float64),
    rectangle_x Nullable(Float64),
    rectangle_y Nullable(Float64),
    rectangle_index Nullable(UInt16),
    detector_lpr_direction Nullable(UInt16),
    detector_lpr_country Nullable(String),
    detector_lpr_plate Nullable(String),
    detector_lpr_best_datetime Nullable(DateTime),
    detector_lpr_best_utc Nullable(DateTime),
    detector_lpr_begin_datetime Nullable(DateTime),
    detector_lpr_begin_utc Nullable(DateTime),
    detector_lpr_end_datetime Nullable(DateTime),
    detector_lpr_end_utc Nullable(DateTime),
    "detector_listedItem_matched_event_time_utc" Nullable(DateTime),
    "detector_listedItem_matched_event_time_datetime" Nullable(DateTime),
    "detector_listedItem_matched_event_id" Nullable(UUID),
    "detector_listedItem_list_id" Nullable(UUID),
    "detector_listedItem_item_id" Nullable(UUID),
    "detector_listedFace_score" Nullable(Float64),
    detector_temperature_time_begin Nullable(DateTime),
    detector_temperature_time_begin_utc Nullable(DateTime),
    detector_temperature_value Nullable(Float64),
    detector_temperature_unit Nullable(String),
    camera_display_id Nullable(String),
    camera_group_id Nullable(String),
    camera_group_name Nullable(String),
    utc_offset Nullable(UInt16),
    object_class Nullable(String)
)
ENGINE = MergeTree()
PARTITION BY (type, toYYYYMM(toDate(time_utc)))
ORDER BY (team__id, type, time_utc)
SETTINGS index_granularity=8192;
