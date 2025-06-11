#!/bin/bash

Xstart=100
Ystart=50
step=20

curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "update"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "figure $Xstart $Ystart"
curl -X POST http://localhost:17000 -d "update"

while true; do
    curl -X POST http://localhost:17000 -d "update"
    for ((i = 0; i < 280; i += step)); do
        curl -X POST http://localhost:17000 -d "move $((-step)) $((step))"
        curl -X POST http://localhost:17000 -d "update"
    done
    for ((i = 600; i < 900; i += step)); do
        curl -X POST http://localhost:17000 -d "move $step $((-step))"
        curl -X POST http://localhost:17000 -d "update"
    done
done