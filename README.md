## Run:

1. `make run-services`

Runs postgres server, creates `generator` database with `alerts` and `events` tables

2. `./smart-generator -configPath=config.json -servers=4 -startDate=2020-11-16 -endDate=2020-11-18`

servers (default: 4) - servers count, on which alert wil be broadcasting

configPath (default: ./config.json) - config describe events amount to generate for each detector type, see sample.config.json

startDate (default: 2020-01-01), endDate (default: now) - time range for events generation
