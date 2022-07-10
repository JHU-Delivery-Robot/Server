#!/usr/bin/env bash

rm -r OSRM/data
mkdir OSRM/data

# download OpenStreetMap data of JHU campus
curl "https://www.openstreetmap.org/api/0.6/map?bbox=-76.6272,39.3227,-76.6144,39.337" > OSRM/data/map.osm

# Run OSRM data processing pipeline
docker run -t -v "${PWD}/OSRM:/OSRM" osrm/osrm-backend osrm-extract -p /OSRM/profiles/wheelchairelektro.lua /OSRM/data/map.osm
docker run -t -v "${PWD}/OSRM:/OSRM" osrm/osrm-backend osrm-partition /OSRM/data/map.osrm
docker run -t -v "${PWD}/OSRM:/OSRM" osrm/osrm-backend osrm-customize /OSRM/data/map.osrm
