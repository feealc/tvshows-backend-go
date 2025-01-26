#! /bin/bash

unset CLOUD_NAME

export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="root"
export DB_PASS="root"
export DB_NAME="root"

# Go
go run main.go
