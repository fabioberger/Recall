Recall: Revision Reminders to Remember Anything
-----------------------------------------------

Environment Variables to set:
- Sendgrid user & key:
	SENDGRID_USER
	SENDGRID_KEY
- Users email
	GMAIL_ADDRESS
- Postgres info
	POSTGRES_PORT_5432_TCP_ADDR (localhost if local, already set if linked psql docker container, /var/run/postgresql on Ubuntu aws)
	POSTGRES_USER
	POSTGRES_PASSWORD

Database Setup:
- In psql:
```
CREATE USER recall WITH PASSWORD 'testing' ENCRYPTED;
CREATE DATABASE recall_dev;
GRANT ALL PRIVILEGES ON DATABASE recall_dev TO recall;
CREATE DATABASE recall_test;
GRANT ALL PRIVILEGES ON DATABASE recall_test TO recall;
CREATE DATABASE recall_prod;
GRANT ALL PRIVILEGES ON DATABASE recall_prod TO recall;
```
- goose --env=development up

Dependencies:
- Negroni (go get github.com/codegangsta/negroni)
- Go-data-parser (go get github.com/albrow/go-data-parser)
- Sendgrid (go get github.com/sendgrid/sendgrid-go)
- gorilla/mux (go get github.com/gorilla/mux)
- Gorp (go get github.com/coopernurse/gorp)
- Goose (go get bitbucket.org/liamstask/goose/cmd/goose)