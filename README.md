osin-json-rest
==================

[![Build Status](https://travis-ci.org/martint17r/osin-mongo-storage.svg?branch=master)](https://travis-ci.org/martint17r/osin-mongo-storage)

This package uses osin with [go-json-rest](https://github.com/ant0ine/go-json-rest) and implements the storage interface for [OSIN](https://github.com/RangelReale/osin) with [MongoDB](http://www.mongodb.org/) using [mgo](http://labix.org/mgo).

[![baby-gopher](https://raw.githubusercontent.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)

Docker
------
The shell scripts under bin/ build a docker image and execute the tests. Make sure that you can run docker without sudo.

Caveats
-------

All structs are serialized as is, i.e. no references are created or resolved.

Currently MongoDB >= 2.6 is required, on 2.4 the TestLoad* Tests fail, but I do not know why.


Related projects
-----
This project's form middleware  inspired by the following projects:

<https://github.com/boonep/go-json-rest-middleware-formjson>

Examples
--------

See the examples subdirectory for integrating into OSIN.

TODO
-------
- [ ] Add the cache support
- [ ] Add the https support 

License
-------
This package is made available under the [MIT License](http://github.com/martint17r/osin-mongo-storage/LICENSE)
