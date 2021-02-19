#!/bin/bash

thisdir="$(dirname "$0")"

"${thisdir}/../../envoy-build/envoy/source/exe/envoy" --config-path "${thisdir}/envoy.yaml"
