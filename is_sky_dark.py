import requests
import gzip
import numpy as np
from functools import lru_cache
import math
from io import BytesIO


@lru_cache(maxsize=None)
def get_tile_data(url):
    response = requests.get(url, timeout=5)
    data = gzip.open(BytesIO(response.content)).read()
    return np.frombuffer(data, dtype=np.int8)


def get_brightness_ratio(lat: float, lon: float, year):
    lon_from_dateline = (lon + 180.0) % 360.0
    lat_from_start = lat + 65.0
    tilex = math.floor(lon_from_dateline / 5.0 + 1)
    tiley = math.floor(lat_from_start / 5.0 + 1)
    if 1 <= tiley <= 28:
        url = f"https://djlorenz.github.io/astronomy/binary_tiles/{year}/binary_tile_{tilex}_{tiley}.dat.gz"
        data_array = get_tile_data(url)
        ix = round(120 * (lon_from_dateline - 5.0 * (tilex - 1) + 1.0 / 240.0))
        iy = round(120 * (lat_from_start - 5.0 * (tiley - 1) + 1.0 / 240.0))
        first_number = 128 * data_array[0] + data_array[1]
        change = 0.0
        for i in range(1, iy):
            change += data_array[600 * i + 1]
        for i in range(1, ix):
            change += data_array[600 * (iy - 1) + 1 + i]
        compressed = first_number + change
        brightness_ratio = (5.0 / 195.0) * (np.exp(0.0195 * compressed) - 1.0)
        return brightness_ratio
    else:
        raise ValueError("Latitude out of range")
