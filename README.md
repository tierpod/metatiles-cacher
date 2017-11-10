metatiles-cacher
================

metatiles-cacher contains:

1) metatiles_cacher - daemon for serving tiles from metatiles cache. If tile not found in cache,
download from remote source and write to metatiles cache.

Contains slippy-map based on [LeafLet][1] for png tiles. For vector tiles, you can use [Tangram][5]
(download and put it in static directory).

2) convert_latlong - converts latitude and longitude to z, x, y format

[Workflow][4]:

```
  request  +------+  not found   +------+
+---------->reader+-------+----->+source|
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



[1]: http://leafletjs.com
[2]: https://github.com/sputnik-maps/gopnik
[3]: https://github.com/apeyroux/gosm
[4]: http://asciiflow.com/
[5]: https://mapzen.com/documentation/vector-tiles/display-tiles/
