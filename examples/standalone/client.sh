#!/bin/bash

thisdir="$(dirname "$0")"

"${thisdir}/../testclient/testclient" -config "${thisdir}/client.json" "${@}"
