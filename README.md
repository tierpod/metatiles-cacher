metatiles-cacher
================

metatiles-cacher contains a few small services:

1) metatiles_cacher - for serving tiles from metatiles cache. Contains slippy-map based on [LeafLet][1].
   If tile not found in cache, download from remote source and write to metatiles cache.

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
