from datetime import datetime
from suncalc import get_times
from pydantic import BaseModel
from pydantic_extra_types.coordinate import Coordinate
from global_land_mask import globe
from is_sky_dark import get_brightness_ratio
import math


class Location(BaseModel):
    coordinate: Coordinate


class BoundingBoxCoords(BaseModel):
    max: Coordinate
    min: Coordinate


def create_grid(bounding_box: BoundingBoxCoords, spacing: int = 5) -> list[Location]:
    min_lat = bounding_box.min.latitude
    max_lat = bounding_box.max.latitude
    min_lon = bounding_box.min.longitude
    max_lon = bounding_box.max.longitude
    lat_spacing = spacing / 111.32  # 1 degree of latitude is approximately 111.32 km
    lon_spacing = spacing / (111.32 * math.cos(math.radians((min_lat + max_lat) / 2)))

    # Calculate the number of grid points in latitude and longitude directions
    num_lat_points = int((max_lat - min_lat) / lat_spacing) + 1
    num_lon_points = int((max_lon - min_lon) / lon_spacing) + 1

    # Create the grid of latitude and longitude coordinates
    grid = []
    for i in range(num_lat_points):
        for j in range(num_lon_points):
            lat = round(min_lat + i * lat_spacing, 4)
            lon = round(min_lon + j * lon_spacing, 4)
            if globe.is_land(lat, lon) and get_brightness_ratio(lat, lon, 2022) < 2:
                grid.append(Location(coordinate=(lat, lon)))
    return grid


UK_BOUNDING_BOX = BoundingBoxCoords(min=Coordinate(50, -4.5), max=Coordinate(53, 1.7))

grid = create_grid(UK_BOUNDING_BOX)
