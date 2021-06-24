# wfs-tiler

[<img src="https://img.shields.io/docker/automated/flaviostutz/wfs-tiler"/>](https://hub.docker.com/r/flaviostutz/wfs-tiler)

A Mapbox Vector Tile provider that uses a WFS3 service as data source.
Useful for exposing huge data from WFS3 servers using vector tiles to avoid waste of unnecessary bandwidth and resource usage.

By using vector tiles, the map will fetch data only on the map portions that are meant to be rendered according to viewport and zoom level.

Any WFS3.0 compliant source can be used (ex.: http://github.com/flaviostutz/wfsgis)

Watch a complete demo at https://youtu.be/pRMtTHFqrX0

Check https://github.com/flaviostutz/map-demos for visualization examples.

## Usage

* Create a docker-compose.yml:

```yml
version: '3.7'

services:

  wfs-tiler:
    image: flaviostutz/wfs-tiler
    ports:
      - 3000:3000
    restart: always
    environment:
      - WFS3_API_URL=http://wfsgis:8080
      - CACHE_CONTROL=public,max-age=3600

  map-demos:
    image: flaviostutz/map-demos
    ports:
      - 8181:80
    environment:
      - MAPBOX_VECTOR_TILE_URL=http://wfs-tiler:3000/tiles/tests/{z}/{x}/{y}.mvt

  wfsgis:
    image: flaviostutz/wfsgis
    ports:
      - 8080:8080
    restart: always
    environment:
      - POSTGRES_HOST=postgis
      - POSTGRES_USERNAME=admin
      - POSTGRES_PASSWORD=admin
      - POSTGRES_DBNAME=admin

  postgis:
    image: mdillon/postgis:11-alpine
    ports:
      - 5432:5432
    environment:
      - POSTGRES_USER=admin
      - POSTGRES_PASSWORD=admin
      - POSTGRES_DB=admin
```

* Run 'docker-compose up'

* Runtime diagram

<img src="demo-diagram.png" width="600" />

* Open sample visualization of the tiles at http://localhost:8181/

## ENVs

* WFS3_API_URL - WFS3 server URL. ex.: http://wfsserver.com
* CACHE_CONTROL - Cache control header added to HTTP responses. defaults to 'no-cache'
* LOG_LEVEL - debug,info,warn,error. defaults to 'info'
* SIMPLIFY_LEVEL - level of geometry simplification to be applied. the wider the zoom, the more simplification is applied. 0 for no simplification. defaults to '10'
* MIN_GEOM_LENGTH - depending on zoom level of the tile, 'small' geometries are hidden. this parameter determines this sense of 'small'. defaults to '3600'
* MAX_ZOOM_LEVEL - max zoom level permited. defaults '20'

## Flow

![WFSDiagram](WFSDiagram.png)
