# wfs-tiler
A Mapbox Vector Tile provider that uses a WFS3 service as data source.
Useful for exposing huge data from WFS3 servers using vector tiles to avoid waste of unnecessary bandwidth and resource usage.

By using vector tiles, the map will fetch data only on the map portions that are meant to be rendered according to viewport and zoom level.

