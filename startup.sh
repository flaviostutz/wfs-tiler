#!/bin/bash

wfs-tiler \
  --loglevel="$LOG_LEVEL" \
  --wfs-url="$WFS3_API_URL" \
  --cache-control="$CACHE_CONTROL" \
  --simplify-level="$SIMPLIFY_LEVEL" \
  --min-geom-length="$MIN_GEOM_LENGTH" \
  --max-zoom-level="$MAX_ZOOM_LEVEL"
