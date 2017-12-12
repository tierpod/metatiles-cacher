metatiles-cacher
================

metatiles-cacher contains:

1) metatiles-cacher - daemon for serving tiles from metatiles cache. If tile not found in cache,
download from remote source and write to metatiles cache.

Contains slippy-map based on [LeafLet][1] for png tiles. For vector tiles, you can use [Tangram][5]
(download and put it in static directory).

2) convert-latlong - converts latitude and longitude to z, x, y format

[Workflow][4]:

```
  request  +------+  not found   +------+
+---------->reader+-------+------>source|
           +---^--+       |      +---+--+
               |          |          |
               |          |          |download
               |          |          |
               |          |      +---v--+
               |          +------>writer|
          found|                 +---+--+
               |                     |
               |                     |write
               |                     |
           +---+---------------------v--+
           |    metatiles file cache    |
           +----------------------------+
```

Based on
--------

* [gopnik tile server][2]
* [gosm library][3]


Entrypoints
-----------

* http://localhost:8080/static/ - slippy-map.

* http://localhost:8080/maps/{style}/{z}/{x}/{y}.{ext} - read tile from metatiles cache. If tile not
  found in cache, fetch from remote source and write to cache.

  Returns http status:

  * StatusInternalServerError - if error occured
  * StatusNotFound - if tile not found in the source, or unknown mimetype
  * StatusNotModified - if tile not modified since last request
  * StatusForbidden - if tile has wrong zoom level
  * StatusOK - if tile serves successful

* http://localhost:8080/fetch/{style}/{z}/{x}/{y}.{ext} - fetch tile from remote source and write to
  metatiles cache.

  Returns http status:

  * StatusInternalServerError - if error occured
  * StatusNotFound - if tile not found in the source, or unknown mimetype
  * StatusForbidden - if tile has wrong zoom level
  * StatusCreated - if tile already in the fetch queue (try later)
  * StatusOK - if tile serves successful

Region files
------------

Regions can be kml files downloaded from geofabrik.de. Or you can convert this kml files to yaml.
Examples can be found at pkg/config/testdata directory.


[1]: http://leafletjs.com
[2]: https://github.com/sputnik-maps/gopnik
[3]: https://github.com/apeyroux/gosm
[4]: http://asciiflow.com/
[5]: https://mapzen.com/documentation/vector-tiles/display-tiles/
