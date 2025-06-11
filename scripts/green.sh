#!/bin/bash

width=15

curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "update"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "bgrect $width $width $((800-width)) $((800-width))"
curl -X POST http://localhost:17000 -d "update"