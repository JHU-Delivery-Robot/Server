import requests
from pprint import pprint
import matplotlib.pyplot as plt
import geopandas
from shapely.geometry import Point

# start and endpoint coordinates in the following format: LONGITUDE,LATITUDE
# both points must be somewhere on campus:
#    Longitude range: []
#    Latitude range: []
START = '-76.619772,39.328796'
END = '-76.620832,39.331518'

r = requests.get('http://localhost:5000/route/v1/bicycle/' + START + ';' + END + '?steps=true&geometries=geojson')
# r is a Response object

print('')
print('OVERVIEW = SIMPLIFIED')
print('')

# will get a list of the geometry coordinates, a list of the maneuver coordinates, a list of the intersection coordinates, 
# and a list of the routestep geometries

# there are several different lists of waypoints that can be extracted from the Route object
# theya re described below (order and numbers are arbitrary)
# 1) Simplified geometry - waypoints produced by the route by default
# 2) Maneuvers - points where the robot does a maneuver (a turn, departing or arriving)
# 3) Intersections - waypoint as every intersecion the passed on the route
# 4) Full geometry - waypoints produced by the route with the parameter overview = full

route = r.json()['routes'][0]

# geometry of the entire route
simple_geometry_points = route['geometry']['coordinates']
simple_geometry_points = [tuple(point) for point in simple_geometry_points]
print('route geometry')
print(simple_geometry_points)
print(len(simple_geometry_points))

# points where a maneuver (depart, turn, arrive) is made
maneuver_points = []
for routestep in route['legs'][0]['steps']:
    maneuver_points.append(tuple(routestep['maneuver']['location']))
print('maneuvers')
print(maneuver_points)
print(len(maneuver_points))

# every intersection the robot crosses on its path
intersection_points = []
for routestep in route['legs'][0]['steps']:
    for intersection in routestep['intersections']:
        intersection_points.append(tuple(intersection['location']))
print('intersection points')
print(intersection_points)
print(len(intersection_points))


print('')
print('OVERVIEW = FULL')
print('')


r2 = requests.get('http://localhost:5000/route/v1/bicycle/' + START + ';' + END + '?steps=true&geometries=geojson&overview=full')
route2 = r2.json()['routes'][0]

# full geometry of the entire route
full_geometry_points = route2['geometry']['coordinates']
full_geometry_points = [tuple(point) for point in full_geometry_points]
print('route2 geometry')
print(full_geometry_points)
print(len(full_geometry_points))

# note: coordinate lists for maneuvers and intersections are exactly the same for the default and overview=full requests

#####################################################################################################################
# We can compare the different waypoint lists by looking at them on maps! :)
# following code saves a png containing 4 maps of campus, with one of the waypoint lists plotted on each

# map of campus
campus_map = geopandas.read_file('./campusMap/roads-line.shp')
fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize = (18.00, 18.00), dpi = 600, tight_layout = True)

# if we wanted to save each map as a different image
'''
fig1,ax1 = plt.subplots(figsize = (19.20,19.20))
fig2,ax2 = plt.subplots()
fig3,ax3 = plt.subplots()
fig4,ax4 = plt.subplots()
'''

campus_map.plot(ax = ax1, alpha = 0.4, color = 'grey')
simpleGS = geopandas.GeoSeries([Point(pt) for pt in simple_geometry_points])
simpleGS.plot(ax = ax1, color = 'blue', markersize = 0.5)
ax1.set_title('Simple Geometry')
#ax1.get_xaxis().set_visible(False)
#ax1.get_yaxis().set_visible(False)
#fig1.savefig('./simpleGeo.png')


campus_map.plot(ax = ax2, alpha = 0.4, color = 'grey')
maneuverGS = geopandas.GeoSeries([Point(pt) for pt in maneuver_points])
maneuverGS.plot(ax = ax2, color = 'red', markersize = 0.5)
ax2.set_title('Maneuvers')
#fig2.savefig('./maneuvers.png')

campus_map.plot(ax = ax3, alpha = 0.4, color = 'grey')
intersectionGS = geopandas.GeoSeries([Point(pt) for pt in intersection_points])
intersectionGS.plot(ax = ax3, color = 'green', markersize = 0.5)
ax3.set_title('Intersections')
#fig3.savefig('./intersections.png')


campus_map.plot(ax = ax4, alpha = 0.4, color = 'grey')
fullGS = geopandas.GeoSeries([Point(pt) for pt in full_geometry_points])
fullGS.plot(ax = ax4, color = 'orange', markersize = 0.5)
ax4.set_title('Full Geometry')
#fig5.savefig('./fullGeo.png')


fig.savefig('./waypoints.png')