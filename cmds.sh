#!/bin/bash

source .env

goose -dir "sql/schema" postgres $DBURL up

docker build -t blank .
docker run -p 8080:$PORT blank

pg_dump 'postgresql://myuser:mypassword@db.myproject.supabase.co:5432/mydatabase > database-dump.sql'
psql 'postgresql://myuser:mypassword@db.myproject.supabase.co:5432/mydatabase' < database-dump.sql