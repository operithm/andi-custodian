#!/bin/bash

docker build -t operithm/andi-custodian . --progress=plain --no-cache

echo "Docker Build Success."

docker image ls