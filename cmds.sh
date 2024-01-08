#!/bin/bash
alias air=$(go env GOPATH)/bin/air&&alias swag=./swag

goose -dir "sql/schema" postgres postgres://root:root@localhost:5432/dev4 up

swag init --parseDependency --parseInternal --parseDepth 1