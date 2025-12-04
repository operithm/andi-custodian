#!/bin/bash

docker-compose up -d

export DATABASE_URL="postgresql://custodian:secure_password_123@localhost:5432/andi_custodian?sslmode=disable"
export TEST_POSTGRES=1

go test -v . -run TestPostgresStore

docker-compose down      # stops containers
## docker-compose down -v   # stops + deletes volume (destroys data)