Recall: Revision Reminders to Remember Anything
-----------------------------------------------

Environment Variables to set:
- Sendgrid user & key:
	SENDGRID_USER
	SENDGRID_KEY
- Postgres info
	POSTGRES_PORT_5432_TCP_ADDR (localhost if local, already set if linked psql docker container, /var/run/postgresql on Ubuntu aws)
- Session Secrets (64 byte hex string):
RECALL_PROD_SECRET
RECALL_DEV_SECRET
RECALL_TEST_SECRET

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

Cron Tab:

check cron log: vi /var/log/cron

Setup Recall Docker Container:

docker pull fabioberger/recall
docker run -d -p 80:4000 -v /root/volume/.bashrc:/.bashrc --name recall --link db:postgres fabioberger/recall
docker exec -it recall /bin/bash
'source /.bashrc'

Setup Postgres Docker Container:

docker pull fabioberger/recall-postgres
docker run --name db -e POSTGRES_PASSWORD=YOUR_PASSWORD -d db




