## Run:

1. `make run-services`

Runs clickhouse and postgres servers, creates `generator` database with `alerts` and `events` tables

2. `./smart-generator -detectorsConfigPath=config.json -servers=4 -startDate=2020-11-16 -endDate=2020-11-18`

servers (default: 1) - servers count, on which alerts wil be broadcasting

teams (default: 2) - teams count, specifies to make possible to part events between different teams

detectorsConfigPath (default: ./config.json) - config describe events amount to generate for each detector type, see sample.config.json

startDate (default: 2020-01-01), endDate (default: now) - time range for events generation

generatorType (default: normal) - specifies how to generate time intervals between events, normal: equal intervals between events in time range
between start end dates, random: picks random timestamps in time range between start and end dates

storageType (default: all) - specifies which storage to use for saving events, possible values: clickhouse, postgres, http, all
