#!/bin/bash

wfs-tiler \
  --loglevel="$LOG_LEVEL" \
  --wfs-url="$WFS3_API_URL" \
  --cache-control="$CACHE_CONTROL"
