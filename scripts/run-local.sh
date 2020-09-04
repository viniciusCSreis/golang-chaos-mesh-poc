#!/bin/bash

make build
skaffold run --tail --port-forward
skaffold delete > /dev/null &