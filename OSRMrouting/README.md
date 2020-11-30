This file takes a start and endpoint and uses the OSRM route service to collect a series of waypoints between them

OSRM is running with the bicycle profile on OSM data for Maryland

To get the OSRM server running on your machine:

    Have Docker installed

    Download the files at this URL: https://drive.google.com/file/d/1vEiwCTzZzxqqrPCnKWbm8-CgFXXfG_l1/view

    Run these 2 commands in terminal:
        docker run -t -v "${PWD}:/data" osrm/osrm-backend osrm-partition /data/maryland-latest.osrm
        docker run -t -v "${PWD}:/data" osrm/osrm-backend osrm-customize /data/maryland-latest.osrm

    Now you're all set up! To run the server:
        docker run -t -i -p 5000:5000 -v "$(pwd):/data" osrm/osrm-backend osrm-routed --algorithm mld /data/maryland-latest.osrm