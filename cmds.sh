#!/bin/bash

source .env

goose -dir "sql/schema" postgres $DBURL up

docker build -t blank .
docker run -p 8080:$PORT blank