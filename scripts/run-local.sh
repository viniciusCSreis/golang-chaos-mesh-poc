#!/bin/bash

make build || exit 1
skaffold run --tail --port-forward
skaffold delete > /dev/null &