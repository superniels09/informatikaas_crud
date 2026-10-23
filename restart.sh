#!/bin/bash
cd "$(dirname -- "$(readlink -f -- "$0")")"
pkill informatikaas
nohup ./informatikaas_crud > output.log 2>&1 &
