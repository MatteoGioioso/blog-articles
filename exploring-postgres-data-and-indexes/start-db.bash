#!/usr/bin/env bash

docker build -t postgres-local . &&
docker run --rm -d \
  -p 5432:5432 \
  -v "${PWD}"/pg_data:/var/lib/postgresql/data \
  --env POSTGRES_DB=test \
  --env POSTGRES_PASSWORD=123 \
  --name postgres \
  postgres-local
