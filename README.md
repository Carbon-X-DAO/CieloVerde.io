# Description
This repository hosts a webserver that connects to a `postgres` DB instance and 
stores information submitted to a form. So, there are a few dependencies.

# Getting Started
## `go` webserver
If you have a `go` compiler installed then you can build and run the webserver
by
```bash
go build ./cmd/server && ./server
```

## `postgres` DB instance
You will need to install `postgres` and configure it. An installation from 
source (not with a `systemd` service target) can be used in the following way.
```bash
initdb db # a `postman` helper binary
cp ./postgresql.conf db/postgresql.conf
pg_ctl -D db -l logfile start
```

When you're all done you can
```bash
pg_ctl -D db -l logfile stop
rm logfile
```