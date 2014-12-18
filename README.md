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
CREATE DATABASE recall_dev;
CREATE DATABASE recall_test;
CREATE DATABASE recall_prod;
```
- goose --env=development up

Dependencies:
- Negroni (go get github.com/codegangsta/negroni)
- Go-data-parser (go get github.com/albrow/go-data-parser)
- Sendgrid (go get github.com/sendgrid/sendgrid-go)
- gorilla/mux (go get github.com/gorilla/mux)
- Gorp (go get github.com/coopernurse/gorp)
- Goose (go get bitbucket.org/liamstask/goose/cmd/goose)