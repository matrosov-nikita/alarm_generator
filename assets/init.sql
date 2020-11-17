CREATE TABLE IF NOT EXISTS alerts (
    type text,
    time_utc timestamp without time zone,
    id uuid,
    alert_id uuid,
    version integer,
    datetime timestamp without time zone,
    year integer,
    month smallint,
    date date,
    "time" text,
    hour smallint,
    hour_x2 smallint,
    hour_x4 smallint,
    min_of_day text,
    min_of_day_x5 text,
    min_of_day_x10 text,
    min_of_day_x15 text,
    min_of_day_x30 text,
    nanoseconds bigint,
    utc_nanoseconds bigint,
    weekday text,
    server_id text,
    server_name text,
    camera_id text,
    camera_name text,
    archive_id text,
    archive_name text,
    reviewer text,
    reviewer_type text,
    severity text,
    state text,
    detector_event_type text,
    macro_event_id text,
    bookmark_id text,
    bookmark_message text,
    bookmark_time_utc timestamp without time zone,
    bookmark_time_datetime timestamp without time zone,
    bookmark_server_id text,
    bookmark_server_name text,
    bookmark_camera_id text,
    bookmark_camera_name text,
    bookmark_archive_id text,
    bookmark_archive_name text,
    bookmark_alert_id text,
    bookmark_group_id text,
    bookmark_user text,
    bookmark_is_protected smallint,
    bookmark_boundary_x double precision,
    bookmark_boundary_y double precision,
    bookmark_boundary_w double precision,
    bookmark_boundary_h double precision,
    bookmark_boundary_index integer,
    bookmark_range_time_begin timestamp without time zone,
    bookmark_range_time_end timestamp without time zone,
    bookmark_geometry_id text,
    bookmark_geometry_alpha integer,
    bookmark_geometry_type text,
    bookmark_geometry_ellipse_center_x double precision,
    bookmark_geometry_ellipse_center_y double precision,
    bookmark_geometry_ellipse_yr double precision,
    bookmark_geometry_ellipse_xr double precision,
    bookmark_geometry_point_x double precision,
    bookmark_geometry_point_y double precision,
    bookmark_geometry_rectangle_x double precision,
    bookmark_geometry_rectangle_y double precision,
    bookmark_geometry_rectangle_w double precision,
    bookmark_geometry_rectangle_h double precision,
    bookmark_geometry_rectangle_index integer,
    initiator text,
    initiator_type text,
    reason_mask integer,
    detector_event_id uuid,
    utc_offset integer,
    camera_display_id text,
    camera_group_id text,
    camera_group_name text
);

CREATE TABLE IF NOT EXISTS events (
    type text,
    time_utc timestamp without time zone,
    id uuid,
    version integer,
    datetime timestamp without time zone,
    year integer,
    month smallint,
    date date,
    "time" text,
    hour smallint,
    hour_x2 smallint,
    hour_x4 smallint,
    min_of_day text,
    min_of_day_x5 text,
    min_of_day_x10 text,
    min_of_day_x15 text,
    min_of_day_x30 text,
    nanoseconds bigint,
    utc_nanoseconds bigint,
    weekday text,
    detector_type text,
    detector_id text,
    detector_name text,
    server_id text,
    server_name text,
    camera_id text,
    camera_name text,
    phase text,
    multi_phase_id uuid,
    recognition_quality double precision,
    detector_face_age smallint,
    detector_face_gender text,
    detector_face_time_begin timestamp without time zone,
    detector_face_time_begin_utc timestamp without time zone,
    rectangle_h double precision,
    rectangle_w double precision,
    rectangle_x double precision,
    rectangle_y double precision,
    rectangle_index integer,
    camera_display_id text,
    camera_group_id text,
    camera_group_name text,
    utc_offset integer,
    object_class text
);
