#!/bin/sh 

docker-compose --profile=debug up -d
xdg-open http://localhost:8081
