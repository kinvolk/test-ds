#!/bin/bash

thisdir="$(dirname "$0")"

"${thisdir}/../testserver/testserver" -config "${thisdir}/server.json"
