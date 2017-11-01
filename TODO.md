TODO
====

- [x] Do I need json for communication between reader and writer?
  Json removed, simple GET http request.

- [x] Add vendoring

- [x] Add zoom filtering (return http forbidden if z < minzoom or z > maxzoom)

- [ ] Add testing metatiles created by renderd, write tests for reader and writer

- [x] Simplify config file (move to yaml or toml?)

- [ ] Add headers png tiles (Expires, Cache-Control, ETag?)

```
HTTP/1.1 200 OK
Date: Thu, 26 Oct 2017 12:51:54 GMT
Server: Apache/2.4.6 (CentOS)
ETag: "60879f10e430929fb2e11cacd7541b55"
Content-Length: 10439
Cache-Control: max-age=522300
Expires: Wed, 01 Nov 2017 13:56:54 GMT
Content-Type: image/png
```

- [x] Add User-Agent option for httpclient?

```
User-Agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0";
```

- [ ] Use keepalive for fetchservice?

- [ ] Write more docs
