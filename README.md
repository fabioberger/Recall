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

Database Setup:
- In psql:
```
CREATE USER recall;
CREATE DATABASE recall_dev;
GRANT ALL PRIVILEGES ON DATABASE recall_dev TO recall;
CREATE DATABASE recall_test;
GRANT ALL PRIVILEGES ON DATABASE recall_test TO recall;
CREATE DATABASE recall_prod;
GRANT ALL PRIVILEGES ON DATABASE recall_prod TO recall;
```

- In recall docker: goose --env=development up
