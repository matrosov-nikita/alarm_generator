## Run:

1. `make run-services`

Runs postgres server, creates `generator` database with `alerts` and `events` tables

2. `./alarm_generator -alertsCount=2500 -servers=4 -startDate=2020-11-16 -endDate=2020-11-18`

servers (default: 4) - servers count, on which alert wil be broadcasting

alertsCount (default: 2000)- alerts count, for each alert _servers*2+1_ events will be generated

startDate (default: 2020-01-01), endDate (default: now) - time range for events generation