#!/bin/bash
alias air=$(go env GOPATH)/bin/air&&alias swag=./swag

goose -dir "sql/schema" postgres postgres://root:root@localhost:5432/localdevdb up
