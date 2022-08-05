#!/bin/sh 

docker-compose -f docker-compose.debug.yml up -d
xdg-open http://localhost:8081
